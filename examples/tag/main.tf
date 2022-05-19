provider "tetration" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://acme.tetrationpreview.com/"
  disable_tls_verification = false
}

resource "tetration_tag" "tag" {
  tenant_name = "acme"
  ip          = "10.0.0.1"
  attributes = {
    Environment = "test"
    Datacenter  = "aws"
    app_name    = "product-service"
  }
}
