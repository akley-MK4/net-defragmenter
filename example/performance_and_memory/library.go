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

func initLibInstance(fns ...func(opt *def.Option)) (retErr error) {
	opt := def.NewOption(fns...)

	libInstance, retErr = implement.NewLibraryInstance(opt)
	if retErr != nil {
		return
	}

	return
}
