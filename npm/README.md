# npm distribution

The npm package is a thin wrapper: `pmock-cli` resolves a platform package
(`pmock-cli-<platform>-<arch>`) and executes its binary.

Windows is spelled `windows` in the package name but stays `win32` in npm's
`os` field: npm's spam filter rejects the `pmock-cli-win32-x64` variant name
with `Package name triggered spam detection`.

## Build

```bash
VERSION=0.1.0 ./scripts/build-npm.sh
```

This calls `scripts/package.sh` if `dist/` is missing, then packs:

- `pmock-cli-darwin-arm64`
- `pmock-cli-darwin-x64`
- `pmock-cli-linux-x64`
- `pmock-cli-linux-arm64`
- `pmock-cli-windows-x64`
- `pmock-cli` (wrapper)

Tarballs land in `npm/build/`.

## Publish

```bash
VERSION=0.1.11 ./scripts/publish-npm.sh
```

The script publishes the five platform packages and the wrapper, skips versions
that already exist, and keeps going when npm rejects one platform package. It
refuses to publish the wrapper when no platform package is available, because
the wrapper pins all of them at its own version.

An account with two-factor auth on writes needs a code for these publishes; it
is reused for every package:

```bash
NPM_OTP=123456 VERSION=0.1.11 ./scripts/publish-npm.sh
```

All names are unscoped and owned by the `pmock` account, which is the identity
behind the `NPM_TOKEN` repository secret.

The token must be a granular access token with **All packages** scope and
**Bypass 2FA** enabled. The scope matters when a platform package is published
for the first time: a token limited to a list of existing packages can never
cover a name that does not exist yet, and the registry answers 403.

The `pmock` account uses `auth-and-writes` two-factor auth, so interactive
commands fail without a code:

```
npm error code EOTP
npm error This operation requires a one-time password from your authenticator.
```

`--otp=<code>` covers interactive use; CI cannot, which is why the repository
secret has to bypass 2FA. Note that npm is restricting bypass-2FA tokens for
account changes and direct publishing, so this may need revisiting.

## Test locally

```bash
npm install --no-save --omit=optional --ignore-scripts \
  npm/build/pmock-cli-darwin-arm64-0.1.0.tgz \
  npm/build/pmock-cli-0.1.0.tgz
./node_modules/.bin/pmock-cli version
```
