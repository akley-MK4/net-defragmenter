package handler

import (
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/internal/common"
	"github.com/google/gopacket"
)

type IHandler interface {
	FastDetect(detectInfo *def.DetectionInfo) (error, def.ErrResultType)
	Collect(fragElem *common.FragElement, fragElemGroup *common.FragElementGroup) (error, def.ErrResultType)
	Reassembly(fragElemGroup *common.FragElementGroup, sharedLayers *common.SharedLayers) (gopacket.Packet, error, def.ErrResultType)
}

var (
	defaultSerializeOptions = gopacket.SerializeOptions{FixLengths: false, ComputeChecksums: false}
)
