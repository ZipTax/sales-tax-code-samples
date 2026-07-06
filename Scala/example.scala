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
 * Dependencies: play-json, scalaj-http.
 * Set your API key in the ZIPTAX_API_KEY environment variable.
 */

import play.api.libs.json._
import scala.util.{Failure, Success, Try}
import scalaj.http.Http

case class ResponseInfo(code: Int, name: String, message: String, definition: String)
case class Metadata(version: String, response: ResponseInfo)

case class BaseRate(
  rate: Double,
  jurType: String,
  jurName: String,
  jurDescription: String,
  jurTaxCode: Option[String]
)

case class ServiceConfig(adjustmentType: String, taxable: String, description: String)

case class ShippingExtended(
  stateName: String,
  stateCode: String,
  rule: String,
  exemptWhenSeparatelyStated: String,
  description: String
)

case class ShippingConfig(
  adjustmentType: String,
  taxable: String,
  description: String,
  shippingExtended: Option[ShippingExtended]
)

case class SourcingRules(adjustmentType: String, description: String, value: String)

case class DisplayRate(name: String, rate: Double)
case class TaxSummary(rate: Double, taxType: String, summaryName: String, displayRates: Seq[DisplayRate])

case class RateRule(
  jurTaxCode: Option[String],
  effectiveDt: Option[Long],
  expiresDt: Option[Long],
  effectiveTaxRate: Option[Double],
  percentTaxable: Option[Double],
  exemptUnder: Option[Double],
  exemptOver: Option[Double],
  taxablePortionOver: Option[Double],
  isDestinationTaxType: Option[Boolean],
  isFoodDrug: Option[Boolean]
)

case class TaxabilityCode(
  id: String,
  stateFIPS: String,
  countyFIPS: String,
  title: String,
  label: String,
  rateActionCode: String,
  rateActionMessage: String,
  rateRules: Seq[RateRule]
)

case class ProductDetail(taxabilityCode: TaxabilityCode)

case class AddressComponents(
  countryCode: String,
  countryName: String,
  stateCode: String,
  state: String,
  county: String,
  city: String,
  street: String,
  postalCode: String,
  houseNumber: String
)

case class AddressDetail(
  normalizedAddress: String,
  incorporated: String,
  geoLat: Double,
  geoLng: Double,
  address: Option[AddressComponents]
)

case class Response(
  metadata: Metadata,
  baseRates: Seq[BaseRate],
  service: Option[ServiceConfig],
  shipping: ShippingConfig,
  sourcingRules: Option[SourcingRules],
  taxSummaries: Seq[TaxSummary],
  productDetail: Option[ProductDetail],
  addressDetail: AddressDetail
)

object ZipTaxClient {

  implicit val responseInfoReads: Reads[ResponseInfo] = Json.reads[ResponseInfo]
  implicit val metadataReads: Reads[Metadata] = Json.reads[Metadata]
  implicit val baseRateReads: Reads[BaseRate] = Json.reads[BaseRate]
  implicit val serviceConfigReads: Reads[ServiceConfig] = Json.reads[ServiceConfig]
  implicit val shippingExtendedReads: Reads[ShippingExtended] = Json.reads[ShippingExtended]
  implicit val shippingConfigReads: Reads[ShippingConfig] = Json.reads[ShippingConfig]
  implicit val sourcingRulesReads: Reads[SourcingRules] = Json.reads[SourcingRules]
  implicit val displayRateReads: Reads[DisplayRate] = Json.reads[DisplayRate]
  implicit val taxSummaryReads: Reads[TaxSummary] = Json.reads[TaxSummary]
  implicit val rateRuleReads: Reads[RateRule] = Json.reads[RateRule]
  implicit val taxabilityCodeReads: Reads[TaxabilityCode] = Json.reads[TaxabilityCode]
  implicit val productDetailReads: Reads[ProductDetail] = Json.reads[ProductDetail]
  implicit val addressComponentsReads: Reads[AddressComponents] = Json.reads[AddressComponents]
  implicit val addressDetailReads: Reads[AddressDetail] = Json.reads[AddressDetail]
  implicit val responseReads: Reads[Response] = Json.reads[Response]

  private val BaseUrl = "https://api.zip-tax.com/request/v60"

  def requestV60(apiKey: String, params: Map[String, String]): Either[String, Response] = {
    val request = params.foldLeft(Http(BaseUrl).param("key", apiKey)) {
      case (req, (name, value)) => req.param(name, value)
    }

    Try(request.asString) match {
      case Success(response) if response.code == 200 =>
        Json.parse(response.body).validate[Response] match {
          case JsSuccess(taxResponse, _) =>
            if (taxResponse.metadata.response.code == 100) Right(taxResponse)
            else Left(s"API error ${taxResponse.metadata.response.code}: ${taxResponse.metadata.response.message}")
          case JsError(errors) => Left(s"Failed to parse JSON: $errors")
        }
      case Success(response)  => Left(s"Unexpected status code: ${response.code}")
      case Failure(exception) => Left(s"Error making API request: ${exception.getMessage}")
    }
  }

