package stats

import (
	"sync/atomic"

	def "github.com/akley-MK4/net-defragmenter/definition"
)

type DetectionStats struct {
	TotalReceivedDetectPacketsNum uint64 `json:",omitempty"`
	//TotalNewDetectInfoNum             uint64
	//TotalReleaseDetectInfoNum         uint64
	TotalFailedDetectEthernetLayerNum uint64     `json:",omitempty"`
	TotalFailedDetectNetworkLayerNum  uint64     `json:",omitempty"`
	TotalNoNetworkLayerHandlerErrNum  uint64     `json:",omitempty"`
	ErrHandlerFastDetectStats         ErrorStats `json:",omitempty"`

	TotalDetectedFragsNum uint64 `json:",omitempty"`
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

func (t *DetectionStatsHandler) AddTotalDetectedFragsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalDetectedFragsNum, delta)
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
