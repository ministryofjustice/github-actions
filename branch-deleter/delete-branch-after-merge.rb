#!/usr/bin/env ruby

require "json"
require "octokit"

require File.join(File.dirname(__FILE__), "github")

gh = GithubClient.new

puts "Merged PR: #{gh.pr_number}"
puts "Deleting branch: #{gh.branch}"

begin
  github.client.delete_branch(repo, branch)
rescue Octokit::UnprocessableEntity
  puts "Branch not found; already deleted?"
end
