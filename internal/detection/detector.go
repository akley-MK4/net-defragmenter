package detection

import (
	"errors"
	"fmt"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/stats"
	"github.com/google/gopacket/layers"
)

type Detector struct {
	pickFragTypeSet map[def.FragType]bool
}

func NewDetector(pickFragTypes []def.FragType) (*Detector, error) {
	pickFragTypeSet := make(map[def.FragType]bool)
	for _, fragTpy := range pickFragTypes {
		if fragTpy <= def.NonFragType || fragTpy >= def.MaxInvalidFragType {
			continue
		}
		pickFragTypeSet[fragTpy] = true
	}

	if len(pickFragTypeSet) <= 0 {
		return nil, errors.New("no valid pick fragment type")
	}

	return &Detector{pickFragTypeSet: pickFragTypeSet}, nil
}

func (t *Detector) FastDetect(pktData []byte, replyDetectInfo *def.DetectionInfo) error {
	statsHandler := stats.GetDetectionStatsHandler()
	statsHandler.AddTotalReceivedDetectPacketsNum(1)

	if err := detectEthernetLayer(pktData, replyDetectInfo); err != nil {
		statsHandler.AddTotalFailedDetectEthernetLayerNum(1)
		return err
	}
	if replyDetectInfo.EthType == layers.EthernetTypeLLC {
		return nil
	}

	if err := t.detectNetworkLayer(replyDetectInfo); err != nil {
		statsHandler.AddTotalFailedDetectNetworkLayerNum(1)
		return err
	}

	if replyDetectInfo.FragType == def.IPV4FragType || replyDetectInfo.FragType == def.IPV6FragType {
		replyDetectInfo.FragGroupId = def.FragGroupID(fmt.Sprintf("%s-%s-%v-%d", replyDetectInfo.SrcIP.String(),
			replyDetectInfo.DstIP.String(), replyDetectInfo.IPProtocol, replyDetectInfo.Identification))
		stats.GetDetectionStatsHandler().AddTotalSuccessfulDetectedFragsNum(1)
		return nil
	}

	// todo Application layer not currently supported
	return nil
}
