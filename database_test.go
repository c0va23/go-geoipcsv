package geoip

import (
  "testing"
  "strings"
  "io"
)
const VALID_DATA = `
network_start_ip,network_mask_length,geoname_id,registered_country_geoname_id,represented_country_geoname_id,postal_code,latitude,longitude,is_anonymous_proxy,is_satellite_provider
::ffff:1.0.64.0,114,1861060,1861060,,,35.6900,139.6900,0,0
::ffff:1.0.32.0,115,1809858,1814991,,,23.1167,113.2500,0,0
::ffff:1.0.16.0,116,1850147,1861060,,,35.6850,139.7514,0,0
`
const INVALID_IP_DATA = `
network_start_ip,network_mask_length,geoname_id,registered_country_geoname_id,represented_country_geoname_id,postal_code,latitude,longitude,is_anonymous_proxy,is_satellite_provider
::GGGG:0.0.0.0,0,0,0,,,0.0,0.0,0,0
`
const INVALID_MASK_DATA = `
network_start_ip,network_mask_length,geoname_id,registered_country_geoname_id,represented_country_geoname_id,postal_code,latitude,longitude,is_anonymous_proxy,is_satellite_provider
:::0.0.0.0,129,0,0,,,0.0,0.0,0,0
`
const INVALID_GEONAME_ID_DATA = `
network_start_ip,network_mask_length,geoname_id,registered_country_geoname_id,represented_country_geoname_id,postal_code,latitude,longitude,is_anonymous_proxy,is_satellite_provider
:::0.0.0.0,0,18446744073709551617,0,,,0.0,0.0,0,0
`
const INVALID_HEADER = `
column1,column2,column3
`
const EMTPY_DATA = `
network_start_ip,network_mask_length,geoname_id,registered_country_geoname_id,represented_country_geoname_id,postal_code,latitude,longitude,is_anonymous_proxy,is_satellite_provider
::ffff:1.2.3.4,0,0,0,,,0.0,0.0,1,0
::ffff:1.2.3.4,0,0,0,,,0.0,0.0,0,1
`

var INVALID_DATA = []string {
  INVALID_IP_DATA,
  INVALID_MASK_DATA,
  INVALID_HEADER,
  INVALID_GEONAME_ID_DATA,
}

func TestLoadDatabase(t *testing.T) {
  var validDatabaseReader io.Reader = strings.NewReader(VALID_DATA)
  validDatabase, validErr := LoadDatabase(&validDatabaseReader)
  if validErr != nil {
    t.Errorf("Return error \"%s\" for valid data", *validErr)
  }
  if nil == validDatabase {
    t.Error("Not return database for valid data")
  } else {
    recordsCount := len(*validDatabase)
    if recordsCount != 3 {
      t.Errorf("Return not valid database for valid data (%d/3)", recordsCount)
    }
  }

  for _, invalidData := range INVALID_DATA {
    var invalidDatabaseReader io.Reader = strings.NewReader(invalidData)
    invalidDatabase, invalidError := LoadDatabase(&invalidDatabaseReader)
    if nil == invalidError {
      t.Error("Not return error for invalid data")
    }
    if nil != invalidDatabase {
      t.Error("Return data for invalid data")
    }
  }

  var emptyDatabaseReader io.Reader = strings.NewReader(EMTPY_DATA)
  emptyDatabase, emptyError := LoadDatabase(&emptyDatabaseReader)
  if nil != emptyError {
    t.Error("Return error for empty data")
  }
  if nil == emptyDatabase {
    t.Error("Not return database for empty database")
  }
  emptyRecordCount := len(*emptyDatabase)
  if emptyRecordCount > 0 {
    t.Errorf("Return %d records for empty database", emptyRecordCount)
  }
}
