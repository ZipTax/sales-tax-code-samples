import requests
import json

class Response:
    def __init__(self, data):
        self.version = data.get("version")
        self.r_code = data.get("rCode")
        self.results = [Result(result) for result in data.get("results", [])]
        self.address_detail = AddressDetail(data.get("addressDetail", {}))

class Result:
    def __init__(self, data):
        self.geo_postal_code = data.get("geoPostalCode")
        self.geo_city = data.get("geoCity")
        self.geo_county = data.get("geoCounty")
        self.geo_state = data.get("geoState")
        self.tax_sales = data.get("taxSales")
        self.tax_use = data.get("taxUse")
        self.txb_service = data.get("txbService")
        self.txb_freight = data.get("txbFreight")
        self.state_sales_tax = data.get("stateSalesTax")
        self.state_use_tax = data.get("stateUseTax")
        self.city_sales_tax = data.get("citySalesTax")
        self.city_use_tax = data.get("cityUseTax")
        self.city_tax_code = data.get("cityTaxCode")
        self.county_sales_tax = data.get("countySalesTax")
        self.county_use_tax = data.get("countyUseTax")
        self.county_tax_code = data.get("countyTaxCode")
        self.district_sales_tax = data.get("districtSalesTax")
        self.district_use_tax = data.get("districtUseTax")
        self.district1_code = data.get("district1Code")
        self.district1_sales_tax = data.get("district1SalesTax")
        self.district1_use_tax = data.get("district1UseTax")
        self.district2_code = data.get("district2Code")
        self.district2_sales_tax = data.get("district2SalesTax")
        self.district2_use_tax = data.get("district2UseTax")
        self.district3_code = data.get("district3Code")
        self.district3_sales_tax = data.get("district3SalesTax")
        self.district3_use_tax = data.get("district3UseTax")
        self.district4_code = data.get("district4Code")
        self.district4_sales_tax = data.get("district4SalesTax")
        self.district4_use_tax = data.get("district4UseTax")
        self.district5_code = data.get("district5Code")
        self.district5_sales_tax = data.get("district5SalesTax")
        self.district5_use_tax = data.get("district5UseTax")
        self.origin_destination = data.get("originDestination")

class AddressDetail:
    def __init__(self, data):
        self.normalized_address = data.get("normalizedAddress")
        self.incorporated = data.get("incorporated")
        self.geo_lat = data.get("geoLat")
        self.geo_lng = data.get("geoLng")

def get_sales_tax(address, api_key):
    try:
        api_url = f"https://api.zip-tax.com/request/v50?key={api_key}&address={requests.utils.quote(address)}"
        response = requests.get(api_url)
        
        if response.status_code != 200:
            raise Exception(f"Unexpected status code: {response.status_code}")

        response_data = response.json()
        return Response(response_data)
    except Exception as e:
        print(f"Error fetching sales tax: {e}")
        return None

def main():
    api_key = "your_api_key_here"  # Replace with your actual API key
    address = "200 Spectrum Center Dr, Irvine, CA 92618"  # Example Address

    tax_info = get_sales_tax(address, api_key)

    if tax_info:
        print(f"Normalized Address: {tax_info.address_detail.normalized_address}")
        print(f"Address Lat/Lng: {tax_info.address_detail.geo_lat}, {tax_info.address_detail.geo_lng}")
        if tax_info.results:
            print(f"Rate: {tax_info.results[0].tax_sales * 100:.2f}%")

if __name__ == "__main__":
    main()