  def printSummary(data: Response): Unit = {
    val detail = data.addressDetail
    println(s"  Normalized Address: ${detail.normalizedAddress}")
    println(s"  Lat/Lng: ${detail.geoLat}, ${detail.geoLng}")
    println(s"  Incorporated: ${detail.incorporated}")

    data.sourcingRules.foreach(s => println(s"  Sourcing: ${s.description} (${s.value})"))
    data.service.foreach(s => println(s"  Services taxable: ${s.taxable} | Freight taxable: ${data.shipping.taxable}"))

    println("  Jurisdiction rates:")
    data.baseRates.foreach { rate =>
      println(f"    ${rate.jurType}%-24s ${rate.jurName}%-20s ${rate.rate * 100}%.3f%%")
    }

    data.taxSummaries.foreach { summary =>
      println(f"  ${summary.summaryName}: ${summary.rate * 100}%.2f%%")
    }
  }

  def main(args: Array[String]): Unit = {
    val apiKey = sys.env.getOrElse("ZIPTAX_API_KEY", "your_api_key_here")

    // 1. Door-level rate lookup by street address, with extended address
    //    components and state shipping rules.
    println("=== Address lookup (extended details) ===")
    requestV60(apiKey, Map(
      "address"               -> "200 Spectrum Center Dr, Irvine, CA 92618",
      "addressDetailExtended" -> "true",
      "shippingExtended"      -> "true"
    )) match {
      case Right(data) =>
        printSummary(data)
        data.addressDetail.address.foreach { c =>
          println(s"  Parsed components: ${c.houseNumber} ${c.street}, ${c.city}, ${c.stateCode} ${c.postalCode}")
        }
        data.shipping.shippingExtended.foreach { ext =>
          println(s"  Shipping rule (${ext.stateCode}): ${ext.rule} — exempt when separately stated: ${ext.exemptWhenSeparatelyStated}")
        }
      case Left(error) => println(s"Error: $error")
    }

    // 2. Lookup by geographic coordinates.
    println("\n=== Coordinate lookup (lat/lng) ===")
    requestV60(apiKey, Map("lat" -> "33.65253", "lng" -> "-117.74794")) match {
      case Right(data) => printSummary(data)
      case Left(error) => println(s"Error: $error")
    }

    // 3. Product-specific tax rules via a Taxability Information Code (TIC).
    //    TIC 40030 = food and food ingredients. Requires the product_rates
    //    plan entitlement.
    println("\n=== Product taxability code (TIC 40030 — food) ===")
    requestV60(apiKey, Map(
      "address"        -> "100 Broadway, Nashville, TN 37201",
      "taxabilityCode" -> "40030"
    )) match {
      case Right(data) =>
        data.productDetail.foreach { detail =>
          val product = detail.taxabilityCode
          println(s"  TIC ${product.id}: ${product.title}")
          println(s"  Rate action: ${product.rateActionCode} — ${product.rateActionMessage}")
          product.rateRules.foreach { rule =>
            println(s"    Jurisdiction ${rule.jurTaxCode.getOrElse("n/a")}: " +
              s"effective ${rule.effectiveDt.getOrElse(0L)}, " +
              s"percent taxable ${rule.percentTaxable.getOrElse(0.0)}%, " +
              s"food/drug: ${rule.isFoodDrug.getOrElse(false)}")
          }
        }
      case Left(error) => println(s"Error: $error")
    }

    // 4. Canadian tax rates. Requires the rate_loc_can plan entitlement.
    println("\n=== Canadian lookup (Toronto, ON) ===")
    requestV60(apiKey, Map("countryCode" -> "CAN", "postalcode" -> "M5V 3L9")) match {
      case Right(data) =>
        data.baseRates.foreach { rate =>
          println(f"  ${rate.jurType}%-4s ${rate.jurName}%-10s ${rate.rate * 100}%.2f%% — ${rate.jurDescription}")
        }
        data.taxSummaries.foreach { summary =>
          summary.displayRates.foreach { display =>
            println(f"  ${summary.summaryName}: ${display.name} ${display.rate * 100}%.2f%%")
          }
        }
      case Left(error) => println(s"Error: $error")
    }

    // 5. Historical rates for a past period (YYYYMM). Requires the
    //    historical data plan entitlement.
    // requestV60(apiKey, Map("address" -> "200 Spectrum Center Dr, Irvine, CA 92618", "historical" -> "202401"))
  }
}
