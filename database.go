package geoip

import (
  "io"
  "encoding/csv"
  "errors"
  "fmt"
)

var HEADER = []string {
  "network_start_ip",
  "network_mask_length",
  "geoname_id",
  "registered_country_geoname_id",
  "represented_country_geoname_id",
  "postal_code",
  "latitude",
  "longitude",
  "is_anonymous_proxy",
  "is_satellite_provider",
}

const (
  IS_ANONYMOUS_PROXY_INDEX = 8
  IS_SATELLITE_PROVIDER_INDEX = 9
)

type Database []*Record

type recordItem struct {
  record *Record
  prev *recordItem
}

func validHeader(header []string) bool {
  if len(header) != len(HEADER) {
    return false
  }
  for index, value := range HEADER {
    if value != header[index] {
      return false
    }
  }
  return true
}

func LoadDatabase(reader *io.Reader) (*Database, *error) {
  csvReader := csv.NewReader(*reader)

  header, headerError := csvReader.Read()
  if headerError != nil {
    return nil, &headerError
  }

  if !validHeader(header) {
    headerError := errors.New("Database has invalid header")
    return nil, &headerError
  }

  var rowCount int
  var item *recordItem

  for {
    rowItems, rowError := csvReader.Read()
    if io.EOF == rowError {
      break
    }
    if nil != rowError {
      return nil, &rowError
    }
    if "1" == rowItems[IS_ANONYMOUS_PROXY_INDEX] ||
        "1" == rowItems[IS_SATELLITE_PROVIDER_INDEX] {
      continue
    }
    rowCount++
    record, recordParseError := ParseRecord(rowItems)
    if nil != recordParseError {
      recordError := fmt.Errorf("Row #%d %#v parsed with error %v",
        rowCount, rowItems, *recordParseError)
      return nil, &recordError
    }
    item = &recordItem {
      record: record,
      prev: item,
    }
  }

  var database Database = make([]*Record, rowCount)

  for index := rowCount - 1; index >= 0; index-- {
    database[index] = item.record
    item = item.prev
  }

  return &database, nil
}

func (database *Database) get(index uint) (*Record) {
  return (*database)[index]
}

func (database *Database) FindRecord(ipAddress *Ipv6Address) (*Record) {
  startIndex := uint(0)
  endIndex := uint(len(*database) - 1)
  for {
    centerIndex := startIndex + (endIndex - startIndex) / 2
    record := database.get(centerIndex)
    switch {
    case record.MatchIpAddress(ipAddress):
      return record
    case startIndex == endIndex:
      return nil
    case startIndex == centerIndex:
      startIndex = endIndex
    case record.ipAddress.Compare(ipAddress) < 0:
      startIndex = centerIndex
    default:
      endIndex = centerIndex
    }
  }
  return nil
}
