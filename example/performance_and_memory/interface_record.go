package main

import (
	"errors"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/implement"
	PCD "github.com/akley-MK4/pep-coroutine/define"
	PCI "github.com/akley-MK4/pep-coroutine/implement"
	"log"
	"sync"
	"time"
)

var (
	interfaceRecordMgr = &InterfaceRecordMgr{
		interfaceRecordMap: make(map[def.InterfaceId]*InterfaceRecord),
		interfaceIdMapping: make(map[string]def.InterfaceId),
	}
)

func initInterfaceRecordMgr(intervalPopFullPktTime time.Duration, maxPopFullPacketsNum int) error {
	interfaceRecordMgr.intervalPopFullPktTime = intervalPopFullPktTime
	interfaceRecordMgr.maxPopFullPacketsNum = maxPopFullPacketsNum
	return nil
}

func getInterfaceRecordMgr() *InterfaceRecordMgr {
	return interfaceRecordMgr
}

type InterfaceRecordMgr struct {
	rwMutex sync.RWMutex

	incInterfaceId         uint16
	intervalPopFullPktTime time.Duration
	maxPopFullPacketsNum   int
	interfaceIdMapping     map[string]def.InterfaceId
	interfaceRecordMap     map[def.InterfaceId]*InterfaceRecord
}

func (t *InterfaceRecordMgr) start() error {
	if err := PCI.CreateAndStartStatelessCoroutine("checkAndPopFullPacketsPeriodically", func(coID PCD.CoId, args ...interface{}) bool {
		t.checkAndPopFullPacketsPeriodically()
		return false
	}); err != nil {
		return err
	}

	return nil
}

func (t *InterfaceRecordMgr) RegisterInterfaceRecord(interfaceName string) def.InterfaceId {
	t.rwMutex.Lock()
	id, exist := t.interfaceIdMapping[interfaceName]
	if !exist {
		t.incInterfaceId += 1
		id = def.InterfaceId(t.incInterfaceId)
		t.interfaceIdMapping[interfaceName] = id
		t.interfaceRecordMap[id] = newInterfaceRecord(interfaceName, id)
	}
	t.rwMutex.Unlock()

	return id
}

func (t *InterfaceRecordMgr) DeregisterInterfaceRecord(interfaceId def.InterfaceId) bool {
	t.rwMutex.Lock()
	record, exist := t.interfaceRecordMap[interfaceId]
	if !exist {
		t.rwMutex.Unlock()
		return false
	}
	delete(t.interfaceRecordMap, interfaceId)
	t.rwMutex.Unlock()

	record.release()
	return true
}

func (t *InterfaceRecordMgr) getInterfaceRecord(interfaceId def.InterfaceId) (retVal *InterfaceRecord) {
	t.rwMutex.RLock()
	retVal = t.interfaceRecordMap[interfaceId]
	t.rwMutex.RUnlock()
	return
}

func (t *InterfaceRecordMgr) checkAndPopFullPacketsPeriodically() {
	for {
		time.Sleep(t.intervalPopFullPktTime)

		fullPks, popErr := getLibInstance().PopFullPackets(t.maxPopFullPacketsNum)
		if popErr != nil {
			log.Printf("Failed to pop full packets from Net-DeFragment Lib, %v\n", popErr)
			continue
		}
		if len(fullPks) <= 0 {
			continue
		}

		for _, fullPkt := range fullPks {
			record := t.getInterfaceRecord(fullPkt.GetInterfaceId())
			if record == nil {
				log.Printf("Unable to find InterfaceRecord with id %v in LibProxy\n", fullPkt.GetInterfaceId())
				continue
			}

			record.reassemblyCapturedFrags(fullPkt)
		}
	}
}

func (t *InterfaceRecordMgr) AssociateCapturedInfo(detectInfo *def.DetectionInfo, capSeconds int64, capNanoseconds int) error {
	record := t.getInterfaceRecord(detectInfo.InterfaceId)
	if record == nil {
		return errors.New("unable to find the specified Record")
	}
	record.associateCapturedInfo(detectInfo.FragGroupId, capSeconds, capNanoseconds)
	return nil
}

type CapturedInfo struct {
	CapSeconds     int64
	CapNanoseconds int
	CreateTs       int64
}

func newInterfaceRecord(interfaceName string, interfaceId def.InterfaceId) *InterfaceRecord {
	return &InterfaceRecord{
		id:         interfaceId,
		name:       interfaceName,
		capInfoMap: make(map[def.FragGroupID]*CapturedInfo),
	}
}

type InterfaceRecord struct {
	name       string
	id         def.InterfaceId
	capInfoMap map[def.FragGroupID]*CapturedInfo
	mutex      sync.Mutex
}

func (t *InterfaceRecord) release() {
	t.mutex.Lock()
	t.capInfoMap = make(map[def.FragGroupID]*CapturedInfo)
	t.mutex.Unlock()
}

func (t *InterfaceRecord) associateCapturedInfo(fragGroupID def.FragGroupID, capSeconds int64, capNanoseconds int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exist := t.capInfoMap[fragGroupID]; exist {
		return
	}

	t.capInfoMap[fragGroupID] = &CapturedInfo{
		CreateTs:       time.Now().Unix(),
		CapSeconds:     capSeconds,
		CapNanoseconds: capNanoseconds,
	}
}

func (t *InterfaceRecord) reassemblyCapturedFrags(fullPkt *def.FullPacket) {
	defer implement.ReleaseFullPacket(fullPkt)

	fragGroupID := fullPkt.GetFragGroupID()
	t.mutex.Lock()
	info, exist := t.capInfoMap[fragGroupID]
	if exist {
		delete(t.capInfoMap, fragGroupID)
	}
	t.mutex.Unlock()

	if info == nil {
		log.Printf("The CapInfo with fragGroup %v dose not exist\n", fragGroupID)
		return
	}

	pkt := fullPkt.GetPacket()
	_ = pkt.Data()

}
