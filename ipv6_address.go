package geoip

import (
  "strings"
  "strconv"
  "errors"
  "regexp"
)

const (
  IPV4_PARTS = 4
  IPV6_BYTES = 16
)

type Ipv6Address [IPV6_BYTES]byte

var IPV4_SUFFIX_PATTERN *regexp.Regexp

func init() {
  var err error
  if IPV4_SUFFIX_PATTERN, err = regexp.Compile(":(\\d{1,3}\\.)+\\d{1,3}$"); err != nil {
    println(err)
  }
}

func parseIpv4Suffix(ipv4Suffix *string) (*[IPV4_PARTS]byte, *error) {
  var ipv4Address [IPV4_PARTS]byte

  ipv4SuffixTrimmed := (*ipv4Suffix)[1:]
  ipv4Parts := strings.Split(ipv4SuffixTrimmed, ".")
  if len(ipv4Parts) != IPV4_PARTS {
    error := errors.New("IPv4 segment can constrain only four parts")
    return nil, &error
  }
  for index, ipv4Part := range ipv4Parts {
    if ipv4Octet, ipv4Err := strconv.ParseUint(ipv4Part, 10, 8); ipv4Err != nil {
      return nil, &ipv4Err
    } else {
      ipv4Address[index] = byte(ipv4Octet)
    }
  }

  return &ipv4Address, nil
}

func parseIpv6Part(ipv6Part *string) (*[]byte, *error) {
  if *ipv6Part == "" {
    emtpyOctets := make([]byte, 0)
    return &emtpyOctets, nil
  }
  parts := strings.Split(*ipv6Part, ":")
  octets := make([]byte, len(parts) * 2)
  for index, part := range parts {
    octet, err := strconv.ParseUint(part, 16, 16)
    if nil != err {
      return nil, &err
    }
    firstOctet := byte(octet & 0xFF)
    secondOctet := byte(octet >> 8)
    octets[index * 2 + 1] = firstOctet
    octets[index * 2] = secondOctet
  }
  return &octets, nil
}

func ParseIpv6Address(ipStr string) (*Ipv6Address, *error) {
  var ipAddress Ipv6Address

  ipv4Suffix := IPV4_SUFFIX_PATTERN.FindString(ipStr)
  hasIpv4Suffix := ipv4Suffix != ""
  if hasIpv4Suffix {
    ipv4Octets, err := parseIpv4Suffix(&ipv4Suffix)
    if err != nil {
      return nil, err
    }
    rightOffset := IPV6_BYTES - IPV4_PARTS
    copy(ipAddress[rightOffset:], (*ipv4Octets)[:])
    ipStr = ipStr[:len(ipStr) - len(ipv4Suffix)]
  }

  ipParts := strings.Split(ipStr, "::")
  if len(ipParts) > 2 {
    error := errors.New("IPv6 address can constrain not more one '::'.")
    return nil, &error
  }

  ipv6Octets, err := parseIpv6Part(&ipParts[0])
  if nil != err {
    return nil, err
  } else {
    copy(ipAddress[:len(*ipv6Octets)], *ipv6Octets)
  }

  if len(ipParts) == 2 {
    ipv6Octets, err := parseIpv6Part(&ipParts[1])
    if nil != err {
      return nil, err
    } else {
      endOffset := IPV6_BYTES
      if hasIpv4Suffix {
        endOffset -= IPV4_PARTS
      }
      startOffset := endOffset - len(*ipv6Octets)
      copy(ipAddress[startOffset:endOffset], *ipv6Octets)
    }
  }

  return &ipAddress, nil
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
