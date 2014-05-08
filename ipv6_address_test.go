package geoip

import (
  "testing"
)

func TestParseIpv6Address(t *testing.T) {
  var validValues = map[string]Ipv6Address {
    "1234:5678:90AB:CDEF:1234:5678:90AB:CDEF": 
      { 0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF, 0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF, },
    "::": 
      { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, },
    "::1":
      { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, },
    "0000:0000:0000:0000:0000:0000:0000:0000": 
      { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, },
    "FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF": 
      { 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, },
    "::FFFF:FFFF:FFFF:FFFF": 
      { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, },
    "::1234:5678:90AB:CDEF": 
      { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF, },
    "FEDC:BA09:8765:4321::": 
      { 0xFE, 0xDC, 0xBA, 0x09, 0x87, 0x65, 0x43, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, },
    "0123:4567::89AB:CDEF": 
      { 0x01, 0x23, 0x45, 0x67, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x89, 0xAB, 0xCD, 0xEF, },
    ":::1.2.3.4":
      { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, },
    "::FFFF:1.2.3.4":
      { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x01, 0x02, 0x03, 0x04, },
    "FFFF:::1.2.3.4":
      { 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, },
  }

  var invalidValues = []string {
    ":::",
    "::::",
    "::1::",
    "::1.2.3.4.5",
    ":",
    "GGGG::",
    "G:",
    "F:",
  }

  for validIpSource, validIp := range validValues {
    parsedIp, ipErr := ParseIpv6Address(validIpSource)
    if ipErr != nil {
      t.Errorf("Error parsing %s: %s", validIpSource, *ipErr)
    }
    if !validIp.Equal(parsedIp) {
      t.Errorf("Invalid value %d for %s", parsedIp, validIpSource)
    }
  }

  for _, invalidIpSource := range invalidValues {
    parsedInvalidIp, invalidIpErr := ParseIpv6Address(invalidIpSource)
    if invalidIpErr == nil {
      t.Errorf("Not return error for invalid ip %s", invalidIpSource)
    }
    if parsedInvalidIp != nil {
      t.Error("Return value for invali ip %s", invalidIpSource)
    }
  }
}

func TestEqual(t *testing.T) {
  var firstAddress = Ipv6Address { 0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF, 0x12, 0x34, 0x56, 0x78, 0x90, 0xAB, 0xCD, 0xEF }
  var secondAddress = Ipv6Address { 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }

  if !firstAddress.Equal(&firstAddress) {
    t.Error("Address not equal youself")
  }

  if firstAddress.Equal(&secondAddress) {
    t.Error("Two diffirent address equal")
  }
}

func TestCompare(t *testing.T) {
  big := "::FFFF"
  bigAddress, _ := ParseIpv6Address(big)
  small := "::0000"
  smallAddress, _ := ParseIpv6Address(small)

  if bigAddress.Compare(smallAddress) != 1 {
    t.Errorf("%s.Compare(%s) != 1", big, small)
  }

  if smallAddress.Compare(bigAddress) != -1 {
    t.Errorf("%s.Compare(%s) != -1", small, big)
  }

  if bigAddress.Compare(bigAddress) != 0 {
    t.Errorf("%s.Compare(%s) != 0", big, big)
  }
}
