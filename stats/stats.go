package stats

var (
	enabledStats = false
)

func EnableStats() {
	enabledStats = true
}

func DisableStats() {
	enabledStats = false
}

type Stats struct {
	//Manager    ManagerStats
	Detection  DetectionStats
	Collection CollectionStats
}

func GetStats() (retStats Stats) {
	//retStats.Manager = GetManagerStatsHandler().getStats()
	retStats.Detection = GetDetectionStatsHandler().getStats()
	retStats.Collection = GetCollectionStatsHandler().getStats()
	return
}
