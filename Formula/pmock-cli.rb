class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/haidaisiganbei/pmock-cli"
  version "0.1.6"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "7771154c5d7451ba4c3d797c7eb69e2c184df49d45caf6127a52ecbdaffde481"
    else
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "109913a5169fb6c0c633aa42b84055c9216095966fa08fc5e11f5c1d3fa97cea"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "4a89ae87d6cd79b11c8cd2ebd2c6c031bfc7dc8965262b139b296224a4e3091f"
    else
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "a8f6b91c75f68f23af86066a12939e592dabca1cd2585113547e234cb7978e3b"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.6", shell_output("#{bin}/pmock-cli --version")
  end
end
