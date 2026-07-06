// zip.tax API v60 examples.
//
// Demonstrates the v60 endpoint features:
//   - Address lookup with normalized address and geocoding
//   - Lookup by geographic coordinates (lat/lng)
//   - Structured jurisdiction rates (baseRates) and tax summaries
//   - Sourcing rules and service/freight taxability
//   - Extended address components and shipping rules
//   - Product taxability codes (TIC) via productDetail
//   - Canadian tax rates (GST/PST/HST)
//
// Uses only the Go standard library. Set your API key in the
// ZIPTAX_API_KEY environment variable.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

const baseURL = "https://api.zip-tax.com/request/v60"

type Response struct {
	Metadata      Metadata       `json:"metadata"`
	BaseRates     []BaseRate     `json:"baseRates"`
	Service       *Service       `json:"service,omitempty"`
	Shipping      Shipping       `json:"shipping"`
	SourcingRules *SourcingRules `json:"sourcingRules,omitempty"`
	TaxSummaries  []TaxSummary   `json:"taxSummaries"`
	ProductDetail *ProductDetail `json:"productDetail,omitempty"`
	AddressDetail AddressDetail  `json:"addressDetail"`
}

type Metadata struct {
	Version  string       `json:"version"`
	Response ResponseInfo `json:"response"`
}

type ResponseInfo struct {
	Code       int    `json:"code"`
	Name       string `json:"name"`
	Message    string `json:"message"`
	Definition string `json:"definition"`
}

type BaseRate struct {
	Rate           float64 `json:"rate"`
	JurType        string  `json:"jurType"`
	JurName        string  `json:"jurName"`
	JurDescription string  `json:"jurDescription"`
	JurTaxCode     *string `json:"jurTaxCode"`
}

type Service struct {
	AdjustmentType string `json:"adjustmentType"`
	Taxable        string `json:"taxable"`
	Description    string `json:"description"`
}

type Shipping struct {
	AdjustmentType   string            `json:"adjustmentType"`
	Taxable          string            `json:"taxable"`
	Description      string            `json:"description"`
	ShippingExtended *ShippingExtended `json:"shippingExtended,omitempty"`
}

type ShippingExtended struct {
	StateName                  string `json:"stateName"`
	StateCode                  string `json:"stateCode"`
	Rule                       string `json:"rule"`
	ExemptWhenSeparatelyStated string `json:"exemptWhenSeparatelyStated"`
	Description                string `json:"description"`
}

type SourcingRules struct {
	AdjustmentType string `json:"adjustmentType"`
	Description    string `json:"description"`
	Value          string `json:"value"`
}

type TaxSummary struct {
	Rate         float64       `json:"rate"`
	TaxType      string        `json:"taxType"`
	SummaryName  string        `json:"summaryName"`
	DisplayRates []DisplayRate `json:"displayRates"`
}

type DisplayRate struct {
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}

type ProductDetail struct {
	TaxabilityCode TaxabilityCode `json:"taxabilityCode"`
}

type TaxabilityCode struct {
	ID                string     `json:"id"`
	StateFIPS         string     `json:"stateFIPS"`
	CountyFIPS        string     `json:"countyFIPS"`
	Title             string     `json:"title"`
	Label             string     `json:"label"`
	RateActionCode    string     `json:"rateActionCode"`
	RateActionMessage string     `json:"rateActionMessage"`
	RateRules         []RateRule `json:"rateRules"`
}

type RateRule struct {
	JurTaxCode           *string  `json:"jurTaxCode"`
	EffectiveDt          *int64   `json:"effectiveDt"`
	ExpiresDt            *int64   `json:"expiresDt"`
	EffectiveTaxRate     *float64 `json:"effectiveTaxRate"`
	PercentTaxable       *float64 `json:"percentTaxable"`
	ExemptUnder          *float64 `json:"exemptUnder"`
	ExemptOver           *float64 `json:"exemptOver"`
	TaxablePortionOver   *float64 `json:"taxablePortionOver"`
	IsDestinationTaxType *bool    `json:"isDestinationTaxType"`
	IsFoodDrug           *bool    `json:"isFoodDrug"`
}

type AddressDetail struct {
	NormalizedAddress string             `json:"normalizedAddress"`
	Incorporated      string             `json:"incorporated"`
	GeoLat            float64            `json:"geoLat"`
	GeoLng            float64            `json:"geoLng"`
	Address           *AddressComponents `json:"address,omitempty"`
}

type AddressComponents struct {
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	StateCode   string `json:"stateCode"`
	State       string `json:"state"`
	County      string `json:"county"`
	City        string `json:"city"`
	Street      string `json:"street"`
	PostalCode  string `json:"postalCode"`
	HouseNumber string `json:"houseNumber"`
}

