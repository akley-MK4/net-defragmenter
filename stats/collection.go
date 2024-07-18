package stats

import (
	def "github.com/akley-MK4/net-defragmenter/definition"
	"sync/atomic"
)

type CollectionStats struct {
	//TotalSuccessfulStartedCollectorsNum uint64
	//TotalFailedStartCollectorsNum       uint64

	TotalFailedDistributionMemberNum uint64 `json:",omitempty"`
	TotalNewFragElementsNum          uint64 `json:",omitempty"`
	TotalAllocatedFragElementsNum    uint64 `json:",omitempty"`
	TotalRecycledFragElementsNum     uint64 `json:",omitempty"`

	TotalAcceptedFragElementsNum   uint64 `json:",omitempty"`
	TotalNotFoundHandlersNum       uint64 `json:",omitempty"`
	TotalFragMapReachedLenLimitNum uint64 `json:",omitempty"`

	TotalNewFragElementGroupsNum             uint64 `json:",omitempty"`
	TotalReleasedFragElementGroupsNum        uint64 `json:",omitempty"`
	TotalReleasedExpiredFragElementGroupsNum uint64 `json:",omitempty"`

	ErrorHandlerCollectStats         ErrorStats `json:",omitempty"`
	TotalSuccessfulCollectedFragsNum uint64     `json:",omitempty"`

	TotalReassemblyNoDelFragGroupsNum uint64     `json:",omitempty"`
	ErrHandlerReassemblyStats         ErrorStats `json:",omitempty"`
	TotalSuccessfulReassemblyFragsNum uint64     `json:",omitempty"`
	TotalPushedFullPacketsNum         uint64     `json:",omitempty"`
	TotalReleasedFullPacketsNum       uint64     `json:",omitempty"`
	TotalForceReleasedFullPacketsNum  uint64     `json:",omitempty"`
	TotalPoppedFullPacketsNum         uint64     `json:",omitempty"`
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

func (t *CollectionStatsHandler) IncTotalFragMapReachedLenLimitNum() {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalFragMapReachedLenLimitNum, 1)
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
