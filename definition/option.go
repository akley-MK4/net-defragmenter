package definition

import (
	"time"
)

func NewOption(fns ...func(opt *Option)) *Option {
	opt := &Option{}
	for _, f := range fns {
		f(opt)
	}

	return opt
}

type IPV6WorkerOption struct {
	QueueLen int
	Interval time.Duration
	RWNum    int
}

type CollectorOption struct {
	MaxCollectorsNum            uint32
	MaxChannelCap               uint32
	MaxFragGroupMapLength       uint32
	MaxFullPktQueueLen          uint32
	MaxFragGroupDurationSeconds int64
	EnableSyncReassembly        bool
}

type StatsOption struct {
	Enable bool
}

type Option struct {
	StatsOption       StatsOption
	PickFragmentTypes []FragType
	CollectorOption   CollectorOption
}
