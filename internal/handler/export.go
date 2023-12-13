package handler

import def "github.com/akley-MK4/net-defragmenter/definition"

var (
	handlerMap = map[def.FragType]IHandler{
		def.IPV6FragType: &IPV6Handler{},
		def.IPV4FragType: &IPV4Handler{},
	}
)

func GetHandler(fragType def.FragType) IHandler {
	return handlerMap[fragType]
}
