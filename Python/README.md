If you're building an application that requires accurate sales tax calculations, the [zip.tax API](https://developers.zip.tax) is an excellent tool to integrate. This guide walks you through how to set up and use the zip.tax API in a Python application.

## Prerequisites

Before getting started, ensure you have the following:

- Basic knowledge of Python.
- A Python development environment set up.
- An API key from [zip.tax](https://www.zip.tax/).

## Step 1: Install Required Libraries

For making HTTP requests, we'll use Python's standard `requests` package. Additionally, we'll use `json` for parsing JSON responses.

## Step 2: Set Up Your Python Project

Create a new project directory and initialize a new module:

```bash
mkdir ziptax-python && cd ziptax-python
```

## Step 3: Write the Code

Here is a complete example of a simple Python application that queries the zip.tax API for sales tax information.

```python
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
```

### Explanation of the Code

1. **API Request:** The `get_sales_tax` function constructs a URL with the API key and an address, makes a GET request, and parses the response.
2. **Response Parsing:** The response JSON is unmarshalled for easy access to sales tax details.
3. **Display Results:** The main function prints the normalized address, lat/lng, and sales tax rate for the specified address code. You can use any of the response values here to output the data you need. 

## Step 4: Run the Application

Save the code to a file (e.g., `main.py`), then run the program:

```bash
python main.py
```

You should see output similar to this:

```
Normalized Address: 200 Spectrum Center Dr, Irvine, CA 92618-5003, United States
Address Lat/Lng: 33.652530, -117.747940
Rate: 7.75%
```

## Conclusion

Integrating the zip.tax API into your Python application is straightforward. By following this guide, you can enhance your application with accurate sales tax information based on address. For more details, refer to the [official documentation](https://developers.zip.tax).

If you have any questions or feedback, feel free to leave a comment below. Happy coding!