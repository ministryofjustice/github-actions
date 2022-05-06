plugin "azurerm" {
    enabled = true
    version = "0.16.0"
    source  = "github.com/terraform-linters/tflint-ruleset-azurerm"
}

rule "terraform_comment_syntax" {
  enabled = false
}

rule "terraform_documented_outputs" {
  enabled = true
}

rule "terraform_documented_variables" {
  enabled = true
}

rule "terraform_module_pinned_source" {
  enabled = false
  style = "flexible"
  default_branches = ["dev"]
}

rule "terraform_module_version" {
  enabled = true
  exact = false # default
}

rule "terraform_naming_convention" {
  enabled = true
  format = "snake_case"
}

rule "terraform_required_providers" {
  enabled = false
}

rule "terraform_required_version" {
  enabled = true
}

rule "terraform_standard_module_structure" {
  enabled = false
}

rule "terraform_typed_variables" {
  enabled = false
}

rule "terraform_unused_declarations" {
  enabled = true
}

rule "terraform_unused_required_providers" {
  enabled = true
}
