#!/usr/bin/env ruby

require "json"
require "octokit"

require File.join(File.dirname(__FILE__), "github")
require File.join(File.dirname(__FILE__), "code_formatter")

def main
  format_terraform_code
  format_ruby_code
  commit_changes "Commit changes made by code formatters"
end

main
