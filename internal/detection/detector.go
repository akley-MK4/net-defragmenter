package detection

import (
	"errors"
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

func (t *Detector) FastDetect(pktBuf []byte, detectInfo *def.DetectionInfo) error {
	statsHandler := stats.GetDetectionStatsHandler()
	statsHandler.AddTotalDetectedPacketsNum(1)

	if err := detectEthernetLayer(pktBuf, detectInfo); err != nil {
		statsHandler.AddTotalFailedDetectEthernetLayerNum(1)
		return err
	}
	if detectInfo.EthType == layers.EthernetTypeLLC {
		return nil
	}

	if err := t.detectNetworkLayer(detectInfo); err != nil {
		statsHandler.AddTotalFailedDetectNetworkLayerNum(1)
		return err
	}

	if detectInfo.FragType == def.IPV4FragType || detectInfo.FragType == def.IPV6FragType {
		stats.GetDetectionStatsHandler().AddTotalSuccessfulDetectedFragsNum(1)
		return nil
	}

	// todo Application layer not currently supported
	return nil
}
