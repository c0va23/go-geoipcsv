package geoip

import (
  "testing"
)

func TestParseRecord(t *testing.T) {
  var validRowItems = []string { "ffff::ffff:1.2.3.4", "123", "456", "0", "", "", "0.0", "0.0", "0,0" }
  validRecord := &Record {
    ipAddress: Ipv6Address { 
      0xFF, 0xFF, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0xFF, 0xFF,
      0x01, 0x02, 0x03, 0x04,
    },
    mask: 123,
    geonameId: 456,
  }
  parsedValidRecord, validRowError := ParseRecord(validRowItems)
  if nil != validRowError {
    t.Errorf("Return error \"%s\" for valid row", validRowError)
  }
  if *validRecord != *parsedValidRecord {
    t.Errorf("Return invalid value for valid row (expected: %#v, parsed: %#v)", validRecord, parsedValidRecord)
  }

  var invalidRows = [][]string {
    { "ffff::ffff:1.2.3.256", "123", "456", "0", "", "", "0.0", "0.0", "0,0" },
    { "ffff::ffff:1.2.3.3", "256", "456", "0", "", "", "0.0", "0.0", "0,0" },
    { "ffff::ffff:1.2.3.3", "123", "18446744073709551617", "0", "", "", "0.0", "0.0", "0,0" },
  }

  for _, invalidRowItems := range invalidRows {
    parsedInvalidRow, invalidRowError := ParseRecord(invalidRowItems)
    if nil == invalidRowError {
      t.Errorf("Not return error for invalid row %s", invalidRowItems)
    }
    if nil != parsedInvalidRow {
      t.Errorf("Return record \"%#v\" for invalid row \"%s\"", *parsedInvalidRow, invalidRowItems)
    }
  }

}

func testMatchIpAddresses(t *testing.T, ipAddress string, maskLength byte, validIpAddresses []string, invalidIpAddresses []string) {
  parsedIpAddress, parsingError := ParseIpv6Address(ipAddress)
  if parsingError != nil {
    t.Errorf("Error parsing ip %s", *parsingError)
    return
  }

  record := Record {
    ipAddress: *parsedIpAddress,
    mask: maskLength,
  }

  if !record.MatchIpAddress(parsedIpAddress) {
    t.Error("Record not match own address")
  }

  for _, validIpAddress := range validIpAddresses {
    parsedIpAddress, parsingError := ParseIpv6Address(validIpAddress)
    if parsingError != nil {
      t.Errorf("Address %s parsed with error %s", validIpAddress, *parsingError)
      return
    }
    if !record.MatchIpAddress(parsedIpAddress) {
      t.Errorf("Record %s/%d not match valid address %s", ipAddress, maskLength, validIpAddress)
    }
  }

  for _, invalidIpAddress := range invalidIpAddresses {
    parsedIpAddress, parsingError := ParseIpv6Address(invalidIpAddress)
    if parsingError != nil {
      t.Errorf("Address %s parsed with error %s", invalidIpAddress, *parsingError)
      return
    }
    if record.MatchIpAddress(parsedIpAddress) {
      t.Errorf("Record %s/%d match not valid address %v", ipAddress, maskLength, invalidIpAddress)
    }
  }

}

func TestMatchIpAddress(t *testing.T) {

  testMatchIpAddresses(t, "::8888", 128, []string { }, []string { "::8887", "::8889" })
  testMatchIpAddresses(t, "::8888", 124, []string { "::8880", "::888F" },
    []string { "::887F", "::8890", "::87FF", "::8900" })
  testMatchIpAddresses(t, "::8888", 120, []string { "::8800", "::88FF" },
    []string { "::87FF", "::8900", "::7FFF", "::9000" })
}
