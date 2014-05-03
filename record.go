package geoip

import (
  "strconv"
  "fmt"
  "bytes"
)

const BYTE_SIZE = 8
const MAX_MASK_VALUE = IPV6_BYTES * BYTE_SIZE

type Record struct {
  ipAddress Ipv6Address
  mask byte
  geonameId uint64
}

func ParseRecord(rowItems []string) (*Record, *error) {
  ipAddress, ipError := ParseIpv6Address(rowItems[0])
  if ipError != nil {
    return nil, ipError
  }
  mask, maskParseError := strconv.ParseUint(rowItems[1], 10, 8)
  if nil != maskParseError {
    maskError := fmt.Errorf("Error parsing mask \"%s\"", maskParseError)
    return nil, &maskError
  }

  if mask > MAX_MASK_VALUE {
    maskError := fmt.Errorf("Mask should have value be between 1 and %d", MAX_MASK_VALUE)
    return nil, &maskError
  }

  geonameId, geonameIdParseError := strconv.ParseUint(rowItems[2], 10, 64)
  if nil != geonameIdParseError {
    geonameIdError := fmt.Errorf("Error parsing geoname_id \"%s\"", geonameIdParseError)
    return nil, &geonameIdError
  }

  record := &Record {
    ipAddress: *ipAddress,
    mask: byte(mask),
    geonameId: geonameId,
  }
  return record, nil
}

func (record *Record) MatchIpAddress(ipAddress *Ipv6Address) bool {
  equalByteCount := record.mask / BYTE_SIZE
  if !bytes.Equal(record.ipAddress[:equalByteCount], ipAddress[:equalByteCount]) {
    return false
  }
  var bitOffset byte = record.mask % BYTE_SIZE
  if bitOffset == 0 {
    return true
  }
  var bitMask byte = 0xFF << (BYTE_SIZE - bitOffset)
  fmt.Printf("%d %d %b\n", record.mask, equalByteCount, bitMask)
  return (record.ipAddress[equalByteCount] & bitMask) == (ipAddress[equalByteCount] & bitMask)
}
