package common

import (
	"bytes"
	"container/list"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/libstats"
	"github.com/google/gopacket/layers"
	"sync"
	"time"
)

type FragmentElement struct {
	GroupID     def.FragmentGroupID
	Type        def.FragmentType
	InMarkValue uint64

	SrcMAC, DstMAC []byte
	SrcIP, DstIP   []byte
	IPProtocol     layers.IPProtocol
	FragOffset     uint16
	MoreFrags      bool
	Identification uint32

	PayloadBuf *bytes.Buffer
}

var (
	fragElementObjectPool = &sync.Pool{
		New: func() any {
			libstats.AddTotalAllocateFragmentElementNum(1)
			return &FragmentElement{
				PayloadBuf: &bytes.Buffer{},
			}
		},
	}
)

func NewFragmentElement() *FragmentElement {
	libstats.AddTotalNewFragmentElementNum(1)
	elem := fragElementObjectPool.Get().(*FragmentElement)
	elem.PayloadBuf.Reset()

	return elem
}

func RecycleFragmentElement(elem *FragmentElement) {
	if elem == nil {
		return
	}

	libstats.AddTotalRecycleFragmentElementNum(1)
	elem.PayloadBuf.Reset()
	fragElementObjectPool.Put(elem)
}

func NewFragmentElementGroup(fragGroupID def.FragmentGroupID) *FragmentElementGroup {
	return &FragmentElementGroup{
		groupID:  fragGroupID,
		createTp: time.Now().Unix(),
		elemList: list.New(),
	}
}

type FragmentElementGroup struct {
	groupID      def.FragmentGroupID
	createTp     int64
	elemList     *list.List
	highest      uint16
	currentLen   uint16
	ptrFinalElem *FragmentElement
	nextProtocol interface{}
	lastSeen     time.Time
}

func (t *FragmentElementGroup) GetID() def.FragmentGroupID {
	return t.groupID
}

func (t *FragmentElementGroup) GetHighest() uint16 {
	return t.highest
}

func (t *FragmentElementGroup) AddHighest(val uint16) uint16 {
	t.highest += val
	return t.highest
}

func (t *FragmentElementGroup) SetHighest(val uint16) {
	t.highest = val
}

func (t *FragmentElementGroup) GetCurrentLen() uint16 {
	return t.currentLen
}

func (t *FragmentElementGroup) SetNextProtocol(proto interface{}) {
	t.nextProtocol = proto
}

func (t *FragmentElementGroup) GetNextProtocol() interface{} {
	return t.nextProtocol
}

func (t *FragmentElementGroup) AddCurrentLen(val uint16) uint16 {
	t.currentLen += val
	return t.currentLen
}

func (t *FragmentElementGroup) CheckFinalElementExists() bool {
	return t.ptrFinalElem != nil
}

func (t *FragmentElementGroup) SetFinalElement(elem *FragmentElement) {
	t.ptrFinalElem = elem
}

func (t *FragmentElementGroup) GetFinalElement() *FragmentElement {
	return t.ptrFinalElem
}

func (t *FragmentElementGroup) PushElementToBack(elem *FragmentElement) {
	t.elemList.PushBack(elem)
}

func (t *FragmentElementGroup) InsertElementToBefore(elem *FragmentElement, mark *list.Element) *list.Element {
	return t.elemList.InsertBefore(elem, mark)
}

func (t *FragmentElementGroup) IterElementList(f func(elem *list.Element) bool) {
	for e := t.elemList.Front(); e != nil; e = e.Next() {
		if !f(e) {
			return
		}
	}
}

func (t *FragmentElementGroup) GetElementListLen() int {
	return t.elemList.Len()
}

func (t *FragmentElementGroup) Release() (clenListLen int) {
	clenListLen = t.cleanUpElementList()
	t.ptrFinalElem = nil
	t.elemList = nil
	return
}

func (t *FragmentElementGroup) cleanUpElementList() int {
	var elems []*list.Element
	for e := t.elemList.Front(); e != nil; e = e.Next() {
		elems = append(elems, e)
	}
	for _, elem := range elems {
		t.elemList.Remove(elem)
		if elem.Value == nil {
			continue
		}

		fragElem, ok := elem.Value.(*FragmentElement)
		if !ok {
			continue
		}
		RecycleFragmentElement(fragElem)
	}

	return len(elems)
}

func (t *FragmentElementGroup) GetCreateTimestamp() int64 {
	return t.createTp
}

func (t *FragmentElementGroup) GetAllElementsPayloadLen() uint16 {
	var totalPayloadLen int
	for e := t.elemList.Front(); e != nil; e = e.Next() {
		elem := e.Value.(*FragmentElement)
		totalPayloadLen += elem.PayloadBuf.Len()
	}
	return uint16(totalPayloadLen)
}