# zip.tax API v60 — Node.js Example

This example demonstrates the [zip.tax API](https://developers.zip.tax) **v60** in Node.js, including:

- Door-level rate lookup by street address (with normalized address and geocoding)
- Rate lookup by geographic coordinates (`lat`/`lng`)
- Structured jurisdiction rates (`baseRates`) and checkout-ready totals (`taxSummaries`)
- Origin/destination sourcing rules and service/freight taxability
- Extended address components (`addressDetailExtended=true`)
- Extended state shipping rules (`shippingExtended=true`)
- Product-specific tax rules via Taxability Information Codes (`taxabilityCode`)
- Canadian GST/PST/HST lookups (`countryCode=CAN`)
- Historical rate lookups (`historical=YYYYMM`)

> Product taxability codes, Canadian rates, and historical data are plan-level entitlements — see [zip.tax pricing](https://zip.tax/pricing).

## Prerequisites

- Node.js 18+ (uses the built-in `fetch` API — no dependencies required)
- An API key from [zip.tax](https://www.zip.tax/)

## Run

```bash
export ZIPTAX_API_KEY=your_api_key_here
node example.js
```

## Sample Output

```
=== Address lookup (extended details) ===
  Normalized Address: 200 Spectrum Center Dr, Irvine, CA 92618-5003, United States
  Lat/Lng: 33.65253, -117.74794
  Incorporated: true
  Sourcing: Destination Based Taxation (D)
  Services taxable: N | Freight taxable: N
  Jurisdiction rates:
    US_STATE_SALES_TAX       CA                   7.250%
    US_COUNTY_SALES_TAX      ORANGE               0.500%
    US_CITY_SALES_TAX        IRVINE               0.000%
    ...
  Total Base Sales Tax: 7.75%
  Total Base Use Tax: 7.75%
  Parsed components: 200 Spectrum Center Dr, Irvine, CA 92618-5003
  Shipping rule (CA): CONDITIONAL — exempt when separately stated: True
```

For the full API reference, see the [official documentation](https://developers.zip.tax).
