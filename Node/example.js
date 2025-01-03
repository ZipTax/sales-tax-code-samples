const axios = require('axios');

async function getSalesTax(address, apiKey) {
    const url = `https://api.zip-tax.com/request/v50?key=${apiKey}&address=${encodeURIComponent(address)}`;

    try {
        const response = await axios.get(url);
        if (response.status !== 200) {
            throw new Error(`Unexpected status code: ${response.status}`);
        }

        return response.data;
    } catch (error) {
        throw new Error(`Failed to fetch sales tax: ${error.message}`);
    }
}

async function main() {
    const apiKey = 'your_api_key_here'; // Replace with your key
    const address = '200 Spectrum Center Dr, Irvine, CA 92618'; // Example address

    try {
        const taxInfo = await getSalesTax(address, apiKey);

        console.log(`Normalized Address: ${taxInfo.addressDetail.normalizedAddress}`);
        console.log(`Address Lat/Lng: ${taxInfo.addressDetail.geoLat}, ${taxInfo.addressDetail.geoLng}`);
        console.log(`Rate: ${(taxInfo.results[0].taxSales * 100).toFixed(2)}%`);
    } catch (error) {
        console.error(`Error fetching sales tax: ${error.message}`);
    }
}

main();