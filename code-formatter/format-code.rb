#!/usr/bin/env ruby

require "json"
require "octokit"

require File.join(File.dirname(__FILE__), "github")
require File.join(File.dirname(__FILE__), "code_formatter")

CodeFormatter.new.run
