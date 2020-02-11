#!/usr/bin/env ruby

require "json"
require "octokit"
require "yaml"

require File.join(File.dirname(__FILE__), "github")

def malformed_yaml_files(gh)
  yaml_files_in_pr(gh).find_all { |file| fails_to_parse?(file) }
end

# Attempt to parse all the yaml/yml files in a PR,
# aside from those with 'secret' in the filename.
# Files with 'secret' in the name are very often
# git-crypted, and so would cause this action to
# fail.
def yaml_files_in_pr(gh)
  gh.files_in_pr
    .grep(/\.(yaml|yml)$/)
    .reject { |f| f =~ /secret/ }
end

def fails_to_parse?(file)
  YAML.safe_load File.read(file) if FileTest.exists?(file)
  false
rescue Psych::SyntaxError
  true
end

############################################################

gh = GithubClient.new

files = malformed_yaml_files(gh)

if files.any?
  file_list = files.map { |f| "  * #{f}" }.join("\n")

  message = <<~EOF
    The following files contain malformed YAML:

    #{file_list}

    Please correct them and resubmit this PR.

  EOF

  gh.reject_pr(message)
  exit 1
end
