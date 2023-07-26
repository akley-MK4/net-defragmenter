# net-defragmenter
## Introduction
This is a network fragmentation reassembly library that now supports fragments of IPV4 and IPV6 types.  
By fast detecting data packets (do not copy memory, do not use 'gopacket.NewPacket'), balancing group collection, 
concurrent defragmentation, object reuse, and other operations,  
a good level of performance can be achieved.
## Directory structure description
### definition
This directory mainly contains some public macro definitions and some public structures.
### example
This directory mainly contains some demonstration examples.
### fragadapter
This is a demonstration of the adapter, which demonstrates how to use net-delimiter library. This adapter is just a demonstration example, you don't need to follow this demonstration example to use the library.
### internal
This directory mainly contains internal implementations that cannot be referenced.
### libstats
Stats data related
### manager
The instance of the library is related, and this instance is exposed to users for use.

## Simple Test Report Reference
### Performance
In a CPU @ 3.20GHz 3.19 GHz 4 threads environment, processing a pcap file with 102146 packets took 473 milliseconds, and restructuring 61277 fragments of the file took 458 milliseconds.

### Stats
The description will be supplemented later.
```json
{
	"Enabled": true,
	"Detection": {
		"TotalReceivedPktNum": 102146,
		"TotalLinkLayerErrNum": 0,
		"TotalHandleNilErrNum": 0,
		"ErrStats": {
			"TotalNewPacketErrNum": 0,
			"TotalSerializeLayersErrNum": 0,
			"TotalFullPacketBufAppendBytesErrNum": 0,
			"TotalIPV4HdrLenInsufficient": 0,
			"TotalIPv6NetWorkerLayerNilNum": 0,
			"TotalIPV6HdrLenInsufficientNum": 0,
			"TotalIPV6FragHdrLenInsufficientNum": 0,
			"TotalUnhandledErrNum": 0
		},
		"TotalPickFragTypeNotExistsNum": 39343,
		"TotalFilterAppLayerErrNum": 0,
		"TotalDetectPassedNum": 61277
	},
	"Collection": {
		"TotalDistributeFragmentFailureNum": 0,
		"TotalNewFragmentElementNum": 61277,
		"TotalAllocateFragmentElementNum": 1197,
		"TotalRecycleFragmentElementNum": 61277,
		"TotalAcceptFragmentElementNum": 61277,
		"TotalHandleNilErrNum": 0,
		"CollectErrStats": {
			"TotalNewPacketErrNum": 0,
			"TotalSerializeLayersErrNum": 0,
			"TotalFullPacketBufAppendBytesErrNum": 0,
			"TotalIPV4HdrLenInsufficient": 0,
			"TotalIPv6NetWorkerLayerNilNum": 0,
			"TotalIPV6HdrLenInsufficientNum": 0,
			"TotalIPV6FragHdrLenInsufficientNum": 0,
			"TotalUnhandledErrNum": 0
		},
		"TotalAcceptFragSuccessfulNum": 61277,
		"ReassemblyErrStats": {
			"TotalNewPacketErrNum": 0,
			"TotalSerializeLayersErrNum": 0,
			"TotalFullPacketBufAppendBytesErrNum": 0,
			"TotalIPV4HdrLenInsufficient": 0,
			"TotalIPv6NetWorkerLayerNilNum": 0,
			"TotalIPV6HdrLenInsufficientNum": 0,
			"TotalIPV6FragHdrLenInsufficientNum": 0,
			"TotalUnhandledErrNum": 0
		},
		"TotalNewFragGroupNum": 16659,
		"TotalDelFragGroupNotExistNum": 0,
		"TotalReleaseFragGroupThReassemblyNum": 16659,
		"TotalReleaseFragGroupThExpiredNum": 0,
		"TotalReassemblyFragNum": 61277,
		"TotalPushFullPktNum": 16659,
		"TotalPopFullPktNum": 16659,
		"TotalReleaseFullPktNum": 0
	}
}
```

### Memory Status
The following comparison shows that there is no memory leak.
#### Initial memory state
```json
{
  "Alloc": "Alloc MBs=1, KBs=1160, Bs=1188464",
  "TotalAlloc": "TotalAlloc MBs=1",
  "Sys": "TotalAlloc MBs=6",
  "Mallocs": 2853,
  "Frees": 110,
  "HeapAlloc": "HeapAlloc MBs=1",
  "HeapSys": "HeapSys MBs=3",
  "HeapInuse": "HeapInuse MBs=2",
  "HeapReleased": "HeapReleased MBs=1",
  "NextGC": "NextGC MBs=5",
  "GCCPUFraction": 0,
  "NumForcedGC": 0
}
```
#### Peak memory status
```json
{
	"Alloc": "Alloc MBs=133, KBs=136759, Bs=140041912",
	"TotalAlloc": "TotalAlloc MBs=364",
	"Sys": "TotalAlloc MBs=150",
	"Mallocs": 1717504,
	"Frees": 1283353,
	"HeapAlloc": "HeapAlloc MBs=133",
	"HeapSys": "HeapSys MBs=139",
	"HeapInuse": "HeapInuse MBs=138",
	"HeapReleased": "HeapReleased MBs=1",
	"NextGC": "NextGC MBs=153",
	"GCCPUFraction": 0.06129849919264683,
	"NumForcedGC": 0
}
```
#### Memory status after garbage collection
```json
{
	"Alloc": "Alloc MBs=1, KBs=1283, Bs=1314176",
	"TotalAlloc": "TotalAlloc MBs=365",
	"Sys": "TotalAlloc MBs=154",
	"Mallocs": 1717664,
	"Frees": 1715580,
	"HeapAlloc": "HeapAlloc MBs=1",
	"HeapSys": "HeapSys MBs=143",
	"HeapInuse": "HeapInuse MBs=2",
	"HeapReleased": "HeapReleased MBs=132",
	"NextGC": "NextGC MBs=8",
	"GCCPUFraction": 0.0009613584960369672,
	"NumForcedGC": 3
}
```



