package main

import (
	"os"
)

func main() {
	//example.LaunchDemoWithPcapReply("E:\\pcaps/small_ipv4frag.pcap")
	//example.LaunchDemoWithPcapReply("E:\\pcaps/ipv6frag.pcap")
	launchDemoWithPcapReply("E:\\pcaps/ipv4frag_10k.pcap")
	os.Exit(0)
}
