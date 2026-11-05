class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.14"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "65a2e01900421a8796aa20fb46d6288a367e9d9205e3d695d4ab0355cb89b49a"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "dcdc233fd3694d73253e6ba4f150261e9ce337355e7a0e502882cc0786676e81"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "909041b447e17659dce99b86fec064724b133dc37c9f07046c0ef8bfbdd07c9e"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "754f08d44b089e296c217bc012430012da9724fb65ddd6c9efe79ff359d26f1e"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.14", shell_output("#{bin}/pmock-cli --version")
  end
end
