package main

import (
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/implement"
)

var (
	libInstance *implement.Library
)

func getLibInstance() *implement.Library {
	return libInstance
}

func initLibInstance() (retErr error) {
	opt := def.NewOption(func(opt *def.Option) {
		//opt.CtrlApiServerOption.Enable = true
		opt.StatsOption.Enable = true

		opt.PickFragmentTypes = []def.FragType{def.IPV4FragType, def.IPV6FragType}

		opt.CollectorOption.MaxCollectorsNum = 30
		opt.CollectorOption.MaxChannelCap = 2000
		opt.CollectorOption.MaxFragGroupMapLength = 30
		opt.CollectorOption.MaxFullPktQueueLen = 10000
		opt.CollectorOption.MaxFragGroupDurationSeconds = 15
	})

	libInstance, retErr = implement.NewLibraryInstance(opt)
	if retErr != nil {
		return
	}

	return
}
