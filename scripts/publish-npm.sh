#!/usr/bin/env bash
# Publish the npm platform packages and wrapper idempotently.
#
#   VERSION=0.1.2 ./scripts/publish-npm.sh
#
# An account with two-factor auth on writes needs a one-time password for a
# direct publish. Pass it once and every package reuses it:
#
#   NPM_OTP=123456 VERSION=0.1.2 ./scripts/publish-npm.sh
#
# Existing versions are skipped, so re-running after a code expires or after a
# partial failure publishes only what is still missing. A platform package that
# npm rejects (for example Windows spam review) is reported but does not block
# the wrapper: npm treats optionalDependencies as optional, so macOS/Linux
# installs still work while the Windows package is being reviewed.
set -uo pipefail

cd "$(dirname "$0")/.."

VERSION="${VERSION:?set VERSION, e.g. VERSION=0.1.2}"
DIR="${DIR:-npm/build}"
REGISTRY="${REGISTRY:-https://registry.npmjs.org}"
failures=0

if [ ! -f "$DIR/pmock-cli-${VERSION}.tgz" ]; then
  echo "missing $DIR/pmock-cli-${VERSION}.tgz; run: VERSION=$VERSION make npm" >&2
  exit 1
fi

echo "npm user: $(npm whoami --registry "$REGISTRY" 2>/dev/null || echo unknown)"
echo "publishing pmock-cli@${VERSION} to ${REGISTRY}"

publish_one() {
  local name="$1"
  local file="$2"
  if npm view "${name}@${VERSION}" version --registry "$REGISTRY" >/dev/null 2>&1; then
    echo "skip ${name}@${VERSION} (already published)"
    return 0
  fi
  echo "→ ${name}@${VERSION}"
  if npm publish "$file" --registry "$REGISTRY" --access public ${NPM_OTP:+--otp="$NPM_OTP"}; then
    echo "published ${name}@${VERSION}"
    return 0
  fi
  echo "WARN: failed to publish ${name}@${VERSION}" >&2
  failures=$((failures + 1))
  return 1
}

available=0
for platform in darwin-arm64 darwin-x64 linux-x64 linux-arm64 win32-x64; do
  if publish_one "pmock-cli-${platform}" "$DIR/pmock-cli-${platform}-${VERSION}.tgz"; then
    available=$((available + 1))
  fi
done

# The wrapper pins every platform package at this version, so publishing it
# while none of them exist puts a CLI that cannot run behind `latest`. One
# rejected platform (Windows review is common) is tolerated: installs on the
# other platforms still work.
if [ "$available" -eq 0 ]; then
  echo "no pmock-cli platform package is available at ${VERSION}; refusing to publish the wrapper" >&2
  exit 1
fi

if publish_one "pmock-cli" "$DIR/pmock-cli-${VERSION}.tgz"; then
  echo
  echo "wrapper available: npm install -g pmock-cli@${VERSION}"
  if [ "$failures" -gt 0 ]; then
    echo "note: $failures platform package(s) still need attention; wrapper installs only matching platforms." >&2
  fi
  exit 0
fi

# The wrapper publish reported failure; npm may still be indexing it. Give the
# registry another minute before failing the job.
for _ in $(seq 1 20); do
  if npm view "pmock-cli@${VERSION}" version --registry "$REGISTRY" >/dev/null 2>&1; then
    echo "wrapper is visible after all"
    exit 0
  fi
  sleep 3
done
echo "wrapper was not published" >&2
exit 1