func requestV60(apiKey string, params url.Values) (*Response, error) {
	params.Set("key", apiKey)

	resp, err := http.Get(baseURL + "?" + params.Encode())
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

	if taxResponse.Metadata.Response.Code != 100 {
		return nil, fmt.Errorf("API error %d: %s",
			taxResponse.Metadata.Response.Code, taxResponse.Metadata.Response.Message)
	}

	return &taxResponse, nil
}

func printSummary(data *Response) {
	fmt.Printf("  Normalized Address: %s\n", data.AddressDetail.NormalizedAddress)
	fmt.Printf("  Lat/Lng: %f, %f\n", data.AddressDetail.GeoLat, data.AddressDetail.GeoLng)
	fmt.Printf("  Incorporated: %s\n", data.AddressDetail.Incorporated)
	if data.SourcingRules != nil {
		fmt.Printf("  Sourcing: %s (%s)\n", data.SourcingRules.Description, data.SourcingRules.Value)
	}
	if data.Service != nil {
		fmt.Printf("  Services taxable: %s | Freight taxable: %s\n", data.Service.Taxable, data.Shipping.Taxable)
	}

	fmt.Println("  Jurisdiction rates:")
	for _, rate := range data.BaseRates {
		fmt.Printf("    %-24s %-20s %.3f%%\n", rate.JurType, rate.JurName, rate.Rate*100)
	}

	for _, summary := range data.TaxSummaries {
		fmt.Printf("  %s: %.2f%%\n", summary.SummaryName, summary.Rate*100)
	}
}

func main() {
	apiKey := os.Getenv("ZIPTAX_API_KEY")
	if apiKey == "" {
		apiKey = "your_api_key_here"
	}

	// 1. Door-level rate lookup by street address, with extended address
	//    components and state shipping rules.
	fmt.Println("=== Address lookup (extended details) ===")
	data, err := requestV60(apiKey, url.Values{
		"address":               {"200 Spectrum Center Dr, Irvine, CA 92618"},
		"addressDetailExtended": {"true"},
		"shippingExtended":      {"true"},
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	printSummary(data)

	if components := data.AddressDetail.Address; components != nil {
		fmt.Printf("  Parsed components: %s %s, %s, %s %s\n",
			components.HouseNumber, components.Street, components.City,
			components.StateCode, components.PostalCode)
	}

	if ext := data.Shipping.ShippingExtended; ext != nil {
		fmt.Printf("  Shipping rule (%s): %s — exempt when separately stated: %s\n",
			ext.StateCode, ext.Rule, ext.ExemptWhenSeparatelyStated)
	}

	// 2. Lookup by geographic coordinates.
	fmt.Println("\n=== Coordinate lookup (lat/lng) ===")
	data, err = requestV60(apiKey, url.Values{
		"lat": {"33.65253"},
		"lng": {"-117.74794"},
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	printSummary(data)

	// 3. Product-specific tax rules via a Taxability Information Code (TIC).
	//    TIC 40030 = food and food ingredients. Requires the product_rates
	//    plan entitlement.
	fmt.Println("\n=== Product taxability code (TIC 40030 — food) ===")
	data, err = requestV60(apiKey, url.Values{
		"address":        {"100 Broadway, Nashville, TN 37201"},
		"taxabilityCode": {"40030"},
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if data.ProductDetail != nil {
		product := data.ProductDetail.TaxabilityCode
		fmt.Printf("  TIC %s: %s\n", product.ID, product.Title)
		fmt.Printf("  Rate action: %s — %s\n", product.RateActionCode, product.RateActionMessage)
		for _, rule := range product.RateRules {
			fmt.Printf("    Jurisdiction %s: effective %d, percent taxable %.0f%%, food/drug: %t\n",
				deref(rule.JurTaxCode), derefInt(rule.EffectiveDt),
				derefFloat(rule.PercentTaxable), derefBool(rule.IsFoodDrug))
		}
	}

	// 4. Canadian tax rates. Requires the rate_loc_can plan entitlement.
	fmt.Println("\n=== Canadian lookup (Toronto, ON) ===")
	data, err = requestV60(apiKey, url.Values{
		"countryCode": {"CAN"},
		"postalcode":  {"M5V 3L9"},
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	for _, rate := range data.BaseRates {
		fmt.Printf("  %-4s %-10s %.2f%% — %s\n", rate.JurType, rate.JurName, rate.Rate*100, rate.JurDescription)
	}
	for _, summary := range data.TaxSummaries {
		for _, display := range summary.DisplayRates {
			fmt.Printf("  %s: %s %.2f%%\n", summary.SummaryName, display.Name, display.Rate*100)
		}
	}

	// 5. Historical rates for a past period (YYYYMM). Requires the
	//    historical data plan entitlement.
	// data, err = requestV60(apiKey, url.Values{
	// 	"address":    {"200 Spectrum Center Dr, Irvine, CA 92618"},
	// 	"historical": {"202401"},
	// })
}

func deref(s *string) string {
	if s == nil {
		return "n/a"
	}
	return *s
}

func derefInt(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func derefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
