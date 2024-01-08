package implement

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

func NewLibraryInstance(opt *def.Option) (*Library, error) {
	if opt == nil {
		return nil, errors.New("opt is a nil pointer")
	}

	inst := &Library{}
	if err := inst.initialize(opt); err != nil {
		return nil, err
	}

	return inst, nil
}

type Library struct {
	status       int32
	detector     *detection.Detector
	collectorMgr *collection.CollectorMgr
}

func (t *Library) initialize(opt *def.Option) error {
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

func (t *Library) Start() error {
	if !atomic.CompareAndSwapInt32(&t.status, def.InitializedStatus, def.StartedStatus) {
		return errors.New("incorrect state")
	}

	return t.collectorMgr.Start()
}

func (t *Library) Stop() {
	if !atomic.CompareAndSwapInt32(&t.status, def.StartedStatus, def.StoppedStatus) {
		return
	}

	t.collectorMgr.Stop()
}

func (t *Library) AsyncProcessPacket(interfaceId def.InterfaceId, pktData []byte, onDetectCompleted def.OnDetectCompleted) (retErr error) {
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
	defer detectInfo.Rest()

	if err := t.detector.FastDetect(interfaceId, pktData, &detectInfo); err != nil {
		return err
	}
	if detectInfo.FragType == def.NonFragType {
		onDetectCompleted(detectInfo.FragType, "")
		return
	}

	onDetectCompleted(detectInfo.FragType, detectInfo.FragGroupId)

	if err := t.collectorMgr.Collect(&detectInfo); err != nil {
		return
	}

	return
}

func (t *Library) PopFullPackets(count int) ([]*def.FullPacket, error) {
	if t.status != def.StartedStatus {
		return nil, fmt.Errorf("manager not started, current status is %v", t.status)
	}

	if t.collectorMgr == nil {
		return nil, errors.New("collectorMgr is a nil pointer")
	}

	return t.collectorMgr.PopFullPackets(count)
}

func (t *Library) FastDetect(interfaceId def.InterfaceId, pktData []byte, replyDetectInfo *def.DetectionInfo) error {
	if t.status != def.StartedStatus {
		return fmt.Errorf("manager not started, current status is %v", t.status)
	}
	if replyDetectInfo == nil {
		return errors.New("the replyDetectInfo is a nil value")
	}

	replyDetectInfo.FragType = def.NonFragType
	return t.detector.FastDetect(interfaceId, pktData, replyDetectInfo)
}

func (t *Library) Collect(detectInfo *def.DetectionInfo) error {
	if t.status != def.StartedStatus {
		return fmt.Errorf("manager not started, current status is %v", t.status)
	}
	if detectInfo == nil {
		return errors.New("the detectInfo is a nil value")
	}

	if detectInfo.FragGroupId == "" {
		return errors.New("the FragGroupId of the detectInfo is not generated")
	}
	if detectInfo.InterfaceId == 0 {
		return errors.New("the InterfaceId of the detectInfo is 0")
	}

	return t.collectorMgr.Collect(detectInfo)
}

func ReleaseFullPacket(fullPkt *def.FullPacket) {
	collection.ReleaseFullPacket(fullPkt)
}
