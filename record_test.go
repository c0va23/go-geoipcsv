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
