package fragadapter_demo

import (
	def "github.com/akley-MK4/net-defragmenter/definition"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

func NewAdapterRecord(id AdapterRecordIdType, inst IAdapterInstance) *AdapterRecord {
	return &AdapterRecord{
		id:              id,
		inst:            inst,
		capturedInfoMap: make(map[def.FragGroupID]*CapturedInfo),
	}
}

type AdapterRecord struct {
	id              AdapterRecordIdType
	inst            IAdapterInstance
	capturedInfoMap map[def.FragGroupID]*CapturedInfo
	mutex           sync.Mutex
}

func (t *AdapterRecord) start() {

}

func (t *AdapterRecord) stop() {

}

func (t *AdapterRecord) close() {

}

func (t *AdapterRecord) release() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	mapLen := len(t.capturedInfoMap)
	if mapLen <= 0 {
		return
	}

	keys := make([]def.FragGroupID, 0, mapLen)
	for key := range t.capturedInfoMap {
		keys = append(keys, key)
	}
	for _, key := range keys {
		delete(t.capturedInfoMap, key)
	}
}

func (t *AdapterRecord) associateCapturedInfo(fragGroupID def.FragGroupID, timestamp time.Time, ifIndex int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exist := t.capturedInfoMap[fragGroupID]; exist {
		return
	}

	t.capturedInfoMap[fragGroupID] = &CapturedInfo{
		CreateTp:  time.Now().Unix(),
		Timestamp: time.Unix(int64(timestamp.Second()), int64(timestamp.Nanosecond())),
		IfIndex:   ifIndex,
	}
}

func (t *AdapterRecord) reassemblyCapturedBuf(fullPkt *def.FullPacket) {
	fragGroupID := fullPkt.GetFragGroupID()
	t.mutex.Lock()
	info, exist := t.capturedInfoMap[fragGroupID]
	if exist {
		delete(t.capturedInfoMap, fragGroupID)
	}
	t.mutex.Unlock()

	if info == nil {
		fullPkt.Pkt = nil
		log.Printf("[warning][reassemblyPcapBuf] The info with fragGroup %v dose not exists\n", fragGroupID)
		return
	}

	pkt := fullPkt.GetPacket()
	pktData := pkt.Data()
	fullPkt.Pkt = nil

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Catch ReassemblyCompletedCallback exception, Recover: %v, Stack: %v", r, string(debug.Stack()))
		}
	}()

	t.inst.ReassemblyCompletedCallback(info.Timestamp, info.IfIndex, pktData)
}
