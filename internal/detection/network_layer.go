package detection

import (
	"fmt"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/internal/handler"
	"github.com/akley-MK4/net-defragmenter/stats"
	"github.com/google/gopacket/layers"
)

func (t *Detector) detectNetworkLayer(detectInfo *def.DetectionInfo) error {
	var mappingFragType def.FragType
	switch detectInfo.EthType {
	case layers.EthernetTypeIPv4:
		mappingFragType = def.IPV4FragType
		break
	case layers.EthernetTypeIPv6:
		mappingFragType = def.IPV6FragType
		break
	default:
		return nil
	}

	if !t.pickFragTypeSet[mappingFragType] {
		return nil
	}

	statsHandler := stats.GetDetectionStatsHandler()
	hd := handler.GetHandler(mappingFragType)
	if hd == nil {
		statsHandler.AddTotalNoNetworkLayerHandlerErrNum(1)
		return fmt.Errorf("handler with fragment type %v dose not exists", mappingFragType)
	}

	detectErr, detectErrType := hd.FastDetect(detectInfo)
	if detectErr != nil {
		statsHandler.AddErrHandlerFastDetectStats(1, detectErrType)
		return detectErr
	}

	return nil
}
