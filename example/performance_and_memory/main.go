package main

import (
	"flag"
	def "github.com/akley-MK4/net-defragmenter/definition"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"runtime"
	"time"
)

// pcaps/ipv4frag_10k.pcap
func main() {
	pcapFilePath := flag.String("pcap_file_path", "", "pcap_file_path=")
	flag.Parse()

	availableNumCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(availableNumCPU)
	log.Printf("The current number of logical CPUs available for the process is %d\n", availableNumCPU)

	if err := initLibInstance(); err != nil {
		log.Printf("Failed to initialize LibInstance, %v\n", err)
		os.Exit(1)
	}
	log.Println("Successfully initialized LibInstance")
	if err := getLibInstance().Start(); err != nil {
		log.Printf("Failed to start LibInstance, %v\n", err)
		os.Exit(1)
	}
	log.Println("Successfully started LibInstance")

	if err := initInterfaceRecordMgr(time.Millisecond*100, 10000); err != nil {
		log.Printf("Failed to initialize InterfaceRecordMgr, %v\n", err)
		os.Exit(1)
	}
	log.Println("Successfully initialized InterfaceRecordMgr")
	if err := getInterfaceRecordMgr().start(); err != nil {
		log.Printf("Failed to start InterfaceRecordMgr, %v\n", err)
		os.Exit(1)
	}
	log.Println("Successfully started LibInstance")
	interfaceId := getInterfaceRecordMgr().RegisterInterfaceRecord("net1")

	initialMemInfo := collectMemoryStatus("Initial memory")

	log.Println("Start loading the pcap file")
	pcapHandle, openErr := pcap.OpenOffline(*pcapFilePath)
	if openErr != nil {
		log.Printf("pcap.OpenOffline failed, Path: %s, Err: %v\n", *pcapFilePath, openErr)
		os.Exit(1)
	}
	log.Println("Successfully loaded the pcap file")

	nowTime := time.Now()
	capSeconds, capNanoseconds := int64(nowTime.Second()), nowTime.Nanosecond()

	packetSource := gopacket.NewPacketSource(pcapHandle, pcapHandle.LinkType())
	for packet := range packetSource.Packets() {
		pktData := packet.Data()
		detectInfo := &def.DetectionInfo{}

		if err := getLibInstance().FastDetect(interfaceId, pktData, detectInfo); err != nil {
			log.Printf("Failed to fast detect packet, %v\n", err)
			continue
		}
		if detectInfo.FragGroupId == "" {
			continue
		}

		if err := getInterfaceRecordMgr().AssociateCapturedInfo(detectInfo, capSeconds, capNanoseconds); err != nil {
			log.Printf("Failed to associate cap info. %v\n", err)
			continue
		}
		if err := getLibInstance().Collect(detectInfo); err != nil {
			log.Printf("Failed to collect frag packet, %v\n", err)
			continue
		}
	}

	pcapHandle.Close()

	collectMemoryStatus("Current memory")

	dt := time.Second * 20
	log.Printf("Wait for %v to complete the reassembly of all fragments\n", dt)
	time.Sleep(dt)

	printStats()

	time.Sleep(time.Second * 1)
	runtime.GC()
	time.Sleep(time.Second * 3)
	runtime.GC()

	finalMemInfo := collectMemoryStatus("Final memory")
	printMemoryStatus(initialMemInfo)

	if (finalMemInfo.AllocMBs - initialMemInfo.AllocMBs) >= 1 {
		log.Println("There is a difference between the final memory and the initial memory size, please check if there is a memory leak")
	} else {
		log.Println("Successfully ran the example without generating any errors")
	}

	os.Exit(0)
}
