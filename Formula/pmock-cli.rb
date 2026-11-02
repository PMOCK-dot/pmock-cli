class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/haidaisiganbei/pmock-cli"
  version "0.1.7"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "86216932b82cb86ac901d8a6d37c86fd14c53e3a7d7b90d093ccbccaa0b9f3c7"
    else
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "175bfc883a03f848683f1f13625fd18ab598a8cb0ce64c26f96346b587aa4598"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "a15611370b4173d7101a23c771cc56dc2c5b63720e5ed16032bdc5ec71edeb8e"
    else
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "8fbf10070314cece9d34a42be7c611f66eac91dcaf1ab414c634c5ac620e0793"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.7", shell_output("#{bin}/pmock-cli --version")
  end
end
