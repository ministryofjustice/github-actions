#!/usr/bin/env ruby

require "json"
require "open3"
require "octokit"

require File.join(File.dirname(__FILE__), "github")

def main
  test_policies
  check_yaml_files
end

# Run the rego tests for our policies
def test_policies
  # Assume rego policies are in the ./policy directory, unless user supplied a
  # different location
  policy_dir = ENV.fetch("POLICY_DIR", "policy")
  cmd = "opa test #{policy_dir}"
  exit 1 unless command_status(cmd)
end

# Test all the YAML files in this PR to ensure they comply with our Rego
# policies
def check_yaml_files
  client = GithubClient.new

  # Assume rego policies are in the ./policy directory, unless user supplied a
  # different location
  policy_dir = ENV.fetch("POLICY_DIR", "policy")

  # Get any additional command-line options for conftest
  conftest_options = ENV.fetch("CONFTEST_OPTIONS", "")

  # We want to test all files, rather than exiting on the first failure, so that
  # the user can see all problems reported in the log. So, we collect all the
  # exit statuses of the conftest commands.
  cmd_statuses = yaml_files_in_pr(client).map { |file|
    cmd = "conftest test -p #{policy_dir} #{conftest_options} #{file}"
    command_status(cmd)
  }

  # Fail the action if conftest failed any YAML files
  exit 1 unless cmd_statuses.all?
end

# Attempt to parse all the yaml/yml files in a PR, aside from those with
# 'secret' in the filename.  Files with 'secret' in the name are very often
# git-crypted, and so would cause this action to fail.
def yaml_files_in_pr(gh)
  gh.files_in_pr
    .grep(/\.(yaml|yml)$/)
    .filter { |f| FileTest.exists?(f) }
    .reject { |f| f =~ /secret/ }
end

def command_status(cmd)
  puts "CMD: #{cmd}"
  stdout, stderr, status = Open3.capture3(cmd)
  puts stdout
  puts "ERROR:\n#{stderr}" unless status.success?
  status.success?
end

main
