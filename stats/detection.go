package stats

import (
	def "github.com/akley-MK4/net-defragmenter/definition"
	"sync/atomic"
)

type DetectionStats struct {
	TotalDetectedPacketsNum uint64

	TotalFailedDetectEthernetLayerNum uint64
	TotalFailedDetectNetworkLayerNum  uint64
	TotalNoNetworkLayerHandlerErrNum  uint64
	ErrHandlerFastDetectStats         ErrorStats

	TotalSuccessfulDetectedFragsNum uint64
	//TotalLinkLayerErrNum          uint64
	//TotalHandleNilErrNum          uint64
	//ErrStats                      LayerPktErrStats
	//TotalPickFragTypeNotExistsNum uint64
	//TotalFilterAppLayerErrNum     uint64
	//TotalDetectPassedNum          uint64
}

var (
	detectionStatsHandler = &DetectionStatsHandler{}
)

func GetDetectionStatsHandler() *DetectionStatsHandler {
	return detectionStatsHandler
}

type DetectionStatsHandler struct {
	stats DetectionStats
}

func (t *DetectionStatsHandler) getStats() DetectionStats {
	return t.stats
}

func (t *DetectionStatsHandler) AddTotalDetectedPacketsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalDetectedPacketsNum, delta)
}

func (t *DetectionStatsHandler) AddTotalFailedDetectEthernetLayerNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalFailedDetectEthernetLayerNum, delta)
}

func (t *DetectionStatsHandler) AddTotalFailedDetectNetworkLayerNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalFailedDetectNetworkLayerNum, delta)
}

func (t *DetectionStatsHandler) AddTotalSuccessfulDetectedFragsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalSuccessfulDetectedFragsNum, delta)
}

func (t *DetectionStatsHandler) AddTotalNoNetworkLayerHandlerErrNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalNoNetworkLayerHandlerErrNum, delta)
}

func (t *DetectionStatsHandler) AddErrHandlerFastDetectStats(delta uint64, errType def.ErrResultType) {
	t.stats.ErrHandlerFastDetectStats.AddTotalNum(delta, errType)
}
