package collection

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/internal/common"
	"github.com/akley-MK4/net-defragmenter/internal/handler"
	"github.com/akley-MK4/net-defragmenter/internal/linkqueue"
	"github.com/akley-MK4/net-defragmenter/stats"
	PCD "github.com/akley-MK4/pep-coroutine/define"
	PCI "github.com/akley-MK4/pep-coroutine/implement"
)

var (
	ErrFragGroupMapReachedLenLimit = errors.New("reached the maximum length limit")
)

func newCollector(id, maxListenChanCap, maxFragGroupMapLength uint32, ptrFullPktQueue *linkqueue.LinkQueue, enableSyncReassembly bool) *Collector {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	return &Collector{
		id:                    id,
		cancelCtx:             cancelCtx,
		cancelFunc:            cancelFunc,
		fragmentChan:          make(chan *common.FragElement, maxListenChanCap),
		fragElemGroupMap:      make(map[def.FragGroupID]*common.FragElementGroup),
		maxFragGroupMapLength: int(maxFragGroupMapLength),
		ptrFullPktQueue:       ptrFullPktQueue,
		sharedLayers:          common.NewSharedLayers(),
		enableSyncReassembly:  enableSyncReassembly,
	}
}

type Collector struct {
	id         uint32
	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	fragmentChan          chan *common.FragElement
	fragElemGroupMap      map[def.FragGroupID]*common.FragElementGroup
	maxFragGroupMapLength int
	ptrFullPktQueue       *linkqueue.LinkQueue
	sharedLayers          *common.SharedLayers

	enableSyncReassembly bool
	syncReassemblyMutex  sync.RWMutex
}

func (t *Collector) start() error {
	return PCI.CreateAndStartStatelessCoroutine(def.CoroutineGroupCollector1, func(_ PCD.CoId, _ ...interface{}) bool {
		t.schedulingCoroutine()
		return false
	})
}

func (t *Collector) close() {
	t.cancelFunc()
}

func (t *Collector) stop() {

}

func (t *Collector) schedulingCoroutine() {
	fragSetCheckTimer := time.NewTicker(time.Second * time.Duration(intervalCheckFragSec))
	sharedLayersRestTimer := time.NewTicker(time.Second * time.Duration(intervalRestSharedLayersFragSec))

loopExit:
	for {
		select {
		case <-fragSetCheckTimer.C:
			t.checkFragmentElementSetExpired()
		case <-sharedLayersRestTimer.C:
			if t.sharedLayers.GetReferencesNum() > 0 {
				t.sharedLayers.Reset()
			}
		case frag, ok := <-t.fragmentChan:
			if !ok {
				break loopExit
			}
			if _, err := t.accept(frag); err != nil {
				// todo
			}
		case <-t.cancelCtx.Done():
			break loopExit
		}
	}

	fragSetCheckTimer.Stop()
	sharedLayersRestTimer.Stop()

	t.stop()
}

func (t *Collector) checkFragmentElementSetExpired() {
	nowTp := time.Now().Unix()
	var expiredGroups []*common.FragElementGroup

	if t.enableSyncReassembly {
		t.syncReassemblyMutex.Lock()
	}

	for _, fragElemGroup := range t.fragElemGroupMap {
		if (nowTp - fragElemGroup.GetCreateTimestamp()) > maxFragGroupDurationSec {
			expiredGroups = append(expiredGroups, fragElemGroup)
		}
	}
	for _, fragElemGroup := range expiredGroups {
		delete(t.fragElemGroupMap, fragElemGroup.GetID())
		fragElemGroup.Release()
	}

	if t.enableSyncReassembly {
		t.syncReassemblyMutex.Unlock()
	}

	if len(expiredGroups) > 0 {
		stats.GetCollectionStatsHandler().AddTotalReleasedExpiredFragElementGroupsNum(uint64(len(expiredGroups)))
	}
}

