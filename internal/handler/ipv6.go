package handler

import (
	"container/list"
	"encoding/binary"
	"fmt"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/internal/common"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type IPV6Handler struct{}

func (t *IPV6Handler) FastDetect(detectInfo *def.DetectionInfo) (retErr error, retErrType def.ErrResultType) {
	if len(detectInfo.EthPayload) <= def.IPV6HdrLen {
		retErr = fmt.Errorf("the IPV6 packet header length less than %d", def.IPV6HdrLen)
		retErrType = def.ErrResultIPV6HdrLenInsufficient
		return
	}

	buf := detectInfo.EthPayload
	buf = buf[def.IPVersionLen+def.IPV6TrafficClassFlowLabelLen+def.IPV6PayloadLen:]
	ipProtocol := layers.IPProtocol(buf[0])
	if ipProtocol != layers.IPProtocolIPv6Fragment {
		detectInfo.IPProtocol = ipProtocol
		return
	}

	detectInfo.FragType = def.IPV6FragType
	detectInfo.TrafficClass = uint8((binary.BigEndian.Uint16(detectInfo.EthPayload[0:2]) >> 4) & 0x00FF)
	detectInfo.FlowLabel = binary.BigEndian.Uint32(detectInfo.EthPayload[0:4]) & 0x000FFFFF

	buf = buf[def.IPV6NextHeaderLen+def.IPV6HopLimitLen:]
	detectInfo.SrcIP = buf[:def.IPV6SrcAddrLen]
	buf = buf[def.IPV6SrcAddrLen:]

	detectInfo.DstIP = buf[:def.IPV6DstAddrLen]
	buf = buf[def.IPV6DstAddrLen:]

	if len(buf) <= def.IPV6FragmentHdrLen {
		retErr = fmt.Errorf("the IPV6 packet fragment header length less than %d", def.IPV6FragmentHdrLen)
		retErrType = def.ErrResultIPV6FragHdrLenInsufficient
		return
	}

	fragHdrBuf := buf[:def.IPV6FragmentHdrLen]
	detectInfo.IPProtocol = layers.IPProtocol(fragHdrBuf[0])
	fragHdrBuf = fragHdrBuf[def.IPV6FragmentNextHeaderLen+def.IPV6FragmentReservedOctetLen:]
	detectInfo.FragOffset = binary.BigEndian.Uint16(fragHdrBuf) >> 3
	detectInfo.MoreFrags = (fragHdrBuf[1] & 0x1) != 0

	fragHdrBuf = fragHdrBuf[def.IPV6FlagsFlagsLen:]
	detectInfo.Identification = binary.BigEndian.Uint32(fragHdrBuf)
	detectInfo.IPPayload = buf[def.IPV6FragmentHdrLen:]

	return
}

func (t *IPV6Handler) Collect(fragElem *common.FragElement, fragElemGroup *common.FragElementGroup) (error, def.ErrResultType) {
	return collectFragElement(fragElem, fragElemGroup)
}

func (t *IPV6Handler) Reassembly(fragElemGroup *common.FragElementGroup,
	sharedLayers *common.SharedLayers) (gopacket.Packet, error, def.ErrResultType) {

	finalElem := fragElemGroup.GetFinalElement()
	payloadLen := fragElemGroup.GetAllElementsPayloadLen()

	sharedLayers.EthFrame.SrcMAC = finalElem.SrcMAC
	sharedLayers.EthFrame.DstMAC = finalElem.DstMAC
	sharedLayers.EthFrame.EthernetType = layers.EthernetTypeIPv6

	sharedLayers.IPV6.Length = payloadLen
	sharedLayers.IPV6.TrafficClass = finalElem.TrafficClass
	sharedLayers.IPV6.FlowLabel = finalElem.FlowLabel
	sharedLayers.IPV6.NextHeader = finalElem.IPProtocol
	sharedLayers.IPV6.SrcIP = finalElem.SrcIP
	sharedLayers.IPV6.DstIP = finalElem.DstIP

	fullPktBuff := sharedLayers.FullIPV6Buff
	if err := gopacket.SerializeLayers(fullPktBuff, defaultSerializeOptions,
		&sharedLayers.EthFrame, &sharedLayers.IPV6); err != nil {
		return nil, err, def.ErrResultSerializeLayers
	}

	freeLen := len(fullPktBuff.Bytes()) - def.EthIPV6HdrLen
	_, appendErr := fullPktBuff.AppendBytes(int(payloadLen) - freeLen)
	if appendErr != nil {
		return nil, appendErr, def.ErrResultFullPacketBufAppendBytes
	}

	payloadSpace := fullPktBuff.Bytes()[def.EthIPV6HdrLen:]
	fragElemGroup.IterElementList(func(elem *list.Element) bool {
		fragElem := elem.Value.(*common.FragElement)
		fragPayloadLen := fragElem.PayloadBuf.Len()
		if fragPayloadLen <= 0 {
			// todo
			return true
		}

		copy(payloadSpace, fragElem.PayloadBuf.Bytes())
		payloadSpace = payloadSpace[fragPayloadLen:]
		return true
	})

	retPkt := gopacket.NewPacket(fullPktBuff.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
	if retPkt.ErrorLayer() != nil {
		return nil, retPkt.ErrorLayer().Error(), def.ErrResultTypeNewPacket
	}
	return retPkt, nil, def.NonErrResultType
}
