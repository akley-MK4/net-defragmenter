package handler

import (
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/internal/common"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type IPV4Handler struct{}

func (t *IPV4Handler) FastDetect(detectInfo *def.DetectionInfo) (retErr error, retErrType def.ErrResultType) {
	if len(detectInfo.EthPayload) <= def.IPV4HdrLen {
		retErr = fmt.Errorf("the IPV4 packet header length less than %d", def.IPV4HdrLen)
		retErrType = def.ErrResultIPV4HdrLenInsufficient
		return
	}

	buf := detectInfo.EthPayload
	buf = buf[def.IPVersionLen+def.IPV4DifferentiatedSvcFieldLen+def.IPV4TotalLengthFieldLen:]
	detectInfo.Identification = uint32(binary.BigEndian.Uint16(buf))
	buf = buf[def.IPV4IdentificationLen:]

	flagsFrags := binary.BigEndian.Uint16(buf)
	buf = buf[def.IPV4FlagsFlagsLen:]

	ipv4Flags := layers.IPv4Flag(flagsFrags >> 13)
	detectInfo.FragOffset = flagsFrags & 0x1FFF
	detectInfo.MoreFrags = (ipv4Flags & layers.IPv4MoreFragments) != 0

	buf = buf[def.IPV4TimeToLiveLen:]
	detectInfo.IPProtocol = layers.IPProtocol(buf[0])

	if !detectInfo.MoreFrags && detectInfo.FragOffset <= 0 {
		return
	}

	detectInfo.FragType = def.IPV4FragType

	buf = buf[def.IPV4ProtocolLen+def.IPV4HeaderChecksumLen:]
	detectInfo.SrcIP = buf[:def.IPV4SourceAddressLen]
	buf = buf[def.IPV4SourceAddressLen:]
	detectInfo.DstIP = buf[:def.IPV4DestinationAddressLen]

	detectInfo.IPPayload = buf[def.IPV4DestinationAddressLen:]
	return
}

func (t *IPV4Handler) Collect(fragElem *common.FragElement, fragElemGroup *common.FragElementGroup) (error, def.ErrResultType) {
	return collectFragElement(fragElem, fragElemGroup)
}

func (t *IPV4Handler) Reassembly(fragElemGroup *common.FragElementGroup,
	sharedLayers *common.SharedLayers) (gopacket.Packet, error, def.ErrResultType) {

	finalElem := fragElemGroup.GetFinalElement()
	payloadLen := fragElemGroup.GetAllElementsPayloadLen()

	sharedLayers.EthFrame.SrcMAC = finalElem.SrcMAC
	sharedLayers.EthFrame.DstMAC = finalElem.DstMAC
	sharedLayers.EthFrame.EthernetType = layers.EthernetTypeIPv4

	sharedLayers.IPV4.Id = uint16(finalElem.Identification)
	sharedLayers.IPV4.Length = payloadLen
	sharedLayers.IPV4.Protocol = finalElem.IPProtocol
	sharedLayers.IPV4.SrcIP = finalElem.SrcIP
	sharedLayers.IPV4.DstIP = finalElem.DstIP

	fullPktBuff := sharedLayers.FullIPV4Buff
	if err := gopacket.SerializeLayers(fullPktBuff, defaultSerializeOptions,
		&sharedLayers.EthFrame, &sharedLayers.IPV4); err != nil {
		return nil, err, def.ErrResultSerializeLayers
	}

	freeLen := len(fullPktBuff.Bytes()) - def.EthIPV4HdrLen
	_, appendErr := fullPktBuff.AppendBytes(int(payloadLen) - freeLen)
	if appendErr != nil {
		return nil, appendErr, def.ErrResultFullPacketBufAppendBytes
	}

	payloadSpace := fullPktBuff.Bytes()[def.EthIPV4HdrLen:]
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

func collectFragElement(fragElem *common.FragElement, fragElemGroup *common.FragElementGroup) (error, def.ErrResultType) {
	fragOffset := fragElem.FragOffset * def.FragOffsetMulNum
	if fragOffset >= fragElemGroup.GetHighest() {
		fragElem.Grouped = true
		fragElemGroup.PushElementToBack(fragElem)
	} else {
		fragElemGroup.IterElementList(func(elem *list.Element) bool {
			exitElem := elem.Value.(*common.FragElement)
			if exitElem.FragOffset == fragElem.FragOffset {
				// todo
				return false
			}
			if exitElem.FragOffset > fragElem.FragOffset {
				fragElem.Grouped = true
				fragElemGroup.InsertElementToBefore(fragElem, elem)
				return false
			}
			return true
		})
	}

	if !fragElem.Grouped {
		return errors.New("ungrouped fragments"), def.ErrResultUngroupedFrag
	}

	fragLength := uint16(fragElem.PayloadBuf.Len())
	if fragElemGroup.GetHighest() < fragOffset+fragLength {
		fragElemGroup.SetHighest(fragOffset + fragLength)
	}

	fragElemGroup.AddCurrentLen(fragLength)

	if !fragElem.MoreFrags {
		fragElemGroup.SetNextProtocol(fragElem.IPProtocol)
		fragElemGroup.SetFinalElement(fragElem)
	}

	return nil, def.NonErrResultType
}
