import play.api.libs.json._
import scala.io.Source
import scala.util.{Failure, Success, Try}
import java.net.URLEncoder
import java.nio.charset.StandardCharsets
import scalaj.http.Http

case class Response(
  version: String,
  rCode: Int,
  results: Seq[Result],
  addressDetail: AddressDetail
)

case class Result(
  geoPostalCode: String,
  geoCity: String,
  geoCounty: String,
  geoState: String,
  taxSales: Double,
  taxUse: Double,
  txbService: String,
  txbFreight: String,
  stateSalesTax: Double,
  stateUseTax: Double,
  citySalesTax: Double,
  cityUseTax: Double,
  cityTaxCode: String,
  countySalesTax: Double,
  countyUseTax: Double,
  countyTaxCode: String,
  districtSalesTax: Double,
  districtUseTax: Double,
  district1Code: String,
  district1SalesTax: Double,
  district1UseTax: Double,
  district2Code: String,
  district2SalesTax: Double,
  district2UseTax: Double,
  district3Code: String,
  district3SalesTax: Double,
  district3UseTax: Double,
  district4Code: String,
  district4SalesTax: Double,
  district4UseTax: Double,
  district5Code: String,
  district5SalesTax: Double,
  district5UseTax: Double,
  originDestination: String
)

case class AddressDetail(
  normalizedAddress: String,
  incorporated: String,
  geoLat: Double,
  geoLng: Double
)

object ZipTaxClient {

  implicit val resultReads: Reads[Result] = Json.reads[Result]
  implicit val addressDetailReads: Reads[AddressDetail] = Json.reads[AddressDetail]
  implicit val responseReads: Reads[Response] = Json.reads[Response]

  def getSalesTax(address: String, apiKey: String): Either[String, Response] = {
    val encodedAddress = URLEncoder.encode(address, StandardCharsets.UTF_8.toString)
    val url = s"https://api.zip-tax.com/request/v50?key=$apiKey&address=$encodedAddress"

    Try(Http(url).asString) match {
      case Success(response) if response.code == 200 =>
        Json.parse(response.body).validate[Response] match {
          case JsSuccess(taxResponse, _) => Right(taxResponse)
          case JsError(errors)           => Left(s"Failed to parse JSON: $errors")
        }
      case Success(response) => Left(s"Unexpected status code: ${response.code}")
      case Failure(exception) => Left(s"Error making API request: ${exception.getMessage}")
    }
  }

  def main(args: Array[String]): Unit = {
    val apiKey = "your_api_key_here" // Replace with your key
    val address = "200 Spectrum Center Dr, Irvine, CA 92618" // Example address

    getSalesTax(address, apiKey) match {
      case Right(taxInfo) =>
        println(s"Normalized Address: ${taxInfo.addressDetail.normalizedAddress}")
        println(s"Address Lat/Lng: ${taxInfo.addressDetail.geoLat}, ${taxInfo.addressDetail.geoLng}")
        if (taxInfo.results.nonEmpty)
          println(f"Rate: ${taxInfo.results.head.taxSales * 100}%.2f%%")
        else
          println("No tax information available.")

      case Left(error) =>
        println(s"Error fetching sales tax: $error")
    }
  }
}
