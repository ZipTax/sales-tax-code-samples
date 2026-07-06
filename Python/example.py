"""zip.tax API v60 examples.

Demonstrates the v60 endpoint features:
  - Address lookup with normalized address and geocoding
  - Lookup by geographic coordinates (lat/lng)
  - Structured jurisdiction rates (baseRates) and tax summaries
  - Sourcing rules and service/freight taxability
  - Extended address components and shipping rules
  - Product taxability codes (TIC) via productDetail
  - Canadian tax rates (GST/PST/HST)

Set your API key in the ZIPTAX_API_KEY environment variable.
"""

import os

import requests

BASE_URL = "https://api.zip-tax.com/request/v60"


def request_v60(api_key, **params):
    """Make a v60 request and return the parsed JSON response."""
    params["key"] = api_key
    response = requests.get(BASE_URL, params=params, timeout=30)
    response.raise_for_status()
    data = response.json()

    code = data.get("metadata", {}).get("response", {}).get("code")
    if code != 100:
        message = data.get("metadata", {}).get("response", {}).get("message", "unknown error")
        raise RuntimeError(f"API error {code}: {message}")

    return data


def print_summary(data):
    """Print the key parts of a v60 response."""
    detail = data["addressDetail"]
    print(f"  Normalized Address: {detail['normalizedAddress']}")
    print(f"  Lat/Lng: {detail['geoLat']}, {detail['geoLng']}")
    print(f"  Incorporated: {detail['incorporated']}")

    sourcing = data.get("sourcingRules")
    if sourcing:
        print(f"  Sourcing: {sourcing['description']} ({sourcing['value']})")
    service = data.get("service")
    if service:
        print(f"  Services taxable: {service['taxable']} | Freight taxable: {data['shipping']['taxable']}")

    print("  Jurisdiction rates:")
    for rate in data["baseRates"]:
        print(f"    {rate['jurType']:<24} {rate['jurName']:<20} {rate['rate'] * 100:.3f}%")

    for summary in data["taxSummaries"]:
        print(f"  {summary['summaryName']}: {summary['rate'] * 100:.2f}%")


def main():
    api_key = os.environ.get("ZIPTAX_API_KEY", "your_api_key_here")

    # 1. Door-level rate lookup by street address, with extended address
    #    components and state shipping rules.
    print("=== Address lookup (extended details) ===")
    data = request_v60(
        api_key,
        address="200 Spectrum Center Dr, Irvine, CA 92618",
        addressDetailExtended="true",
        shippingExtended="true",
    )
    print_summary(data)

    components = data["addressDetail"].get("address")
    if components:
        print(f"  Parsed components: {components['houseNumber']} {components['street']}, "
              f"{components['city']}, {components['stateCode']} {components['postalCode']}")

    shipping_ext = data["shipping"].get("shippingExtended")
    if shipping_ext:
        print(f"  Shipping rule ({shipping_ext['stateCode']}): {shipping_ext['rule']} — "
              f"exempt when separately stated: {shipping_ext['exemptWhenSeparatelyStated']}")

    # 2. Lookup by geographic coordinates.
    print("\n=== Coordinate lookup (lat/lng) ===")
    data = request_v60(api_key, lat=33.65253, lng=-117.74794)
    print_summary(data)

    # 3. Product-specific tax rules via a Taxability Information Code (TIC).
    #    TIC 40030 = food and food ingredients. Requires the product_rates
    #    plan entitlement.
    print("\n=== Product taxability code (TIC 40030 — food) ===")
    data = request_v60(
        api_key,
        address="100 Broadway, Nashville, TN 37201",
        taxabilityCode="40030",
    )
    product = data.get("productDetail", {}).get("taxabilityCode")
    if product:
        print(f"  TIC {product['id']}: {product['title']}")
        print(f"  Rate action: {product['rateActionCode']} — {product['rateActionMessage']}")
        for rule in product["rateRules"]:
            print(f"    Jurisdiction {rule['jurTaxCode']}: "
                  f"effective {rule['effectiveDt']}, "
                  f"percent taxable {rule['percentTaxable']}%, "
                  f"food/drug: {rule['isFoodDrug']}")

    # 4. Canadian tax rates. Requires the rate_loc_can plan entitlement.
    print("\n=== Canadian lookup (Toronto, ON) ===")
    data = request_v60(api_key, countryCode="CAN", postalcode="M5V 3L9")
    for rate in data["baseRates"]:
        print(f"  {rate['jurType']:<4} {rate['jurName']:<10} {rate['rate'] * 100:.2f}% — {rate['jurDescription']}")
    for summary in data["taxSummaries"]:
        for display in summary["displayRates"]:
            print(f"  {summary['summaryName']}: {display['name']} {display['rate'] * 100:.2f}%")

    # 5. Historical rates for a past period (YYYYMM). Requires the
    #    historical data plan entitlement.
    # data = request_v60(api_key, address="200 Spectrum Center Dr, Irvine, CA 92618", historical="202401")


if __name__ == "__main__":
    main()
