#!/usr/bin/env ruby

require "json"
require "octokit"

require File.join(File.dirname(__FILE__), "github")

def format_terraform_code
  terraform_directories_in_pr.each do |dir|
    terraform11?(dir) && format_terraform11(dir)
    terraform12?(dir) && format_terraform12(dir)
  end
end

def format_terraform11
  execute "terraform fmt #{dir}"
  _stdout, stderr, status = execute "terraform validate -check-variables=false #{dir}"
  raise "terraform validate failed:\n#{stderr}" unless status.success?
end

def format_terraform12
  execute "terraform12 fmt #{dir}"
  _stdout, stderr, status = execute "terraform12 init && terraform12 validate #{dir}"
  raise "terraform12 validate failed:\n#{stderr}" unless status.success?
end

def terraform11?(dir)
  FileTest.directory?(dir) && !FileTest.exists?(File.join(dir, "versions.tf"))
end

def terraform12?(dir)
  FileTest.directory?(dir) && FileTest.exists?(File.join(dir, "versions.tf"))
end

def format_ruby_code
  ruby_files_in_pr.each do |file|
    execute "standardrb --fix #{file}" if FileTest.exists?(file)
  end
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

############################################################

format_terraform_code
format_ruby_code
commit_changes "Commit changes made by code formatters"
