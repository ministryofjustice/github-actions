#!/usr/bin/env ruby

require "json"
require "octokit"

require File.join(File.dirname(__FILE__), "github")

NAMESPACE_REGEX = %r[namespaces.live-1.cloud-platform.service.justice.gov.uk]

def namespaces_touched_by_pr
  files_in_pr
    .grep(NAMESPACE_REGEX)
    .map { |f| File.dirname(f) }
    .map { |f| f.split("/") }
    .map { |arr| arr[2] }
    .sort
    .uniq
end

############################################################

namespaces = namespaces_touched_by_pr

# PRs which touch no namespaces are fine
# PRs which touch one namespace are fine
if namespaces.size > 1
  namespace_list = namespaces.map {|n| "  * #{n}"}.join("\n")

  message = <<~EOF
  This PR affects multiple namespaces

  #{namespace_list}

  Please submit a separate PR for each namespace.

  EOF

  reject_pr(message)
end
