class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.15"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "f5d25d49b9e47a89b1a34ae7ceafcd8b048a73eaf246064bf5dc55bf862dbb0f"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "df65352b60e33ff1874e3a2dcf7f38ebc7fd5174e7b75421df477738a85f77ca"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "f879c4864f91c2a0f10db23d6bb2e844e4c4ca1c15dd41f69af5f95f42be8118"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "79c50d2badce5e04e30d53661c10365f8a1e3ca1da54bc169bd99c5659c3e630"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.15", shell_output("#{bin}/pmock-cli --version")
  end
end
