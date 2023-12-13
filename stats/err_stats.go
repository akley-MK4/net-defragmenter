package stats

import (
	def "github.com/akley-MK4/net-defragmenter/definition"
	"sync/atomic"
)

type ErrorStats struct {
	TotalErrResultTypeNewPacketsNum             uint64 `json:"TotalErrResultTypeNewPacketsNum,omitempty"`
	TotalErrResultSerializeLayersNum            uint64 `json:"TotalErrResultSerializeLayersNum,omitempty"`
	TotalErrResultFullPacketBufAppendBytes      uint64 `json:"TotalErrResultFullPacketBufAppendBytes,omitempty"`
	TotalErrResultIPV4HdrLenInsufficientNum     uint64 `json:"TotalErrResultIPV4HdrLenInsufficientNum,omitempty"`
	TotalErrResultIPV6NetworkLayerNilNum        uint64 `json:"TotalErrResultIPV6NetworkLayerNilNum,omitempty"`
	TotalErrResultIPV6HdrLenInsufficientNum     uint64 `json:"TotalErrResultIPV6HdrLenInsufficientNum,omitempty"`
	TotalErrResultIPV6FragHdrLenInsufficientNum uint64 `json:"TotalErrResultIPV6FragHdrLenInsufficientNum,omitempty"`

	TotalUnknownErrNum uint64 `json:"TotalUnknownErrNum,omitempty"`
}

func (t *ErrorStats) AddTotalNum(delta uint64, errType def.ErrResultType) {
	if !enabledStats {
		return
	}

	switch errType {
	case def.ErrResultTypeNewPacket:
		atomic.AddUint64(&t.TotalErrResultTypeNewPacketsNum, delta)
		break
	case def.ErrResultSerializeLayers:
		atomic.AddUint64(&t.TotalErrResultSerializeLayersNum, delta)
		break
	case def.ErrResultFullPacketBufAppendBytes:
		atomic.AddUint64(&t.TotalErrResultFullPacketBufAppendBytes, delta)
		break
	case def.ErrResultIPV4HdrLenInsufficient:
		atomic.AddUint64(&t.TotalErrResultIPV4HdrLenInsufficientNum, delta)
		break
	case def.ErrResultIPV6NetworkLayerNil:
		atomic.AddUint64(&t.TotalErrResultIPV6NetworkLayerNilNum, delta)
		break
	case def.ErrResultIPV6HdrLenInsufficient:
		atomic.AddUint64(&t.TotalErrResultIPV6HdrLenInsufficientNum, delta)
		break
	case def.ErrResultIPV6FragHdrLenInsufficient:
		atomic.AddUint64(&t.TotalErrResultIPV6FragHdrLenInsufficientNum, delta)
		break
	default:
		atomic.AddUint64(&t.TotalUnknownErrNum, delta)
	}
}
