package definition

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type FragType int8

const (
	NonFragType FragType = iota

	IPV4FragType
	IPV6FragType
	PFCPFragType

	MaxInvalidFragType
)

type FragGroupID string

var (
	layerEnumMapping = map[interface{}]FragType{
		layers.EthernetTypeIPv4: IPV4FragType,
		layers.EthernetTypeIPv6: IPV6FragType,
	}
)

type OnDetectCompleted func(fragType FragType, fragGroupID FragGroupID)

type FullPacket struct {
	InterfaceId InterfaceId
	FragGroupID FragGroupID
	Pkt         gopacket.Packet
	FragElemLen int
}

func (t *FullPacket) GetInterfaceId() InterfaceId {
	return t.InterfaceId
}

func (t *FullPacket) GetFragGroupID() FragGroupID {
	return t.FragGroupID
}

func (t *FullPacket) GetPacket() gopacket.Packet {
	return t.Pkt
}

type DetectionInfo struct {
	SrcMAC, DstMAC []byte
	EthType        layers.EthernetType
	EthPayload     []byte

	SrcIP, DstIP   net.IP
	IPProtocol     layers.IPProtocol
	FragType       FragType
	FragOffset     uint16
	MoreFrags      bool
	Identification uint32
	IPPayload      []byte

	InterfaceId InterfaceId
	FragGroupId FragGroupID
}

func (t *DetectionInfo) Reset() {
	t.SrcMAC = nil
	t.DstMAC = nil
	t.EthPayload = nil
	t.SrcIP = nil
	t.DstIP = nil
	t.IPPayload = nil
	t.FragGroupId = ""
}

type ReplyParse struct {
	SrcIP          string
	DstIP          string
	Protocol       interface{}
	Identification uint32
}

type LayerHeaders struct {
	Eth  *layers.Ethernet
	IPV4 *layers.IPv4
	IPV6 *layers.IPv6
}
