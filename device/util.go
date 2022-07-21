package device

import (
	"strconv"
)

//decoding the byte to IP address
func decodeIP(addr []byte) string {
	var outputIP string
	for i, val := range addr {
		if len(addr)-1 != i {
			outputIP += strconv.Itoa(int(byte(val))) + "."
		} else {
			outputIP += strconv.Itoa(int(byte(val)))
		}
	}
	return outputIP
}

//decoding the protocol
func decodeProtocol(protocol []byte) string {
	// ICMP will have - 1 , TCP - 6 , UDP - 17
	for _, proto := range protocol {
		if int(byte(proto)) == 1 {
			return "icmp"
		} else if int(byte(proto)) == 6 {
			return "tcp"
		} else if int(byte(proto)) == 17 {
			return "udp"
		}
	}
	return ""
}
