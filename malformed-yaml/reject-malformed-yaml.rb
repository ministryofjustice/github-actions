#!/usr/bin/env ruby

require "json"
require "octokit"
require "yaml"

require File.join(File.dirname(__FILE__), "github")

def malformed_yaml_files
  yaml_files_in_pr.map do |file|
    YAML.load File.read(file) if FileTest.exists?(file)
    nil
  rescue Psych::SyntaxError
    file
  end.compact
end

def yaml_files_in_pr
  files_in_pr.grep(/\.(yaml|yml)$/)
end

############################################################

files = malformed_yaml_files

if files.any?
  file_list = files.map {|f| "  * #{f}"}.join("\n")

  message = <<~EOF
  The following files contain malformed YAML:

  #{file_list}

  Please correct them and resubmit this PR.

  EOF

  reject_pr(message)
end
