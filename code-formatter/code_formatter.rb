class CodeFormatter
  attr_reader :executor, :github_client

  def initialize(args = {})
    @executor = args.fetch(:executor) { Executor.new }
    @github_client = args.fetch(:github_client) { GithubClient.new(executor: executor) }
  end

  def run
    format_terraform_code
    format_ruby_code
    github_client.commit_changes "Commit changes made by code formatters"
  end

  private

  def format_terraform_code
    terraform_directories_in_pr.each do |dir|
      executor.execute "terraform fmt #{dir}"
    end
  end

  def format_ruby_code
    ruby_files_in_pr.each do |file|
      executor.execute "standardrb --fix #{file}" if FileTest.exists?(file)
    end
  end

  def terraform_directories_in_pr
    terraform_files_in_pr
      .map { |f| File.dirname(f) }
      .sort
      .uniq
  end

  def ruby_files_in_pr
    github_client.files_in_pr.grep(/\.rb$/)
  end

  def terraform_files_in_pr
    github_client.files_in_pr.grep(/\.tf$/)
  end
end
