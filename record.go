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
  NETWORK_START_IP_INDEX = 0
  NETWORK_MASK_LENGTH_INDEX = 1
  GEONAME_ID_INDEX = 2
  IS_ANONYMOUS_PROXY_INDEX = 8
  IS_SATELLITE_PROVIDER_INDEX = 9
  TRUE_VALUE = "1"
)

var (
  COLUMN_COUNT = len(HEADER)
  USELESS_RECORS_ERROR = errors.New("Anonymous proxy or satellite provider")
)

type Record struct {
  ipAddress Ipv6Address
  mask byte
  geonameId uint32
}

func ParseRecord(rowItems []string) (*Record, *error) {
  if COLUMN_COUNT != len(rowItems) {
    lenErr := fmt.Errorf("len(%#v) != %s", rowItems, COLUMN_COUNT)
    return nil, &lenErr
  }

  if TRUE_VALUE == rowItems[IS_ANONYMOUS_PROXY_INDEX] ||
      TRUE_VALUE == rowItems[IS_SATELLITE_PROVIDER_INDEX] {
    return nil, &USELESS_RECORS_ERROR
  }

  ipStr := rowItems[NETWORK_START_IP_INDEX]
  ipAddress, ipError := ParseIpv6Address(ipStr)
  if ipError != nil {
    return nil, ipError
  }
  maskStr := rowItems[NETWORK_MASK_LENGTH_INDEX]
  mask, maskParseError := strconv.ParseUint(maskStr, 10, 8)
  if nil != maskParseError {
    maskError := fmt.Errorf("Error parsing mask \"%s\"", maskParseError)
    return nil, &maskError
  }

  if mask > MAX_MASK_VALUE {
    maskError := fmt.Errorf("Mask should have value be between 1 and %d", MAX_MASK_VALUE)
    return nil, &maskError
  }

  geonameIdStr := rowItems[GEONAME_ID_INDEX]
  geonameId, geonameIdParseError := strconv.ParseUint(geonameIdStr, 10, 32)
  if nil != geonameIdParseError {
    geonameIdError := fmt.Errorf("Error parsing geoname_id \"%s\"", geonameIdParseError)
    return nil, &geonameIdError
  }

  record := &Record {
    ipAddress: *ipAddress,
    mask: byte(mask),
    geonameId: uint32(geonameId),
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
