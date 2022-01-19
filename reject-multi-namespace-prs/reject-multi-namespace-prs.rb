#!/usr/bin/env ruby

require "json"
require "octokit"

require File.join(File.dirname(__FILE__), "github")

NAMESPACE_REGEX = %r{namespaces.(live|live-1).cloud-platform.service.justice.gov.uk}

gh = GithubClient.new

def namespaces_touched_by_pr(gh)
  gh.files_in_pr
    .grep(NAMESPACE_REGEX)
    .map { |f| File.dirname(f) }
    .map { |f| f.split("/") }
    .map { |arr| arr[2] }
    .sort
    .uniq
end

############################################################

namespaces = namespaces_touched_by_pr(gh)

# PRs which touch no namespaces are fine
# PRs which touch one namespace are fine
if namespaces.size > 1
  namespace_list = namespaces.map { |n| "  * #{n}" }.join("\n")

  message = <<~EOF
    This PR affects multiple namespaces
     #{namespace_list}
     Please submit a separate PR for each namespace.

  EOF

  gh.reject_pr(message)
  exit 1
end
