package geoip

import (
  "strconv"
  "fmt"
)

const MAX_MASK_VALUE = 128

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
  mask, maskError := strconv.ParseUint(rowItems[1], 10, 8)
  if nil != maskError {
    return nil, &maskError
  }

  if mask > MAX_MASK_VALUE {
    error := fmt.Errorf("Mask should have value be between 1 and %d", MAX_MASK_VALUE)
    return nil, &error
  }

  geonameId, geonameidError := strconv.ParseUint(rowItems[2], 10, 64)
  if nil != geonameidError {
    return nil, &geonameidError
  }

  record := &Record {
    ipAddress: *ipAddress,
    mask: byte(mask),
    geonameId: geonameId,
  }
  return record, nil
}
