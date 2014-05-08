package geoip

import (
  "strconv"
  "fmt"
  "bytes"
  "errors"
)

const (
  BYTE_SIZE = 8
  MAX_MASK_VALUE = IPV6_BYTES * BYTE_SIZE
  IS_ANONYMOUS_PROXY_INDEX = 8
  IS_SATELLITE_PROVIDER_INDEX = 9
)

var (
  COLUMN_COUNT = len(HEADER)
  USELESS_RECORS_ERROR = errors.New("Anonymous proxy or satellite provider")
)

type Record struct {
  ipAddress Ipv6Address
  mask byte
  geonameId uint64
}

func ParseRecord(rowItems []string) (*Record, *error) {
  if COLUMN_COUNT != len(rowItems) {
    lenErr := fmt.Errorf("len(%#v) != %s", rowItems, COLUMN_COUNT)
    return nil, &lenErr
  }

  if "1" == rowItems[IS_ANONYMOUS_PROXY_INDEX] ||
      "1" == rowItems[IS_SATELLITE_PROVIDER_INDEX] {
    return nil, &USELESS_RECORS_ERROR
  }

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
  return (record.ipAddress[equalByteCount] & bitMask) == (ipAddress[equalByteCount] & bitMask)
}
