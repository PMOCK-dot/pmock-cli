class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.13"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "d087e278275cd615d633affaad67d1520d45e631a4fc197dcf9eb98395c35d47"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "55d74282442658656ba86bc56c545f990e637aa2c46b9a6464e93fbf3124f9fe"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "04c13b3edbb4d4eea261cfcf6eaf378af5e3ca78ea693af88cd06c32dc365921"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "53ba9548ac41629dd07c874a83e8cb7b020df7912292421d09cf5ce561badc46"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.13", shell_output("#{bin}/pmock-cli --version")
  end
end
