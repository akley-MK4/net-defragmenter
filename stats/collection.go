package stats

import (
	def "github.com/akley-MK4/net-defragmenter/definition"
	"sync/atomic"
)

type CollectionStats struct {
	//TotalSuccessfulStartedCollectorsNum uint64
	//TotalFailedStartCollectorsNum       uint64

	TotalFailedDistributionMemberNum uint64 `json:"TotalFailedDistributionMemberNum,omitempty"`
	TotalNewFragElementsNum          uint64 `json:"TotalNewFragElementsNum,omitempty"`
	TotalAllocatedFragElementsNum    uint64 `json:"TotalAllocatedFragElementsNum,omitempty"`
	TotalRecycledFragElementsNum     uint64 `json:"TotalRecycledFragElementsNum,omitempty"`

	TotalAcceptedFragElementsNum uint64 `json:"TotalAcceptedFragElementsNum,omitempty"`
	TotalNotFoundHandlersNum     uint64 `json:"TotalNotFoundHandlersNum,omitempty"`

	TotalNewFragElementGroupsNum             uint64 `json:"TotalNewFragElementGroupsNum,omitempty"`
	TotalReleasedFragElementGroupsNum        uint64 `json:"TotalReleasedFragElementGroupsNum,omitempty"`
	TotalReleasedExpiredFragElementGroupsNum uint64 `json:"TotalReleasedExpiredFragElementGroupsNum,omitempty"`

	ErrorHandlerCollectStats         ErrorStats `json:"ErrorHandlerCollectStats,omitempty"`
	TotalSuccessfulCollectedFragsNum uint64     `json:"TotalSuccessfulCollectedFragsNum,omitempty"`

	TotalReassemblyNoDelFragGroupsNum uint64     `json:"TotalReassemblyNoDelFragGroupsNum,omitempty"`
	ErrHandlerReassemblyStats         ErrorStats `json:"ErrHandlerReassemblyStats,omitempty"`
	TotalSuccessfulReassemblyFragsNum uint64     `json:"TotalSuccessfulReassemblyFragsNum,omitempty"`
	TotalPushedFullPacketsNum         uint64     `json:"TotalPushedFullPacketsNum,omitempty"`
	TotalReleasedFullPacketsNum       uint64     `json:"TotalReleasedFullPacketsNum,omitempty"`
	TotalForceReleasedFullPacketsNum  uint64     `json:"TotalForceReleasedFullPacketsNum,omitempty"`
	TotalPoppedFullPacketsNum         uint64     `json:"TotalPoppedFullPacketsNum,omitempty"`
}

var (
	collectionStatsHandler = &CollectionStatsHandler{}
)

func GetCollectionStatsHandler() *CollectionStatsHandler {
	return collectionStatsHandler
}

type CollectionStatsHandler struct {
	stats CollectionStats
}

func (t *CollectionStatsHandler) getStats() CollectionStats {
	return t.stats
}

//func (t *CollectionStatsHandler) AddTotalSuccessfulStartedCollectorsNum(delta uint64) {
//	if !enabledStats {
//		return
//	}
//	atomic.AddUint64(&t.stats.TotalSuccessfulStartedCollectorsNum, delta)
//}
//
//func (t *CollectionStatsHandler) AddTotalFailedStartCollectorsNum(delta uint64) {
//	if !enabledStats {
//		return
//	}
//	atomic.AddUint64(&t.stats.TotalFailedStartCollectorsNum, delta)
//}

func (t *CollectionStatsHandler) AddTotalFailedDistributionMemberNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalFailedDistributionMemberNum, delta)
}

func (t *CollectionStatsHandler) AddTotalNewFragElementsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalNewFragElementsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalAllocatedFragElementsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalAllocatedFragElementsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalRecycledFragElementsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalRecycledFragElementsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalAcceptedFragElementsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalAcceptedFragElementsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalNotFoundHandlersNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalNotFoundHandlersNum, delta)
}

func (t *CollectionStatsHandler) AddTotalNewFragElementGroupsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalNewFragElementGroupsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalReleasedFragElementGroupsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalReleasedFragElementGroupsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalReleasedExpiredFragElementGroupsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalReleasedExpiredFragElementGroupsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalErrCollectNum(delta uint64, errType def.ErrResultType) {
	t.stats.ErrorHandlerCollectStats.AddTotalNum(delta, errType)
}

func (t *CollectionStatsHandler) AddTotalErrReassemblyNum(delta uint64, errType def.ErrResultType) {
	t.stats.ErrHandlerReassemblyStats.AddTotalNum(delta, errType)
}

func (t *CollectionStatsHandler) AddTotalReassemblyNoDelFragGroupsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalReassemblyNoDelFragGroupsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalSuccessfulReassemblyFragsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalSuccessfulReassemblyFragsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalPushedFullPacketsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalPushedFullPacketsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalForceReleasedFullPacketsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalForceReleasedFullPacketsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalReleasedFullPacketsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalReleasedFullPacketsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalPoppedFullPacketsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalPoppedFullPacketsNum, delta)
}

func (t *CollectionStatsHandler) AddTotalSuccessfulCollectedFragsNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalSuccessfulCollectedFragsNum, delta)
}
