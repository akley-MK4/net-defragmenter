package manager

import (
	"errors"
	"fmt"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/internal/collection"
	"github.com/akley-MK4/net-defragmenter/internal/detection"
	"github.com/akley-MK4/net-defragmenter/stats"
	"runtime/debug"
	"sync/atomic"
)

func NewManager(opt *def.Option) (*Manager, error) {
	if opt == nil {
		return nil, errors.New("opt is a nil pointer")
	}

	mgr := &Manager{}
	if err := mgr.initialize(opt); err != nil {
		return nil, err
	}

	return mgr, nil
}

type Manager struct {
	status       int32
	detector     *detection.Detector
	collectorMgr *collection.CollectorMgr
}

func (t *Manager) initialize(opt *def.Option) error {
	if opt.StatsOption.Enable {
		stats.EnableStats()
	} else {
		stats.DisableStats()
	}

	detector, newDetectorErr := detection.NewDetector(opt.PickFragmentTypes)
	if newDetectorErr != nil {
		return fmt.Errorf("NewDetector failed, %v", newDetectorErr)
	}

	collectorMgr, newCollectorErr := collection.NewCollectorMgr(opt.CollectorOption)
	if newCollectorErr != nil {
		return fmt.Errorf("NewCollectorMgr failed, %v", newCollectorErr)
	}

	t.detector = detector
	t.collectorMgr = collectorMgr
	t.status = def.InitializedStatus

	return nil
}

func (t *Manager) Start() error {
	if !atomic.CompareAndSwapInt32(&t.status, def.InitializedStatus, def.StartedStatus) {
		return errors.New("incorrect state")
	}

	return t.collectorMgr.Start()
}

func (t *Manager) Stop() {
	if !atomic.CompareAndSwapInt32(&t.status, def.StartedStatus, def.StoppedStatus) {
		return
	}

	t.collectorMgr.Stop()
}

func (t *Manager) AsyncProcessPacket(pktBuf []byte, inMarkValue uint64, onDetectCompleted def.OnDetectCompleted) (retErr error) {
	defer func() {
		if r := recover(); r != nil {
			retErr = fmt.Errorf("catch AsyncProcessPacket exception, Recover: %v, Stack: %v", r, string(debug.Stack()))
			stats.GetManagerStatsHandler().AddTotalCrashedAsyncProcessNum(1)
		}
		if retErr != nil {
			stats.GetManagerStatsHandler().AddTotalFailedCallAsyncProcessNum(1)
		}
	}()

	if t.status != def.StartedStatus {
		retErr = fmt.Errorf("manager not started, current status is %v", t.status)
		return
	}
	if onDetectCompleted == nil {
		retErr = errors.New("the onDetectCompleted is a nil value")
		return
	}

	var detectInfo def.DetectionInfo
	detectInfo.FragType = def.NonFragType

	if err := t.detector.FastDetect(pktBuf, &detectInfo); err != nil {
		return err
	}
	if detectInfo.FragType == def.NonFragType {
		onDetectCompleted(detectInfo.FragType, "")
		return
	}

	fragGroupID := detectInfo.GenFragGroupID()
	onDetectCompleted(detectInfo.FragType, fragGroupID)

	t.collectorMgr.Collect(fragGroupID, &detectInfo, inMarkValue)

	detectInfo.Rest()
	return
}

func (t *Manager) PopFullPackets(count int) ([]*def.FullPacket, error) {
	if t.status != def.StartedStatus {
		return nil, fmt.Errorf("manager not started, current status is %v", t.status)
	}

	if t.collectorMgr == nil {
		return nil, errors.New("collectorMgr is a nil pointer")
	}

	return t.collectorMgr.PopFullPackets(count)
}
