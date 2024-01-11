plugin "aws" {
    enabled = true
    version = "0.29.0"
    source  = "github.com/terraform-linters/tflint-ruleset-aws"
}

plugin "terraform" {
    enabled = true
    version = "0.5.0"
    preset  = "recommended"
    source  = "github.com/terraform-linters/tflint-ruleset-terraform"
}

rule "terraform_deprecated_index" {
    enabled = false
}
