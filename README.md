# Domain Information Retrieval Package

## Overview

This Go package provides a robust solution for domain validation, IP resolution, and geolocation information retrieval. It offers a simple and efficient way to extract detailed information about a given domain.

## Features

- Domain format validation
- DNS resolution check
- IP address retrieval
- Geolocation information extraction
- Multiple location data providers
- Error handling
- Easy-to-use API

## Installation

```bash
go get github.com/ByteBreach/domaininfo
```

## Usage Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/yourusername/domaininfo"
)

func main() {
    domain := "google.com"
    
    // Validate domain and retrieve information
    domainInfo, err := domaininfo.ValidateDomain(domain)
    if err != nil {
        log.Fatalf("Domain validation failed: %v", err)
    }

    // Print domain details
    fmt.Printf("Domain: %s\n", domainInfo.CleanDomain)
    fmt.Printf("IP Address: %s\n", domainInfo.IPAddress)
    
    // Print location information
    location := domainInfo.Location
    fmt.Printf("City: %s\n", location.City)
    fmt.Printf("Region: %s\n", location.Region)
    fmt.Printf("Country: %s\n", location.Country)
    fmt.Printf("Coordinates: %.4f, %.4f\n", 
        location.Latitude, 
        location.Longitude)
}
```

## Package Structure

### Main Types

- `DomainInfo`: Contains comprehensive domain details
  - `OriginalInput`: Original domain input
  - `CleanDomain`: Sanitized domain name
  - `IPAddress`: Resolved IP address
  - `Location`: Geographical location details

- `LocationDetails`: Geographical information
  - `IP`: IP address
  - `City`: City name
  - `Region`: Region/State
  - `Country`: Country name
  - `Latitude`: Geographical latitude
  - `Longitude`: Geographical longitude

## Functions

### `ValidateDomain(input string) (*DomainInfo, error)`

Primary function to validate and retrieve domain information.

#### Parameters
- `input`: Domain name or URL to validate

#### Returns
- `*DomainInfo`: Detailed domain information
- `error`: Validation or retrieval error

### Validation Steps

1. Clean and normalize domain input
2. Validate domain format
3. Check DNS resolution
4. Retrieve IP address
5. Fetch geolocation information

## Geolocation Providers

The package uses multiple geolocation providers to ensure reliable location data:
- ipapi.co
- ipinfo.io
- freegeoip.app

## Error Handling

Comprehensive error handling for various scenarios:
- Invalid domain format
- DNS resolution failure
- IP address retrieval issues
- Location data fetch problems

## Performance Considerations

- Concurrent geolocation provider checking
- Minimal external dependencies
- Fast domain validation

## Limitations

- Requires active internet connection
- Geolocation accuracy depends on external services
- Rate limits may apply from geolocation providers

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License

## Dependencies

- Go 1.16+
- Standard library packages

## Disclaimer

Location data accuracy is not guaranteed and depends on external geolocation services.
