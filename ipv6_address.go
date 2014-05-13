package geoipcsv

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	IPV6_BYTES          = 16
	IPV4_PARTS          = 4
	IPV4_PREFIX         = "::ffff:"
	IPV4_PART_SEPARATOR = "."
)

var (
	IPV4_PREFIX_LEN   = len(IPV4_PREFIX)
	IPV4_FIXED_PART   = []byte{0x00, 0x00, 0xFF, 0xFF}
	IPV4_FIXED_OFFSET = IPV6_BYTES - IPV4_PARTS - len(IPV4_FIXED_PART)
)

type Ipv6Address [IPV6_BYTES]byte

func parseIpv4(ipAddress string) (*Ipv6Address, *error) {
	ipv4Suffix := ipAddress[IPV4_PREFIX_LEN:]
	ipv4Parts := strings.Split(ipv4Suffix, IPV4_PART_SEPARATOR)
	if len(ipv4Parts) != IPV4_PARTS {
		error := fmt.Errorf("IPv4 segment can constrain only %d parts", IPV4_PARTS)
		return nil, &error
	}

	ipv4Address := new(Ipv6Address)
	copy(ipv4Address[IPV4_FIXED_OFFSET:], IPV4_FIXED_PART)

	for partIndex, ipv4Part := range ipv4Parts {
		ipv4Octet, ipv4Err := strconv.ParseUint(ipv4Part, 10, 8)
		if ipv4Err != nil {
			return nil, &ipv4Err
		} else {
			index := IPV6_BYTES - IPV4_PARTS + partIndex
			ipv4Address[index] = byte(ipv4Octet)
		}
	}

	return ipv4Address, nil
}

func parseIpv6(ipStr string) (*Ipv6Address, *error) {
	ipParts := strings.Split(ipStr, "::")

	if len(ipParts) > 2 {
		error := errors.New("IPv6 address can constrain not more one '::'.")
		return nil, &error
	}

	ipAddress := new(Ipv6Address)

	leftOctets, err := parseIpv6Part(&ipParts[0])
	if nil != err {
		return nil, err
	}
	leftOctetCount := len(*leftOctets)
	if leftOctetCount > IPV6_BYTES {
		err := fmt.Errorf("Ip address %s contains %d octets, but ipv6 can contain ony %d octects.",
			ipAddress, leftOctetCount, IPV6_BYTES)
		return nil, &err
	}
	copy(ipAddress[:], *leftOctets)

	if len(ipParts) == 2 {
		rightOctets, err := parseIpv6Part(&ipParts[1])
		if nil != err {
			return nil, err
		}
		rightOctetCount := len(*rightOctets)
		if rightOctetCount+leftOctetCount > IPV6_BYTES {
			err := fmt.Errorf("Ip address %s contains %d octets, but ipv6 can contain ony %d octects.",
				ipAddress, rightOctetCount+leftOctetCount, IPV6_BYTES)
			return nil, &err
		}
		offset := IPV6_BYTES - rightOctetCount
		copy(ipAddress[offset:], *rightOctets)
	}

	return ipAddress, nil
}

func parseIpv6Part(ipv6Part *string) (*[]byte, *error) {
	if *ipv6Part == "" {
		emtpyOctets := make([]byte, 0)
		return &emtpyOctets, nil
	}
	parts := strings.Split(*ipv6Part, ":")
	octets := make([]byte, len(parts)*2)
	for index, part := range parts {
		octet, err := strconv.ParseUint(part, 16, 16)
		if nil != err {
			return nil, &err
		}
		firstOctet := byte(octet & 0xFF)
		secondOctet := byte(octet >> 8)
		octets[index*2+1] = firstOctet
		octets[index*2] = secondOctet
	}
	return &octets, nil
}

func ParseIpv6Address(ipStr string) (*Ipv6Address, *error) {
	isIpv4 := strings.HasPrefix(ipStr, IPV4_PREFIX) && strings.Index(ipStr[IPV4_PREFIX_LEN:], IPV4_PART_SEPARATOR) > 0
	if isIpv4 {
		return parseIpv4(ipStr)
	} else {
		return parseIpv6(ipStr)
	}
}

func (address *Ipv6Address) Equal(otherAddress *Ipv6Address) bool {
	if otherAddress == nil {
		return false
	}
	for index, value := range address {
		if value != otherAddress[index] {
			return false
		}
	}
	return true
}

func (address *Ipv6Address) Compare(otherAddress *Ipv6Address) int {
	return bytes.Compare(address[:], otherAddress[:])
}
