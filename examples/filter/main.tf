provider "tetration" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://acme.tetrationpreview.com"
  disable_tls_verification = false
}

resource "tetration_filter" "filter" {
  name         = "Terraform created filter"
  query        = <<EOF
                    {
                      "type": "eq",
                      "field": "ip",
                      "value": "10.0.0.1"
                    }
          EOF
  app_scope_id = "5ed6890c497d4f55eb5c585c"
  primary      = true
  public       = false
}
