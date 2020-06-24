#!/usr/bin/env ruby

require "json"
require "octokit"
require "yaml"
require 'colorize'

require File.join(File.dirname(__FILE__), "github")

PATTERN = "cluster-admin"

def yaml_files(gh)
  yaml_files_in_pr(gh).each do |file|
    hash = YAML.load_file(file) 
    pattern_text = Regexp.new(PATTERN, :nocase)
    recurse(hash, pattern) do |path, value|
      line = "#{path}:\t#{value}"
      line = line.gsub(pattern) {|match| match.green }
      puts line
    end
end

def recurse(obj, pattern, current_path = [], &block)
  if obj.is_a?(String)
    path = current_path.join('.')
    if obj =~ pattern || path =~ pattern
      yield [path, obj]
    end
  elsif obj.is_a?(Hash)
    obj.each do |key, value|
      recurse(value, pattern, current_path + [key], &block)
    end
  end
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

############################################################

gh = GithubClient.new

privileges_code = yaml_files(gh)

  message = <<~EOF
    The YAML files contain below code which will grant the user escalated privileges:

    #{file_list}

    Please correct them and resubmit this PR.

  EOF

  gh.reject_pr(message)
  exit 1
end


