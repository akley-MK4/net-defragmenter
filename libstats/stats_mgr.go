package libstats

import (
	def "github.com/akley-MK4/net-defragmenter/definition"
)

var (
	mgr = &StatsMgr{}
)

type StatsMgr struct {
	Enabled    bool
	Detection  DetectionStats
	Collection CollectionStats
}

func InitStatsMgr(opt def.StatsOption) {
	mgr.Enabled = opt.Enable
}

func GetStatsMgr() StatsMgr {
	return *mgr
}