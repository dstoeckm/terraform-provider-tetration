provider "tetration" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://acme.tetrationpreview.com"
  disable_tls_verification = false
}

resource "tetration_user" "user_one" {
  enable_existing = true
  email           = "joe+100@acme.com"
  first_name      = "Joe"
  last_name       = "Bloggs 100"
  app_scope_id    = "5ce71503497d4f2c23af85b7"
}

resource "tetration_user" "user_two" {
  enable_existing = true
  email           = "joe+200@acme.com"
  first_name      = "Joe"
  last_name       = "Bloggs 200"
  app_scope_id    = "5ce71503497d4f2c23af85b7"
}

resource "tetration_scope" "scope" {
  short_name          = "Terraform created scope"
  short_query_type    = "eq"
  short_query_field   = "ip"
  short_query_value   = "192.168.0.1"
  parent_app_scope_id = "5ce71503497d4f2c23af85b7"
}

resource "tetration_role" "read_role" {
  name                = "read_role"
  app_scope_id        = "5ce71503497d4f2c23af85b7"
  access_app_scope_id = tetration_scope.scope.id
  access_type         = "scope_read"
  user_ids            = [tetration_user.user_one.id, tetration_user.user_two.id]
  description         = "role which provides read-only access to role_your_own_application"
}

resource "tetration_role" "dev_role" {
  name                = "dev_role"
  app_scope_id        = "5ce71503497d4f2c23af85b7"
  access_app_scope_id = tetration_scope.scope.id
  access_type         = "developer"
  user_ids            = [tetration_user.user_two.id]
  description         = "role which provides developer access to role_your_own_application"
}
