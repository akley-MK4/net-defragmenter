package collection

import (
	"context"
	"fmt"
	"time"

	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/internal/common"
	"github.com/akley-MK4/net-defragmenter/internal/handler"
	"github.com/akley-MK4/net-defragmenter/internal/linkqueue"
	"github.com/akley-MK4/net-defragmenter/stats"
	PCD "github.com/akley-MK4/pep-coroutine/define"
	PCI "github.com/akley-MK4/pep-coroutine/implement"
)

func newCollector(id, maxListenChanCap uint32, ptrFullPktQueue *linkqueue.LinkQueue) *Collector {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	return &Collector{
		id:               id,
		cancelCtx:        cancelCtx,
		cancelFunc:       cancelFunc,
		fragmentChan:     make(chan *common.FragElement, maxListenChanCap),
		fragElemGroupMap: make(map[def.FragGroupID]*common.FragElementGroup),
		ptrFullPktQueue:  ptrFullPktQueue,
		sharedLayers:     common.NewSharedLayers(),
	}
}

type Collector struct {
	id         uint32
	cancelCtx  context.Context
	cancelFunc context.CancelFunc

	fragmentChan     chan *common.FragElement
	fragElemGroupMap map[def.FragGroupID]*common.FragElementGroup
	ptrFullPktQueue  *linkqueue.LinkQueue
	sharedLayers     *common.SharedLayers
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
			break
		case <-sharedLayersRestTimer.C:
			if t.sharedLayers.GetReferencesNum() > 0 {
				t.sharedLayers.Reset()
			}
			break
		case frag, ok := <-t.fragmentChan:
			if !ok {
				break loopExit
			}
			if err := t.accept(frag); err != nil {
				// todo
			}
			break
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
	for _, fragElemGroup := range t.fragElemGroupMap {
		if (nowTp - fragElemGroup.GetCreateTimestamp()) > maxFragGroupDurationSec {
			expiredGroups = append(expiredGroups, fragElemGroup)
		}
	}

	for _, fragElemGroup := range expiredGroups {
		delete(t.fragElemGroupMap, fragElemGroup.GetID())
		fragElemGroup.Release()
	}

	if len(expiredGroups) > 0 {
		stats.GetCollectionStatsHandler().AddTotalReleasedExpiredFragElementGroupsNum(uint64(len(expiredGroups)))
	}
}

func (t *Collector) accept(fragElem *common.FragElement) error {
	statsHandler := stats.GetCollectionStatsHandler()
	statsHandler.AddTotalAcceptedFragElementsNum(1)

	hd := handler.GetHandler(fragElem.Type)
	if hd == nil {
		statsHandler.AddTotalNotFoundHandlersNum(1)
		return fmt.Errorf("handler with fragment type %v dose not exists", fragElem.Type)
	}

	fragElemGroup, exist := t.fragElemGroupMap[fragElem.GroupID]
	if !exist {
		t.fragElemGroupMap[fragElem.GroupID] = common.NewFragElementGroup(fragElem.GroupID)
		fragElemGroup = t.fragElemGroupMap[fragElem.GroupID]
	}

	collectErr, collectErrType := hd.Collect(fragElem, fragElemGroup)
	if collectErr != nil {
		statsHandler.AddTotalErrCollectNum(1, collectErrType)
		return collectErr
	}

	statsHandler.AddTotalSuccessfulCollectedFragsNum(1)
	if err := t.checkAndReassembly(fragElemGroup, fragElem, hd); err != nil {
		return err
	}

	return nil
}

func (t *Collector) checkAndReassembly(fragElemGroup *common.FragElementGroup, fragElem *common.FragElement, hd handler.IHandler) error {
	if !fragElemGroup.CheckFinalElementExists() || fragElemGroup.GetHighest() != fragElemGroup.GetCurrentLen() {
		return nil
	}

	statsHandler := stats.GetCollectionStatsHandler()
	if _, exist := t.fragElemGroupMap[fragElemGroup.GetID()]; exist {
		delete(t.fragElemGroupMap, fragElemGroup.GetID())
	} else {
		statsHandler.AddTotalReassemblyNoDelFragGroupsNum(1)
	}

	fragElemListLen := fragElemGroup.GetElementListLen()
	defer func() {
		fragElemGroup.Release()
	}()

	pkt, reassemblyErr, errType := hd.Reassembly(fragElemGroup, t.sharedLayers)
	t.sharedLayers.UpdateReferences()

	if reassemblyErr != nil {
		statsHandler.AddTotalErrReassemblyNum(1, errType)
		return reassemblyErr
	}
	statsHandler.AddTotalSuccessfulReassemblyFragsNum(uint64(fragElemListLen))

	statsHandler.AddTotalPushedFullPacketsNum(1)
	t.ptrFullPktQueue.SafetyPutValue(&def.FullPacket{
		InterfaceId:  fragElem.InterfaceId,
		FragGroupID:  fragElem.GroupID,
		Pkt:          pkt,
		FragElemsNum: fragElemListLen,
	})

	return nil
}

func (t *Collector) pushFragmentElement(fragElem *common.FragElement) {
	t.fragmentChan <- fragElem
}
