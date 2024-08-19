package common

import (
	"bytes"
	"container/list"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/stats"
	"github.com/google/gopacket/layers"
	"sync"
	"time"
)

type FragElement struct {
	GroupID     def.FragGroupID
	Type        def.FragType
	InterfaceId def.InterfaceId

	SrcMAC, DstMAC []byte
	SrcIP, DstIP   []byte
	TOS            uint8
	TrafficClass   uint8
	FlowLabel      uint32
	IPProtocol     layers.IPProtocol
	FragOffset     uint16
	MoreFrags      bool
	Identification uint32

	PayloadBuf *bytes.Buffer
	Grouped    bool
}

var (
	fragElementObjectPool = &sync.Pool{
		New: func() any {
			stats.GetCollectionStatsHandler().AddTotalAllocatedFragElementsNum(1)
			return &FragElement{
				PayloadBuf: &bytes.Buffer{},
			}
		},
	}
)

func NewFragElement() *FragElement {
	stats.GetCollectionStatsHandler().AddTotalNewFragElementsNum(1)
	elem := fragElementObjectPool.Get().(*FragElement)
	elem.PayloadBuf.Reset()
	elem.Grouped = false
	return elem
}

func RecycleFragElement(elem *FragElement) {
	if elem == nil {
		return
	}

	stats.GetCollectionStatsHandler().AddTotalRecycledFragElementsNum(1)
	elem.PayloadBuf.Reset()
	fragElementObjectPool.Put(elem)
}

func NewFragElementGroup(fragGroupID def.FragGroupID) *FragElementGroup {
	stats.GetCollectionStatsHandler().AddTotalNewFragElementGroupsNum(1)
	return &FragElementGroup{
		groupID:  fragGroupID,
		createTp: time.Now().Unix(),
		elemList: list.New(),
	}
}

type FragElementGroup struct {
	groupID      def.FragGroupID
	createTp     int64
	elemList     *list.List
	highest      uint16
	currentLen   uint16
	ptrFinalElem *FragElement
	nextProtocol interface{}
	lastSeen     time.Time
}

func (t *FragElementGroup) GetID() def.FragGroupID {
	return t.groupID
}

func (t *FragElementGroup) GetHighest() uint16 {
	return t.highest
}

func (t *FragElementGroup) AddHighest(val uint16) uint16 {
	t.highest += val
	return t.highest
}

func (t *FragElementGroup) SetHighest(val uint16) {
	t.highest = val
}

func (t *FragElementGroup) GetCurrentLen() uint16 {
	return t.currentLen
}

func (t *FragElementGroup) SetNextProtocol(proto interface{}) {
	t.nextProtocol = proto
}

func (t *FragElementGroup) GetNextProtocol() interface{} {
	return t.nextProtocol
}

func (t *FragElementGroup) AddCurrentLen(val uint16) uint16 {
	t.currentLen += val
	return t.currentLen
}

func (t *FragElementGroup) CheckFinalElementExists() bool {
	return t.ptrFinalElem != nil
}

func (t *FragElementGroup) SetFinalElement(elem *FragElement) {
	t.ptrFinalElem = elem
}

func (t *FragElementGroup) GetFinalElement() *FragElement {
	return t.ptrFinalElem
}

func (t *FragElementGroup) PushElementToBack(elem *FragElement) {
	t.elemList.PushBack(elem)
}

func (t *FragElementGroup) InsertElementToBefore(elem *FragElement, mark *list.Element) *list.Element {
	return t.elemList.InsertBefore(elem, mark)
}

func (t *FragElementGroup) IterElementList(f func(elem *list.Element) bool) {
	for e := t.elemList.Front(); e != nil; e = e.Next() {
		if !f(e) {
			return
		}
	}
}

func (t *FragElementGroup) GetElementListLen() int {
	return t.elemList.Len()
}

func (t *FragElementGroup) Release() (clenListLen int) {
	clenListLen = t.cleanUpElementList()
	t.ptrFinalElem = nil
	t.elemList = nil
	stats.GetCollectionStatsHandler().AddTotalReleasedFragElementGroupsNum(1)
	return
}

func (t *FragElementGroup) cleanUpElementList() int {
	if t.elemList == nil {
		return 0
	}

	var elems []*list.Element
	for e := t.elemList.Front(); e != nil; e = e.Next() {
		elems = append(elems, e)
	}
	for _, elem := range elems {
		v := t.elemList.Remove(elem)
		if v == nil {
			continue
		}

		fragElem, ok := v.(*FragElement)
		if !ok {
			continue
		}
		RecycleFragElement(fragElem)
	}

	return len(elems)
}

func (t *FragElementGroup) GetCreateTimestamp() int64 {
	return t.createTp
}

func (t *FragElementGroup) GetAllElementsPayloadLen() uint16 {
	var totalPayloadLen int
	for e := t.elemList.Front(); e != nil; e = e.Next() {
		elem := e.Value.(*FragElement)
		totalPayloadLen += elem.PayloadBuf.Len()
	}
	return uint16(totalPayloadLen)
}
