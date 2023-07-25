package detection

import (
	"errors"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/akley-MK4/net-defragmenter/libstats"
	"github.com/google/gopacket/layers"
)

type Detector struct {
	pickFragTypeSet map[def.FragmentType]bool
}

func NewDetector(pickFragTypes []def.FragmentType) (*Detector, error) {
	pickFragTypeSet := make(map[def.FragmentType]bool)
	for _, fragTpy := range pickFragTypes {
		if fragTpy <= def.InvalidFragType || fragTpy >= def.MaxInvalidFragType {
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
	if err := detectEthernetLayer(pktBuf, detectInfo); err != nil {
		libstats.AddTotalDetectLinkLayerErrNum(1)
		return err
	}
	if detectInfo.EthType == layers.EthernetTypeLLC {
		return nil
	}

	if err := t.detectNetworkLayer(detectInfo); err != nil {
		return err
	}

	if detectInfo.FragType == def.IPV4FragType || detectInfo.FragType == def.IPV6FragType {
		libstats.AddTotalDetectPassedNum(1)
		return nil
	}

	// Application layer not currently supported
	if true {
		libstats.AddTotalPickFragTypeNotExistsNum(1)
		return nil
	}

	if err := detectApplicationLayer(detectInfo); err != nil {
		libstats.AddTotalDetectAppLayerErrNum(1)
		return err
	}

	return nil
}
