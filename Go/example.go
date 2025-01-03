package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Response struct {
	Version       string        `json:"version"`
	RCode         int           `json:"rCode"`
	Results       []Result      `json:"results"`
	AddressDetail AddressDetail `json:"addressDetail"`
}

type Result struct {
	GeoPostalCode     string  `json:"geoPostalCode"`
	GeoCity           string  `json:"geoCity"`
	GeoCounty         string  `json:"geoCounty"`
	GeoState          string  `json:"geoState"`
	TaxSales          float64 `json:"taxSales"`
	TaxUse            float64 `json:"taxUse"`
	TxbService        string  `json:"txbService"`
	TxbFreight        string  `json:"txbFreight"`
	StateSalesTax     float64 `json:"stateSalesTax"`
	StateUseTax       float64 `json:"stateUseTax"`
	CitySalesTax      float64 `json:"citySalesTax"`
	CityUseTax        float64 `json:"cityUseTax"`
	CityTaxCode       string  `json:"cityTaxCode"`
	CountySalesTax    float64 `json:"countySalesTax"`
	CountyUseTax      float64 `json:"countyUseTax"`
	CountyTaxCode     string  `json:"countyTaxCode"`
	DistrictSalesTax  float64 `json:"districtSalesTax"`
	DistrictUseTax    float64 `json:"districtUseTax"`
	District1Code     string  `json:"district1Code"`
	District1SalesTax float64 `json:"district1SalesTax"`
	District1UseTax   float64 `json:"district1UseTax"`
	District2Code     string  `json:"district2Code"`
	District2SalesTax float64 `json:"district2SalesTax"`
	District2UseTax   float64 `json:"district2UseTax"`
	District3Code     string  `json:"district3Code"`
	District3SalesTax float64 `json:"district3SalesTax"`
	District3UseTax   float64 `json:"district3UseTax"`
	District4Code     string  `json:"district4Code"`
	District4SalesTax float64 `json:"district4SalesTax"`
	District4UseTax   float64 `json:"district4UseTax"`
	District5Code     string  `json:"district5Code"`
	District5SalesTax float64 `json:"district5SalesTax"`
	District5UseTax   float64 `json:"district5UseTax"`
	OriginDestination string  `json:"originDestination"`
}

type AddressDetail struct {
	NormalizedAddress string  `json:"normalizedAddress"`
	Incorporated      string  `json:"incorporated"`
	GeoLat            float64 `json:"geoLat"`
	GeoLng            float64 `json:"geoLng"`
}

func getSalesTax(address string, apiKey string) (*Response, error) {
	url := fmt.Sprintf("https://api.zip-tax.com/request/v50?key=%s&address=%s", apiKey, url.QueryEscape(address))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var taxResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&taxResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &taxResponse, nil
}

func main() {
	apiKey := "your_api_key_here"                         // Replace with your key
	address := "200 Spectrum Center Dr, Irvine, CA 92618" // Example address

	taxInfo, err := getSalesTax(address, apiKey)
	if err != nil {
		log.Fatalf("Error fetching sales tax: %v", err)
	}

	fmt.Printf("Normalized Address: %s\n", taxInfo.AddressDetail.NormalizedAddress)
	fmt.Printf("Address Lat/Lng: %f, %f\n", taxInfo.AddressDetail.GeoLat, taxInfo.AddressDetail.GeoLng)
	fmt.Printf("Rate: %.2f%%\n", taxInfo.Results[0].TaxSales*100)
}
