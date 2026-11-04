class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.12"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "1bfd1c3410f5f6dae77d5fbe22eb70e8fbdd366e294a5f32896955d79acca517"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "59df45a9e00a6b16269febbcd1d1b6e222fc72f239aab2dd5695d38219d7df6f"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "9deda962b4460f2df8815d5a93e51ff24c18a25b6278c53dd4593c3a375ee184"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "9f448e3fd663de105543ac3579cc2f467cbf84bfb6c4239c6395e622e2cfb2c1"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.12", shell_output("#{bin}/pmock-cli --version")
  end
end
