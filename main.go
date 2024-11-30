package domaininfo

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type DomainInfo struct {
	OriginalInput string
	CleanDomain   string
	IPAddress     string
	Location      *LocationDetails
}

type LocationDetails struct {
	IP          string  `json:"ip"`
	City        string  `json:"city,omitempty"`
	Region      string  `json:"region,omitempty"`
	Country     string  `json:"country_name,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
}

func ValidateDomain(input string) (*DomainInfo, error) {
	cleanDomain := cleanDomainInput(input)

	if !isValidDomainFormat(cleanDomain) {
		return nil, fmt.Errorf("invalid domain format")
	}

	if !checkDNSResolution(cleanDomain) {
		return nil, fmt.Errorf("cannot resolve domain")
	}

	ipAddress, err := getIPAddress(cleanDomain)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve IP: %v", err)
	}

	location, err := getIPLocation(ipAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch location: %v", err)
	}

	return &DomainInfo{
		OriginalInput: input,
		CleanDomain:   cleanDomain,
		IPAddress:     ipAddress,
		Location:      location,
	}, nil
}

func cleanDomainInput(input string) string {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		parsedURL, err := url.Parse(input)
		if err == nil {
			input = parsedURL.Hostname()
		}
	}

	input = strings.TrimPrefix(input, "www.")
	return strings.TrimSpace(input)
}

func isValidDomainFormat(domain string) bool {
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z]{2,})+$`)
	return domainRegex.MatchString(domain)
}

func getIPAddress(domain string) (string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil || len(ips) == 0 {
		return "", err
	}
	return ips[0].String(), nil
}

func checkDNSResolution(domain string) bool {
	_, err := net.LookupIP(domain)
	return err == nil
}

func getIPLocation(ip string) (*LocationDetails, error) {
	locationProviders := []func(string) (*LocationDetails, error){
		getIPAPILocation,
		getIPInfoLocation,
		getFreeGeoIPLocation,
	}

	for _, provider := range locationProviders {
		location, err := provider(ip)
		if err == nil && location != nil {
			return location, nil
		}
	}

	return nil, fmt.Errorf("could not fetch location from any provider")
}

func getIPAPILocation(ip string) (*LocationDetails, error) {
	resp, err := http.Get(fmt.Sprintf("https://ipapi.co/%s/json/", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var location LocationDetails
	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, err
	}

	if location.City == "" && location.Country == "" {
		return nil, fmt.Errorf("no location data")
	}

	return &location, nil
}

func getIPInfoLocation(ip string) (*LocationDetails, error) {
	resp, err := http.Get(fmt.Sprintf("https://ipinfo.io/%s/json", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	location := &LocationDetails{}
	if loc, ok := data["loc"].(string); ok {
		coords := strings.Split(loc, ",")
		if len(coords) == 2 {
			fmt.Sscanf(coords[0], "%f", &location.Latitude)
			fmt.Sscanf(coords[1], "%f", &location.Longitude)
		}
	}

	location.City, _ = data["city"].(string)
	location.Region, _ = data["region"].(string)
	location.Country, _ = data["country"].(string)

	if location.City == "" && location.Country == "" {
		return nil, fmt.Errorf("no location data")
	}

	return location, nil
}

func getFreeGeoIPLocation(ip string) (*LocationDetails, error) {
	resp, err := http.Get(fmt.Sprintf("https://freegeoip.app/json/%s", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var location LocationDetails
	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, err
	}

	if location.City == "" && location.Country == "" {
		return nil, fmt.Errorf("no location data")
	}

	return &location, nil
}
