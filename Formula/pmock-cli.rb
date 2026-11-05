class PmockCli < Formula
  desc "PMOCK open API command line client"
  homepage "https://github.com/PMOCK-dot/pmock-cli"
  version "0.1.16"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "bfae164a4fd403a8e2dbcbdf0b1a8bd6b98e12c61d14ac1a59b3ec72720a7539"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "9a96e09d43f0a034f205be73c614a8108210594bff347b579bdd328e629aa239"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_arm64.tar.gz"
      sha256 "05449a97c45e4f608c7c294583bdd90a349857444af1b625f6cbb6029e48ced4"
    else
      url "https://github.com/PMOCK-dot/pmock-cli/releases/download/v#{version}/pmock-cli_#{version}_linux_amd64.tar.gz"
      sha256 "a4089c39e956fdd043f6bb35244809fd7cce42c7a10f08c429b8b529c8dcdc45"
    end
  end

  def install
    bin.install "pmock-cli"
  end

  test do
    assert_match "0.1.15", shell_output("#{bin}/pmock-cli --version")
  end
end
