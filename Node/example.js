/**
 * zip.tax API v60 examples.
 *
 * Demonstrates the v60 endpoint features:
 *   - Address lookup with normalized address and geocoding
 *   - Lookup by geographic coordinates (lat/lng)
 *   - Structured jurisdiction rates (baseRates) and tax summaries
 *   - Sourcing rules and service/freight taxability
 *   - Extended address components and shipping rules
 *   - Product taxability codes (TIC) via productDetail
 *   - Canadian tax rates (GST/PST/HST)
 *
 * Uses the built-in fetch API (Node.js 18+). No dependencies required.
 * Set your API key in the ZIPTAX_API_KEY environment variable.
 */

const BASE_URL = 'https://api.zip-tax.com/request/v60';

async function requestV60(apiKey, params) {
    const query = new URLSearchParams({ key: apiKey, ...params });
    const response = await fetch(`${BASE_URL}?${query}`);

    if (!response.ok) {
        throw new Error(`Unexpected status code: ${response.status}`);
    }

    const data = await response.json();
    const { code, message } = data.metadata?.response ?? {};
    if (code !== 100) {
        throw new Error(`API error ${code}: ${message ?? 'unknown error'}`);
    }

    return data;
}

function printSummary(data) {
    const detail = data.addressDetail;
    console.log(`  Normalized Address: ${detail.normalizedAddress}`);
    console.log(`  Lat/Lng: ${detail.geoLat}, ${detail.geoLng}`);
    console.log(`  Incorporated: ${detail.incorporated}`);

    const sourcing = data.sourcingRules;
    if (sourcing) {
        console.log(`  Sourcing: ${sourcing.description} (${sourcing.value})`);
    }
    if (data.service) {
        console.log(`  Services taxable: ${data.service.taxable} | Freight taxable: ${data.shipping.taxable}`);
    }

    console.log('  Jurisdiction rates:');
    for (const rate of data.baseRates) {
        console.log(`    ${rate.jurType.padEnd(24)} ${rate.jurName.padEnd(20)} ${(rate.rate * 100).toFixed(3)}%`);
    }

    for (const summary of data.taxSummaries) {
        console.log(`  ${summary.summaryName}: ${(summary.rate * 100).toFixed(2)}%`);
    }
}

async function main() {
    const apiKey = process.env.ZIPTAX_API_KEY ?? 'your_api_key_here';

    // 1. Door-level rate lookup by street address, with extended address
    //    components and state shipping rules.
    console.log('=== Address lookup (extended details) ===');
    let data = await requestV60(apiKey, {
        address: '200 Spectrum Center Dr, Irvine, CA 92618',
        addressDetailExtended: 'true',
        shippingExtended: 'true',
    });
    printSummary(data);

    const components = data.addressDetail.address;
    if (components) {
        console.log(`  Parsed components: ${components.houseNumber} ${components.street}, ` +
            `${components.city}, ${components.stateCode} ${components.postalCode}`);
    }

    const shippingExt = data.shipping.shippingExtended;
    if (shippingExt) {
        console.log(`  Shipping rule (${shippingExt.stateCode}): ${shippingExt.rule} — ` +
            `exempt when separately stated: ${shippingExt.exemptWhenSeparatelyStated}`);
    }

    // 2. Lookup by geographic coordinates.
    console.log('\n=== Coordinate lookup (lat/lng) ===');
    data = await requestV60(apiKey, { lat: 33.65253, lng: -117.74794 });
    printSummary(data);

    // 3. Product-specific tax rules via a Taxability Information Code (TIC).
    //    TIC 40030 = food and food ingredients. Requires the product_rates
    //    plan entitlement.
    console.log('\n=== Product taxability code (TIC 40030 — food) ===');
    data = await requestV60(apiKey, {
        address: '100 Broadway, Nashville, TN 37201',
        taxabilityCode: '40030',
    });
    const product = data.productDetail?.taxabilityCode;
    if (product) {
        console.log(`  TIC ${product.id}: ${product.title}`);
        console.log(`  Rate action: ${product.rateActionCode} — ${product.rateActionMessage}`);
        for (const rule of product.rateRules) {
            console.log(`    Jurisdiction ${rule.jurTaxCode}: effective ${rule.effectiveDt}, ` +
                `percent taxable ${rule.percentTaxable}%, food/drug: ${rule.isFoodDrug}`);
        }
    }

    // 4. Canadian tax rates. Requires the rate_loc_can plan entitlement.
    console.log('\n=== Canadian lookup (Toronto, ON) ===');
    data = await requestV60(apiKey, { countryCode: 'CAN', postalcode: 'M5V 3L9' });
    for (const rate of data.baseRates) {
        console.log(`  ${rate.jurType.padEnd(4)} ${rate.jurName.padEnd(10)} ${(rate.rate * 100).toFixed(2)}% — ${rate.jurDescription}`);
    }
    for (const summary of data.taxSummaries) {
        for (const display of summary.displayRates) {
            console.log(`  ${summary.summaryName}: ${display.name} ${(display.rate * 100).toFixed(2)}%`);
        }
    }

    // 5. Historical rates for a past period (YYYYMM). Requires the
    //    historical data plan entitlement.
    // data = await requestV60(apiKey, { address: '200 Spectrum Center Dr, Irvine, CA 92618', historical: '202401' });
}

main().catch((error) => {
    console.error(`Error: ${error.message}`);
    process.exit(1);
});
