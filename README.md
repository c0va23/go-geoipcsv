# geoipcsv
Golang library for parsing and searching by csv-database MaxMind GeoIP (http://dev.maxmind.com/geoip/geoip2/geolite2/). 
Database fully loaded into memory.
Search perform via binary searching.

## Usage

### Import
```go
import "github.com/c0va23/go-geoipcsv"
```

### Parse databse
```go
var databaseReader io.Reader
database, databaseErr := geoipcsv.LoadDatabase(&databaseReader)
```

### Parse IPv6-address
```go
ipSrc := "::ffff:8.8.8.8"
ipAddress, ipAddressErr := geoicsv.ParseIpv6Address(ipSrc)
```

### Search ip-address
```go
record := database.FindRecord(ipAddress)
```

### Geoname ID
```go
println(record.GeonameId())
```
