#!/usr/bin/env ruby

require "json"
require "octokit"
require "yaml"

require File.join(File.dirname(__FILE__), "github")

# can expand this list spliting with spaces. e.g %w(cluster-admin root webops)
STRING_LIST = %w[cluster-admin]

# Output the yaml file and the code if any of the strings in the STRING_LIST is
# present in the files.

# The below code recurse through each of hash and array and pattern match if
# any of the string is present in any of the key or value field.

def main(gh)
  pattern = Regexp.union(STRING_LIST)
  yaml_files_in_pr(gh).find_all { |file|
    if FileTest.exists?(file)
      hash = YAML.load_file(file)
      recurse(hash, pattern) do |path, value|
        line = "#{path}:\t#{value}"
        if pattern.match?(line)
          message = <<~EOF
            The YAML file below
            
            #{file}
            
            contain the code 

            #{line}

            which will grant the user escalated privileges.

            Please correct them and resubmit this PR.

          EOF

          gh.reject_pr(message)
          exit 1
        end
      end
    end
  }
end

def recurse(obj, pattern, current_path = [], &block)
  if obj.is_a?(String)
    path = current_path.join(".")
    if obj =~ pattern || path =~ pattern
      yield [path, obj]
    end
  elsif obj.is_a?(Hash)
    obj.each do |key, value|
      recurse(value, pattern, current_path + [key], &block)
    end
  elsif obj.is_a?(Array)
    obj.each do |value|
      recurse(value, pattern, current_path, &block)
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

main(gh)
