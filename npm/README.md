# npm distribution

The npm package is a thin wrapper: `pmock-cli` resolves a platform package
(`pmock-cli-<platform>-<arch>`) and executes its binary.

## Build

```bash
VERSION=0.1.0 ./scripts/build-npm.sh
```

This calls `scripts/package.sh` if `dist/` is missing, then packs:

- `pmock-cli-darwin-arm64`
- `pmock-cli-darwin-x64`
- `pmock-cli-linux-x64`
- `pmock-cli-linux-arm64`
- `pmock-cli-win32-x64`
- `pmock-cli` (wrapper)

Tarballs land in `npm/build/`.

## Publish

```bash
for f in npm/build/*-*.tgz; do npm publish "$f" --access public; done
npm publish npm/build/pmock-cli-0.1.0.tgz --access public
```

All names are unscoped so they can be published from the `pmock` npm account.
That account does not own the existing packages yet, so the current owner adds
it once per package (the win32 package was never published, so it is excluded):

```bash
for p in pmock-cli pmock-cli-darwin-arm64 pmock-cli-darwin-x64 \
         pmock-cli-linux-x64 pmock-cli-linux-arm64; do
  npm owner add pmock "$p"
done
```

`release.yml` authenticates as `pmock` through the `NPM_TOKEN` repository
secret, so the token has to come from that account.

## Test locally

```bash
npm install --no-save --omit=optional --ignore-scripts \
  npm/build/pmock-cli-darwin-arm64-0.1.0.tgz \
  npm/build/pmock-cli-0.1.0.tgz
./node_modules/.bin/pmock-cli version
```
