plugin "aws" {
    enabled = true
    version = "0.17.0"
    source  = "github.com/terraform-linters/tflint-ruleset-aws"
}

plugin "terraform" {
    enabled = true
    version = "0.1.0"
    preset  = "recommended"
    source  = "github.com/terraform-linters/tflint-ruleset-terraform"
}

rule "terraform_deprecated_index" {
    enabled = false
}
