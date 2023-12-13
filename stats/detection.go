package stats

import (
	def "github.com/akley-MK4/net-defragmenter/definition"
	"sync/atomic"
)

type DetectionStats struct {
	TotalReceivedDetectPacketsNum uint64 `json:"TotalReceivedDetectPacketsNum,omitempty"`
	//TotalNewDetectInfoNum             uint64
	//TotalReleaseDetectInfoNum         uint64
	TotalFailedDetectEthernetLayerNum uint64     `json:"TotalFailedDetectEthernetLayerNum,omitempty"`
	TotalFailedDetectNetworkLayerNum  uint64     `json:"TotalFailedDetectNetworkLayerNum,omitempty"`
	TotalNoNetworkLayerHandlerErrNum  uint64     `json:"TotalNoNetworkLayerHandlerErrNum,omitempty"`
	ErrHandlerFastDetectStats         ErrorStats `json:"ErrHandlerFastDetectStats,omitempty"`

	TotalSuccessfulDetectedFragsNum uint64 `json:"TotalSuccessfulDetectedFragsNum,omitempty"`
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

func (t *DetectionStatsHandler) AddTotalReceivedDetectPacketsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalReceivedDetectPacketsNum, delta)
}

//func (t *DetectionStatsHandler) AddTotalNewDetectInfoNum(delta uint64) {
//	if !enabledStats {
//		return
//	}
//	atomic.AddUint64(&t.stats.TotalNewDetectInfoNum, delta)
//}
//
//func (t *DetectionStatsHandler) AddTotalReleaseDetectInfoNum(delta uint64) {
//	if !enabledStats {
//		return
//	}
//	atomic.AddUint64(&t.stats.TotalReleaseDetectInfoNum, delta)
//}

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
