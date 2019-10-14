#!/usr/bin/env ruby

require "json"
require "octokit"

def github
  unless ENV.key?("GITHUB_TOKEN")
    raise "No GITHUB_TOKEN env var found. Please make this available via the github actions workflow\nhttps://help.github.com/en/articles/virtual-environments-for-github-actions#github_token-secret"
  end

  @client ||= Octokit::Client.new(access_token: ENV["GITHUB_TOKEN"])
end

def event
  unless ENV.key?("GITHUB_EVENT_PATH")
    raise "No GITHUB_EVENT_PATH env var found. This script is designed to run via github actions, which will provide the github event via this env var."
  end

  @evt ||= JSON.parse File.read(ENV["GITHUB_EVENT_PATH"])
end

def repo
  name = event.dig("repository", "name")
  owner = event.dig("repository", "owner", "login")
  [owner, name].join("/")
end

def pr_number
  event.dig("pull_request", "number")
end

def branch
  event.dig("pull_request", "head", "ref")
end

def commit_files(branch, files, commit_message)
  ref = "heads/#{branch}"
  sha_latest_commit = github.ref(repo, ref).object.sha
  sha_base_tree = github.commit(repo, sha_latest_commit).commit.tree.sha
  changes = create_blobs files
  sha_new_tree = github.create_tree(repo, changes, {:base_tree => sha_base_tree }).sha
  sha_new_commit = github.create_commit(repo, commit_message, sha_new_tree, sha_latest_commit).sha
  github.update_ref(repo, ref, sha_new_commit)
end

def create_blobs(files)
  files.map do |file_name|
    content = File.read(file_name)
    blob_sha = github.create_blob(repo, Base64.encode64(content), "base64")
    { :path => file_name, :mode => "100644", :type => "blob", :sha => blob_sha }
  end
end

def format_terraform_code
  terraform_directories_in_pr.each do |dir|
    execute "terraform fmt #{dir}"
  end
  commit_changes "Commit changes made by `terraform fmt`"
end

def format_ruby_code
  ruby_files_in_pr.each do |file|
    execute "standardrb --fix #{file}"
  end
  commit_changes "Commit changes made by `standardrb --fix`"
end

def terraform_directories_in_pr
  terraform_files_in_pr
    .map { |f| File.dirname(f) }
    .sort
    .uniq
end

def ruby_files_in_pr
  files_in_pr.grep(/\.rb$/)
end

def terraform_files_in_pr
  files_in_pr.grep(/\.tf$/)
end

def files_in_pr
  github.pull_request_files(repo, pr_number)
    .map(&:filename)
    .sort
    .uniq
end

def execute(cmd)
  puts "Running: #{cmd}"
  `#{cmd}`
end

def commit_changes(message)
  files = modified_files
  if files.any?
    puts "Committing changes to:\n  #{files.join("\n  ")}"
    commit_files(branch, files, message)
  end
end

def modified_files
  execute("git status --porcelain=1 --untracked-files=no")
    .split("\n")
    .map { |line| line.sub(" M ", "") }
end

############################################################

format_terraform_code
format_ruby_code
