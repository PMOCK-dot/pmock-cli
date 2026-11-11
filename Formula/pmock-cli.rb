class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.20"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "7ff25535eb1d1aaa910d9372c678fa490ce4041f042f7ea337db2d80b900e55b"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "63621d2649f723eb6237e2391a64cddda7ef9f2def256879d5764ae6127adcf3"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "7de91d7ec4f1db94daf2293317313f9d3c621f64339759d5b6ad28616bed08dc"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "250e29b2baf750cb14adaac8bb64f49878dc83aa47d361bf805b7b4ba57f7823"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.20", shell_output("#{bin}/pmock-cli --version")
  end
end
