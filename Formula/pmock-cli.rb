class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.11"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "7908ca28745e410e086ed287e479e1cc4e9405f223c189002a68e6272efcaffe"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "f16eb76d3c198b549016c9462f02afbc3a4b59330b7534245a0bc053e64ef93f"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "6f4db8c1f52a81c110c4351c0288f1e487f353a9385d8cc3891dd025471b6071"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "e4eeca788e0ff85215d7479f466256822a9f81fc035d07dd06b49ef63bc8092c"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.11", shell_output("#{bin}/pmock-cli --version")
  end
end
