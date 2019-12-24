require "spec_helper"

describe CodeFormatter do
  let(:executor) { double(Executor) }
  let(:github_client) { double(GithubClient, files_in_pr: files, commit_changes: nil) }

  let(:params) {
    {
      executor: executor,
      github_client: github_client,
    }
  }

  subject(:cf) { described_class.new(params) }

  before do
    files.each { |f| allow(FileTest).to receive(:exists?).with(f).and_return(true) }
  end

  context "when PR updates ruby files" do
    let(:files) { ["foo.rb"] }

    it "formats ruby code" do
      expect(executor).to receive(:execute).with("standardrb --fix foo.rb")
      cf.run
    end
  end

  context "when PR updates terraform files" do
    let(:files) { ["resources/foo.tf", "anotherdir/bar.tf"] }

    it "formats terraform code directories" do
      expect(executor).to receive(:execute).with("terraform fmt resources")
      expect(executor).to receive(:execute).with("terraform fmt anotherdir")
      cf.run
    end
  end
end
