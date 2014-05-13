package geoipcsv

import (
	"io"
	"strings"
	"testing"
)

const VALID_DATA = `
network_start_ip,network_mask_length,geoname_id,registered_country_geoname_id,represented_country_geoname_id,postal_code,latitude,longitude,is_anonymous_proxy,is_satellite_provider
::ffff:1.0.64.0,114,1861060,1861060,,,35.6900,139.6900,0,0
::ffff:1.0.32.0,115,1809858,1814991,,,23.1167,113.2500,0,0
::ffff:1.0.16.0,116,1850147,1861060,,,35.6850,139.7514,0,0
::ffff:1.0.16.0,116,1850147,1861060,,,35.6850,139.7514,0,1
::ffff:1.0.16.0,116,1850147,1861060,,,35.6850,139.7514,1,0
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

var INVALID_DATA = []string{
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
		for recordIndex, record := range *validDatabase {
			if nil == record {
				t.Errorf("Database have nil on row #%d", recordIndex)
			}
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

func TestFindRecord(t *testing.T) {

	var ipAddresses = []string{
		"::ffff:1.1.1.1",
		"::ffff:2.2.2.2",
		"::ffff:3.3.3.3",
	}
	var validIpAddresses = []string{
		"::ffff:1.1.1.1",
		"::ffff:1.1.1.0",
		"::ffff:1.1.1.255",
		"::ffff:2.2.2.2",
		"::ffff:2.2.2.0",
		"::ffff:2.2.2.255",
		"::ffff:3.3.3.3",
		"::ffff:3.3.3.0",
		"::ffff:3.3.3.255",
	}
	var invalidIpAddresses = []string{
		"::ffff:4.4.4.4",
		"::ffff:5.5.5.5",
		"::ffff:3.3.4.0",
	}

	database := make(Database, len(ipAddresses))
	for index, ipAddress := range ipAddresses {
		parsedIpAddress, parsingError := ParseIpv6Address(ipAddress)
		if nil != parsingError {
			t.Error("Error parsing ip %s", parsingError)
			return
		}
		database[index] = &Record{ipAddress: *parsedIpAddress, mask: 120}
	}

	for _, ipAddress := range validIpAddresses {
		parsedIpAddress, parsingError := ParseIpv6Address(ipAddress)
		if nil != parsingError {
			t.Errorf("Error parsing ip %s", parsingError)
			return
		}
		record := database.FindRecord(parsedIpAddress)
		if nil == record {
			t.Errorf("Not found record for ip %s", ipAddress)
		}
	}

	for _, ipAddress := range invalidIpAddresses {
		parsedIpAddress, parsingError := ParseIpv6Address(ipAddress)
		if nil != parsingError {
			t.Errorf("Error parsing ip %s", parsingError)
			return
		}
		record := database.FindRecord(parsedIpAddress)
		if nil != record {
			t.Errorf("Found record %#v for invalid ip %s", record, ipAddress)
		}
	}
}
