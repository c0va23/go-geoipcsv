package geoip

import (
  "testing"
)

func TestParseRecord(t *testing.T) {
  var validRowItems = []string { "::ffff:1.2.3.4", "123", "456", "0", "", "", "0.0", "0.0", "0", "0" }
  validRecord := &Record {
    ipAddress: Ipv6Address { 
      0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0xFF, 0xFF,
      0x01, 0x02, 0x03, 0x04,
    },
    mask: 123,
    geonameId: 456,
  }
  parsedValidRecord, validRowError := ParseRecord(validRowItems)
  if nil != validRowError {
    t.Errorf("Return error \"%v\" for valid row", *validRowError)
  }
  if nil == parsedValidRecord {
    t.Errorf("Return nil for valid row %v", validRowItems)
  } else if *validRecord != *parsedValidRecord {
    t.Errorf("Return invalid value for valid row (expected: %v, parsed: %v)", *validRecord, *parsedValidRecord)
  }

  var invalidRows = [][]string {
    { "::ffff:1.2.3.256", "123", "456", "0", "", "", "0.0", "0.0", "0", "0" },
    { "::ffff:1.2.3.3", "256", "456", "0", "", "", "0.0", "0.0", "0", "0" },
    { "::ffff:1.2.3.3", "123", "18446744073709551617", "0", "", "", "0.0", "0.0", "0", "0" },
    { "::ffff:1.2.3.4", "123", "456", "0", "", "", "0.0", "0.0", "1", "0" },
    { "::ffff:1.2.3.4", "123", "456", "0", "", "", "0.0", "0.0", "0", "1" },
    { "::ffff:1.2.3.4", "123", "456", "0", "", "", "0.0", "0.0" },
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
  testMatchIpAddresses(t, "::8888", 124, []string { "::8880", "::888f" },
    []string { "::887f", "::8890", "::87ff", "::8900" })
  testMatchIpAddresses(t, "::8888", 120, []string { "::8800", "::88ff" },
    []string { "::87ff", "::8900", "::7fff", "::9000" })
}
