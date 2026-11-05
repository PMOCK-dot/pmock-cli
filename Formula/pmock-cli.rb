class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.17"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "b53ba8ffb50325bf1b682e855678f14906cca463d1c530f7defc038264ce7add"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "d71ad2ec2b2d3884e9a76ac2321fd2ba267fc18f90ff860662e875fbaecc1454"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "8167486d46235b8746e4b7bca54fbd51c735f4797d877e2c9e7861593f110b68"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "9bd46c48f67892ccf6a5c0a022cb96d01d130b9a3d428b9101dccac52b0418d5"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.17", shell_output("#{bin}/pmock-cli --version")
  end
end
