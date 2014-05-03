package geoip

import (
  "io"
  "encoding/csv"
  "errors"
)

var HEADER = []string {
  "network_start_ip",
  "network_mask_length",
  "geoname_id",
  "registered_country_geoname_id",
  "represented_country_geoname_id",
  "postal_code",
  "latitude","longitude",
  "is_anonymous_proxy",
  "is_satellite_provider",
}

type Database []*Record

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
  var database Database = make([]*Record, 0)

  header, headerError := csvReader.Read()
  if headerError != nil {
    return nil, &headerError
  }

  if !validHeader(header) {
    headerError := errors.New("Database has invalid header")
    return nil, &headerError
  }

  for {
    rowItems, rowError := csvReader.Read()
    if io.EOF == rowError {
      break
    }
    if nil != rowError {
      return nil, &rowError
    }
    record, recordError := ParseRecord(rowItems)
    if nil != recordError {
      return nil, recordError
    }
    database = append(database, record)
  }

  return &database, nil
}

