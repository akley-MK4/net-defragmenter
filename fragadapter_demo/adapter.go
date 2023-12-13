package fragadapter_demo

import (
	"encoding/json"
	"errors"
	"fmt"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/manager"
	"github.com/akley-MK4/net-defragmenter/stats"
	PCD "github.com/akley-MK4/pep-coroutine/define"
	PCI "github.com/akley-MK4/pep-coroutine/implement"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type CapturedInfo struct {
	Timestamp time.Time
	IfIndex   int
	CreateTp  int64
}

const (
	initializeStatus int32 = iota
	initializedStatus
	startedStatus
	stoppedStatus
)

const (
	maxPullFullPacketsNum      = 5000
	popPullPktInterval         = time.Second * time.Duration(3)
	enableStatsFile            = true
	intervalUpdateStatsFileMin = 3
	statsFilePath              = "/tmp/vtap_defragment.json"
)

type AdapterRecordIdType uint32

type IAdapterInstance interface {
	ReassemblyCompletedCallback(timestamp time.Time, ifIndex int, buf []byte)
}

type NewDeFragmentLibFunc func() (IDeFragmentLib, error)

type IDeFragmentLib interface {
	Start() error
	Stop()
	AsyncProcessPacket(pktBuf []byte, userMarkValue uint32, onDetectCompleted def.OnDetectCompleted) error
	PopFullPackets(count int) ([]*def.FullPacket, error)
	FastDetect(pktData []byte, replyDetectInfo *def.DetectionInfo) error
	Collect(detectInfo *def.DetectionInfo, userMarkValue uint32) error
}

var (
	adapterInstance *DeFragmentAdapter
)

func InitializeAdapterInstance() error {
	if adapterInstance != nil {
		return nil
	}

	inst, newInstErr := NewDeFragmentAdapter()
	if newInstErr != nil {
		return newInstErr
	}

	adapterInstance = inst
	log.Println("[info][DeFragmentAdapter] AdapterInstance initialization successful")
	return nil
}

func GetAdapterInstance() *DeFragmentAdapter {
	return adapterInstance
}

func NewDeFragmentAdapter() (*DeFragmentAdapter, error) {
	opt := def.NewOption(func(opt *def.Option) {
		//opt.CtrlApiServerOption.Enable = true
		opt.StatsOption.Enable = true

		opt.PickFragmentTypes = []def.FragType{def.IPV4FragType, def.IPV6FragType}

		opt.CollectorOption.MaxCollectorsNum = 30
		opt.CollectorOption.MaxChannelCap = 2000
		opt.CollectorOption.MaxFullPktQueueLen = 10000
	})

	lib, newLibErr := manager.NewManager(opt)
	if newLibErr != nil {
		return nil, newLibErr
	}

	adapter := &DeFragmentAdapter{
		status:    initializedStatus,
		lib:       lib,
		recordMap: make(map[AdapterRecordIdType]*AdapterRecord),
	}

	return adapter, nil
}

type DeFragmentAdapter struct {
	status      int32
	lib         IDeFragmentLib
	incRecordId AdapterRecordIdType
	recordMap   map[AdapterRecordIdType]*AdapterRecord
	rwMutex     sync.RWMutex
}

func (t *DeFragmentAdapter) Start() error {
	if !atomic.CompareAndSwapInt32(&t.status, initializedStatus, startedStatus) {
		return errors.New("incorrect state")
	}

	if err := t.lib.Start(); err != nil {
		return err
	}

	nowTime := time.Now()
	if err := PCI.CreateAndStartStatelessCoroutine(PCD.CoGroup(fmt.Sprintf("listenReassemblyCompleted-%v", nowTime.String())), func(coID PCD.CoId, args ...interface{}) bool {
		t.listenReassemblyCompleted()
		return false
	}); err != nil {
		return err
	}

	if enableStatsFile {
		if err := PCI.CreateAndStartStatelessCoroutine(PCD.CoGroup(fmt.Sprintf("updateStatsFilePeriodically-%v", nowTime.String())), func(coID PCD.CoId, args ...interface{}) bool {
			updateStatsFilePeriodically()
			return false
		}); err != nil {
			return err
		}
	}

	return nil
}

func (t *DeFragmentAdapter) Stop() {
	if !atomic.CompareAndSwapInt32(&t.status, startedStatus, stoppedStatus) {
		return
	}

	t.clearUpRecords()
	t.lib.Stop()
}

func (t *DeFragmentAdapter) getRecord(id AdapterRecordIdType) *AdapterRecord {
	t.rwMutex.RLock()
	record := t.recordMap[id]
	t.rwMutex.RUnlock()
	return record
}

func (t *DeFragmentAdapter) clearUpRecords() {
	delMap := make(map[AdapterRecordIdType]*AdapterRecord)
	t.rwMutex.Lock()
	for id, record := range t.recordMap {
		delMap[id] = record
	}
	t.rwMutex.Unlock()

	for _, record := range delMap {
		record.release()
	}
}

func (t *DeFragmentAdapter) RegisterInstance(inst IAdapterInstance) (retId AdapterRecordIdType) {
	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()

	t.incRecordId += 1
	retId = t.incRecordId
	t.recordMap[retId] = NewAdapterRecord(retId, inst)
	log.Printf("[info][DeFragmentAdapter] Registered a new PCAP instance, RecordId: %v\n", retId)

	return
}

func (t *DeFragmentAdapter) DeregisterInstance(id AdapterRecordIdType) {
	t.rwMutex.Lock()
	delInstRecord, exist := t.recordMap[id]
	if !exist {
		t.rwMutex.Unlock()
		log.Printf("[warning][DeFragmentAdapter] Deregister failed, The record %v dose not exists\n", id)
		return
	}

	delete(t.recordMap, id)
	t.rwMutex.Unlock()

	delInstRecord.release()
	log.Printf("[info][DeFragmentAdapter] Deregistered instance with RecordId %v\n", id)
}

func (t *DeFragmentAdapter) AsyncProcessPacket(id AdapterRecordIdType, timestamp time.Time, ifIndex int, buf []byte) bool {
	t.rwMutex.RLock()
	record := t.recordMap[id]
	t.rwMutex.RUnlock()
	if record == nil {
		log.Printf("[warning][DeFragmentAdapter] AsyncProcessPacket failed, The record %v dose not exists\n", id)
		return false
	}

	var fragGroupID def.FragGroupID
	var processErr error
	processErr = t.lib.AsyncProcessPacket(buf, uint32(record.id), func(fragType def.FragType, fragGroupID def.FragGroupID) {
		record.associateCapturedInfo(fragGroupID, timestamp, ifIndex)
	})
	if processErr != nil {
		log.Printf("[warning][DeFragmentAdapter] AsyncProcessPacket failed, %v\n", processErr)
		return false
	}

	return fragGroupID != ""
}

func (t *DeFragmentAdapter) SyncProcessPacket(id AdapterRecordIdType, timestamp time.Time, ifIndex int, buf []byte) bool {
	t.rwMutex.RLock()
	record := t.recordMap[id]
	t.rwMutex.RUnlock()
	if record == nil {
		log.Printf("[warning][DeFragmentAdapter] SyncProcessPacket failed, The record %v dose not exists\n", id)
		return false
	}

	var detectInfo def.DetectionInfo
	defer detectInfo.Rest()

	if err := t.lib.FastDetect(buf, &detectInfo); err != nil {
		log.Printf("[warning][DeFragmentAdapter] SyncProcessPacket, FastDetect failed, %v\n", err)
		return false
	}
	if detectInfo.FragGroupId == "" {
		return false
	}

	record.associateCapturedInfo(detectInfo.FragGroupId, timestamp, ifIndex)
	if err := t.lib.Collect(&detectInfo, uint32(id)); err != nil {
		log.Printf("[warning][DeFragmentAdapter] SyncProcessPacket, Collect failed, %v\n", err)
		return false
	}

	return true
}

func (t *DeFragmentAdapter) listenReassemblyCompleted() {
	for {
		time.Sleep(popPullPktInterval)

		fullPktList, popErr := t.lib.PopFullPackets(maxPullFullPacketsNum)
		if popErr != nil {
			log.Printf("[warning][DeFragmentAdapter] Call listenReassemblyCompleted failed, PopFullPackets error, %v\n", popErr)
			continue
		}
		if len(fullPktList) <= 0 {
			continue
		}

		for _, pkt := range fullPktList {
			recordId := AdapterRecordIdType(pkt.GetUserMarkValue())
			record := t.getRecord(recordId)
			if record == nil {
				log.Printf("[warning][DeFragmentAdapter] Call listenReassemblyCompleted failed, The record %v dose not exists\n", pkt.GetUserMarkValue())
				continue
			}

			record.reassemblyCapturedBuf(pkt)
		}
	}
}

func updateStatsFilePeriodically() {
	interval := time.Minute * intervalUpdateStatsFileMin
	for {
		time.Sleep(interval)
		if err := updateStatsFile(); err != nil {
			log.Printf("[warning][DeFragmentAdapter] updateStatsFilePeriodically failed, %v\n", err)
		}
	}
}

func updateStatsFile() (retErr error) {
	defer func() {
		if r := recover(); r != nil {
			retErr = fmt.Errorf("catch updateStatsFile exception, Recover: %v, Stack: %v", r, string(debug.Stack()))
		}
	}()

	f, openErr := os.OpenFile(statsFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if openErr != nil {
		retErr = openErr
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			retErr = err
		}
	}()

	d, marshalErr := json.Marshal(stats.GetStats())
	if marshalErr != nil {
		retErr = marshalErr
		return
	}

	_, writeErr := f.Write(d)
	if writeErr != nil {
		retErr = writeErr
		return
	}

	return
}
