class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.19"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "d9b10184d59bb16ac5a70c102d2d2c8c281cd6a8d30e4142325000f32185c8bc"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "76c69fc737c14a7c444ca3c19767145e9099a53076d0e86cbed5eda997fbb7c3"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "bc95485678b55ba03262aa47dd4e101c1c8303c4bbeb9f985ecd20d462b33384"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "c219f2d99c1ecea47de7e2b95e16ac8c74d9e4360d5f55e1e545e7387b330cb7"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.19", shell_output("#{bin}/pmock-cli --version")
  end
end
