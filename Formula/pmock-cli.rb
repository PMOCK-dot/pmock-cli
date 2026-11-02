class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/haidaisiganbei/pmock-cli"
  version "0.1.8"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "fad512ccd29abedc9914b4b8784410f362c49ca86267a3336df8cbcfda371bbc"
    else
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "91e674c0fa9c6ecb908e44f84f8053a585a71d79a25d51495c71816239e85369"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "748300f869b7edb81c04671a6a31b0ad49adffb48d483343aae18bf0e04103b9"
    else
      url "https://github.com/haidaisiganbei/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "745c1fc465e1a5c8aa9861333a3ea2b86a04d7487780683433e1cadaf9838adc"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.8", shell_output("#{bin}/pmock-cli --version")
  end
end
