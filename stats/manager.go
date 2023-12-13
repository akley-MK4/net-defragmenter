package stats

import "sync/atomic"

type ManagerStats struct {
	TotalFailedCallAsyncProcessNum uint64
	TotalCrashedAsyncProcessNum    uint64
}

var (
	managerStatsHandler = &ManagerStatsHandler{}
)

func GetManagerStatsHandler() *ManagerStatsHandler {
	return managerStatsHandler
}

type ManagerStatsHandler struct {
	stats ManagerStats
}

func (t *ManagerStatsHandler) getStats() ManagerStats {
	return t.stats
}

func (t *ManagerStatsHandler) AddTotalCrashedAsyncProcessNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalCrashedAsyncProcessNum, delta)
}

func (t *ManagerStatsHandler) AddTotalFailedCallAsyncProcessNum(delta uint64) {
	if !enabledStats {
		return
	}
	atomic.AddUint64(&t.stats.TotalFailedCallAsyncProcessNum, delta)
}
