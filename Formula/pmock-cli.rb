class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.18"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "f4836f4b7372c763dc739e6bd5a83296e92ba833bf01bfcfd5b715b3883947a2"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "cb656f958d763b971a29b9dec63d838f1e7719057eb0c725f1c115b5ee06cf91"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "ee56db9181005c924aed33c227e9259c21a955995879c9c5a1aceda66882b44c"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "97cfda15f9c98b3ec454cc33878d51e0a5d9a71e3b87995b40d356e4bc09957d"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.18", shell_output("#{bin}/pmock-cli --version")
  end
end