func (t *Collector) accept(fragElem *common.FragElement) (*def.FullPacket, error) {
	statsHandler := stats.GetCollectionStatsHandler()
	statsHandler.AddTotalAcceptedFragElementsNum(1)

	hd := handler.GetHandler(fragElem.Type)
	if hd == nil {
		common.RecycleFragElement(fragElem)
		statsHandler.AddTotalNotFoundHandlersNum(1)
		return nil, fmt.Errorf("handler with fragment type %v dose not exists", fragElem.Type)
	}

	if t.enableSyncReassembly {
		t.syncReassemblyMutex.Lock()
	}
	fragElemGroup, exist := t.fragElemGroupMap[fragElem.GroupID]
	if !exist {
		if len(t.fragElemGroupMap) >= t.maxFragGroupMapLength {
			common.RecycleFragElement(fragElem)
			statsHandler.IncTotalFragMapReachedLenLimitNum()
			return nil, ErrFragGroupMapReachedLenLimit
		}

		t.fragElemGroupMap[fragElem.GroupID] = common.NewFragElementGroup(fragElem.GroupID)
		fragElemGroup = t.fragElemGroupMap[fragElem.GroupID]
	}
	if t.enableSyncReassembly {
		t.syncReassemblyMutex.Unlock()
	}

	collectErr, collectErrType := hd.Collect(fragElem, fragElemGroup)
	if collectErr != nil {
		common.RecycleFragElement(fragElem)
		statsHandler.AddTotalErrCollectNum(1, collectErrType)
		return nil, collectErr
	}

	statsHandler.AddTotalCollectedFragsNum(1)
	var sharedLayers *common.SharedLayers = nil
	if t.enableSyncReassembly {
		sharedLayers = common.NewSharedLayers()
	} else {
		sharedLayers = t.sharedLayers
	}

	return t.checkAndReassembly(fragElemGroup, fragElem, hd, sharedLayers)
}

func (t *Collector) checkAndReassembly(fragElemGroup *common.FragElementGroup, fragElem *common.FragElement, hd handler.IHandler,
	sharedLayers *common.SharedLayers) (*def.FullPacket, error) {

	if !fragElemGroup.CheckFinalElementExists() || fragElemGroup.GetHighest() != fragElemGroup.GetCurrentLen() {
		return nil, nil
	}

	statsHandler := stats.GetCollectionStatsHandler()

	if t.enableSyncReassembly {
		t.syncReassemblyMutex.Lock()
	}
	if _, exist := t.fragElemGroupMap[fragElemGroup.GetID()]; exist {
		delete(t.fragElemGroupMap, fragElemGroup.GetID())
	} else {
		statsHandler.AddTotalReassemblyNoDelFragGroupsNum(1)
	}
	if t.enableSyncReassembly {
		t.syncReassemblyMutex.Unlock()
	}

	fragElemListLen := fragElemGroup.GetElementListLen()
	defer func() {
		fragElemGroup.Release()
	}()

	pkt, reassemblyErr, errType := hd.Reassembly(fragElemGroup, sharedLayers)
	if t.enableSyncReassembly {
		sharedLayers.Reset()
	} else {
		sharedLayers.UpdateReferences()
	}

	if reassemblyErr != nil {
		statsHandler.AddTotalErrReassemblyNum(1, errType)
		return nil, reassemblyErr
	}
	statsHandler.AddTotalReassemblyFragsNum(uint64(fragElemListLen))

	statsHandler.AddTotalReassemblyFullPacketsNum(1)
	fullPkt := &def.FullPacket{
		InterfaceId:  fragElem.InterfaceId,
		FragGroupID:  fragElem.GroupID,
		Pkt:          pkt,
		FragElemsNum: fragElemListLen,
	}

	if t.enableSyncReassembly {
		return fullPkt, nil
	}

	t.ptrFullPktQueue.SafetyPutValue(&def.FullPacket{
		InterfaceId:  fragElem.InterfaceId,
		FragGroupID:  fragElem.GroupID,
		Pkt:          pkt,
		FragElemsNum: fragElemListLen,
	})

	return nil, nil
}

func (t *Collector) pushFragmentElement(fragElem *common.FragElement) {
	t.fragmentChan <- fragElem
}
