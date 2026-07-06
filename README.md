# Sales Tax Automation Code Samples

Welcome to the official `sales-tax-code-samples` repository! This repository contains examples of how to integrate and automate sales tax calculations using the [zip.tax API](https://zip.tax/) **v60** — the latest and most capable version of the API. These code samples are designed to help developers quickly implement and adapt the API in their applications.

## ✨ What's New in v60

Version 60 is a major upgrade over previous API versions, introducing a completely restructured response format and a set of powerful new capabilities:

| Feature | Description |
|---|---|
| **Structured jurisdiction rates (`baseRates`)** | Every component rate (state, county, city, special districts) is returned as a discrete object with jurisdiction type, name, description, and tax code — no more fixed `district1`–`district5` columns. |
| **Tax summaries (`taxSummaries`)** | Pre-computed total sales and use tax rates with display-ready breakdowns, ready to show at checkout. |
| **Sourcing rules (`sourcingRules`)** | Explicit origin- vs destination-based sourcing determination for the queried location. |
| **Product taxability codes (`taxabilityCode`)** | Pass a TIC (Taxability Information Code) to get product-specific tax rules — reduced food rates, exemption thresholds, per-volume taxes, effective dates, and more — returned in `productDetail`. |
| **Canadian tax rates (`countryCode=CAN`)** | GST / PST / HST / QST lookups for Canadian postal codes, with province-level breakdowns. |
| **Historical rate lookups (`historical=YYYYMM`)** | Query the tax rates that were in effect for a past period — essential for amended returns and audits. |
| **Extended address details (`addressDetailExtended=true`)** | Fully parsed address components (street, house number, city, county, state, postal code, country) from the geocoder. |
| **Extended shipping rules (`shippingExtended=true`)** | State-level freight taxability rules, including whether separately-stated shipping is exempt. |
| **Service & freight taxability** | `service` and `shipping` blocks indicate whether services and freight are taxable at the destination. |
| **Tennessee Single Article Tax (`sat_item_total`)** | Correct local tax cap calculation for high-value single articles in Tennessee. |
| **Lookup by address, postal code, or coordinates** | Door-level accuracy via street address or `lat`/`lng`, or all applicable rates for a postal code. |
| **JSON and XML** | Choose your response format with `format=json` (default) or `format=xml`. |

> Some features (product taxability codes, Canadian rates, historical data) are plan-level entitlements. See [zip.tax pricing](https://zip.tax/pricing) for details.

## 🚀 Quick Start

All lookups are a single GET request:

```
GET https://api.zip-tax.com/request/v60?key=YOUR_API_KEY&address=200 Spectrum Center Dr, Irvine, CA 92618
```

Example v60 response (abbreviated):

```json
{
  "metadata": {
    "version": "v60",
    "response": { "code": 100, "name": "RESPONSE_CODE_SUCCESS", "message": "Successful API Request." }
  },
  "baseRates": [
    { "rate": 0.0725, "jurType": "US_STATE_SALES_TAX",  "jurName": "CA",     "jurDescription": "US State Sales Tax",  "jurTaxCode": "06" },
    { "rate": 0.005,  "jurType": "US_COUNTY_SALES_TAX", "jurName": "ORANGE", "jurDescription": "US County Sales Tax", "jurTaxCode": "30" },
    { "rate": 0,      "jurType": "US_CITY_SALES_TAX",   "jurName": "IRVINE", "jurDescription": "US City Sales Tax",   "jurTaxCode": null }
  ],
  "service":  { "adjustmentType": "SERVICE_TAXABLE", "taxable": "N", "description": "Services non-taxable" },
  "shipping": { "adjustmentType": "FREIGHT_TAXABLE", "taxable": "N", "description": "Freight non-taxable" },
  "sourcingRules": { "adjustmentType": "ORIGIN_DESTINATION", "description": "Destination Based Taxation", "value": "D" },
  "taxSummaries": [
    { "rate": 0.0775, "taxType": "SALES_TAX", "summaryName": "Total Base Sales Tax", "displayRates": [ { "name": "Total Rate", "rate": 0.0775 } ] },
    { "rate": 0.0775, "taxType": "USE_TAX",   "summaryName": "Total Base Use Tax",   "displayRates": [ { "name": "Total Rate", "rate": 0.0775 } ] }
  ],
  "addressDetail": {
    "normalizedAddress": "200 Spectrum Center Dr, Irvine, CA 92618-5003, United States",
    "incorporated": "true",
    "geoLat": 33.65253,
    "geoLng": -117.74794
  }
}
```

## 🔎 Request Parameters

| Parameter | Description |
|---|---|
| `key` | **Required.** Your API key. |
| `address` | Street address for a door-level rate lookup (geocoded and normalized). |
| `postalcode` | US ZIP or Canadian postal code — returns all applicable rates for the code. |
| `lat`, `lng` | Geographic coordinates for a location-based lookup. |
| `city`, `county`, `state` | Optional filters to narrow postal-code lookups. |
| `countryCode` | `USA` (default) or `CAN` for Canadian tax rates. |
| `taxabilityCode` | TIC code for product-specific tax rules (e.g. `40030` = food & food ingredients). |
| `historical` | `YYYYMM` period for historical rate lookups (e.g. `202401`). |
| `addressDetailExtended` | `true` to include parsed address components in `addressDetail.address`. |
| `shippingExtended` | `true` to include state-level shipping rules in `shipping.shippingExtended`. |
| `sat_item_total` | Item total for Tennessee Single Article Tax calculation. |
| `adjustment` | Sourcing adjustment: `auto`, `origin`, or `destination`. |
| `format` | `json` (default) or `xml`. |

## 📂 Repository Structure

```
├── Go/
│   ├── example.go
│   └── README.md
├── Node/
│   ├── example.js
│   └── README.md
├── Python/
│   ├── example.py
│   └── README.md
├── Scala/
│   ├── example.scala
│   └── README.md
└── README.md
```

Each language folder contains a runnable example demonstrating the core v60 features: address lookup, coordinate lookup, product taxability codes, Canadian rates, and extended address/shipping details.

## 🛠️ Prerequisites

Before using the code samples, ensure you have:

1. A [zip.tax API Key](https://zip.tax/pricing).
2. Basic knowledge of the programming language you're working with.
3. Installed dependencies for the selected code sample (refer to the specific sample's README).

## 🔧 Getting Started

1. Clone this repository:

    ```bash
    git clone https://github.com/ZipTax/sales-tax-code-samples.git
    cd sales-tax-code-samples
    ```

2. Navigate to the directory of the programming language you want to use (e.g., `Python/`, `Go/`, `Node/`, or `Scala/`).

3. Set your API key as an environment variable:

    ```bash
    export ZIPTAX_API_KEY=your_api_key_here
    ```

4. Follow the language-specific instructions to install dependencies and run the sample.

### Example: Python

```bash
cd Python
pip install requests
python example.py
```

## 📟 Response Codes

| Code | Meaning |
|---|---|
| 100 | Successful API request |
| 101 | Invalid or missing API key |
| 104 | Invalid postal code |
| 105 | Invalid query string |
| 107 | Feature/version not enabled for your plan |
| 108 | Rate limit exceeded (HTTP 429) |
| 109 | Address missing, incomplete, or invalid |
| 110 | Valid query, but no result found |
| 111 | Invalid `historical` parameter |
| 112 | International rates not enabled for this key |
| 113 | Product rate rules not enabled for this key |

## 📚 Documentation

- [API Reference](https://developers.zip.tax)

## 🤝 Contributing

We welcome contributions! If you have improvements, additional language examples, or bug fixes, please:

1. Fork this repository.
2. Create a new branch (`feature/my-feature` or `fix/issue-name`).
3. Submit a pull request.

## 📝 License

This project is licensed under the [MIT License](LICENSE).

---

### 🌟 Support

If you encounter any issues or have questions, feel free to open an [issue](https://github.com/ZipTax/sales-tax-code-samples/issues) or contact our [support team](https://zip.tax/contact).

Happy Coding! 🚀
