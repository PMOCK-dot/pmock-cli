#!/usr/bin/env bash
#
# pmock-cli 一键发版。
#
#   ./scripts/release.sh                        自动推断版本并发布
#   ./scripts/release.sh 0.1.17                 指定版本（v 前缀可省略）
#   ./scripts/release.sh --dry-run              只打印计划，不动仓库
#   ./scripts/release.sh --skip-formula         发版但不回填 Homebrew formula
#   ./scripts/release.sh --formula-only 0.1.17  只回填 formula（发布失败后补跑）
#   ./scripts/release.sh -y 0.1.17              跳过确认
#
# 流程：校验工作区 → 打 annotated tag → 推 main 和 tag → 等 CI 出包 →
#       下载 checksums.txt → 回填 Formula 的 version/sha256 → 提交并推送。
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

SLUG="${REPO_SLUG:-PMOCK-dot/pmock-cli}"
BRANCH="${RELEASE_BRANCH:-main}"
FORMULA="Formula/pmock-cli.rb"
WAIT_SECONDS="${RELEASE_WAIT_SECONDS:-900}"

DRY_RUN=0
ASSUME_YES=0
SKIP_FORMULA=0
FORMULA_ONLY=""
VERSION_ARG=""
TMP=""

cleanup() {
	[ -n "$TMP" ] && rm -rf "$TMP"
	return 0
}
trap cleanup EXIT

die() { printf 'error: %s\n' "$*" >&2; exit 1; }
info() { printf '==> %s\n' "$*"; }
ok() { printf '  ok  %s\n' "$*"; }

usage() {
	cat <<'EOF'
用法: scripts/release.sh [版本] [选项]

  版本                形如 0.1.17 或 v0.1.17；省略时按提交自动推断
                      （BREAKING -> major，feat -> minor，其余 -> patch）

选项:
  --dry-run           只打印将要执行的步骤，不修改仓库
  --formula-only VER  不发布，只按已发布的 checksums 回填 Formula
  --skip-formula      发版后不回填 Formula
  -y, --yes           跳过确认
  -h, --help          显示本帮助

环境变量:
  REPO_SLUG           默认 PMOCK-dot/pmock-cli
  RELEASE_BRANCH      默认 main
  RELEASE_WAIT_SECONDS  等待 CI 产物的上限，默认 900
EOF
	exit 0
}

while [ $# -gt 0 ]; do
	case "$1" in
	-h | --help) usage ;;
	--dry-run) DRY_RUN=1; shift ;;
	-y | --yes) ASSUME_YES=1; shift ;;
	--skip-formula) SKIP_FORMULA=1; shift ;;
	--formula-only)
		[ $# -ge 2 ] || die "--formula-only 需要版本号"
		FORMULA_ONLY="$2"
		shift 2
		;;
	-*) die "未知参数：$1（--help 查看用法）" ;;
	*)
		[ -z "$VERSION_ARG" ] || die "多余参数：$1"
		VERSION_ARG="$1"
		shift
		;;
	esac
done

normalize_version() {
	local v="${1#v}"
	[[ "$v" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "版本号格式不对：$1（示例：0.1.17）"
	printf '%s\n' "$v"
}

latest_tag() { git describe --tags --abbrev=0 --match 'v*' 2>/dev/null || true; }

suggest_version() {
	local last base major minor patch range subjects
	last="$(latest_tag)"
	if [ -n "$last" ]; then
		base="${last#v}"
		range="$last..HEAD"
	else
		base="0.0.0"
		range="HEAD"
	fi
	IFS=. read -r major minor patch <<<"$base"
	subjects="$(git log --format='%s%n%b' $range)"
	if grep -qE '^[A-Za-z]+(\([^)]*\))?!:' <<<"$subjects" ||
		grep -qE 'BREAKING[ -]CHANGE' <<<"$subjects"; then
		major=$((major + 1))
		minor=0
		patch=0
	elif grep -qE '^feat(\([^)]*\))?!?:' <<<"$subjects"; then
		minor=$((minor + 1))
		patch=0
	else
		patch=$((patch + 1))
	fi
	printf '%s.%s.%s\n' "$major" "$minor" "$patch"
}

require_clean_repo() {
	git rev-parse --is-inside-work-tree >/dev/null 2>&1 || die "当前目录不是 git 仓库"
	local dirty branch behind
	dirty="$(git status --porcelain)"
	[ -z "$dirty" ] || die "工作区有未提交改动，先 commit 再发版：
$dirty"
	branch="$(git rev-parse --abbrev-ref HEAD)"
	[ "$branch" = "$BRANCH" ] || die "当前在 ${branch}，发版请在 $BRANCH 上"
	git fetch --quiet origin "$BRANCH" || die "git fetch origin $BRANCH 失败（检查网络或 SSH key）"
	behind="$(git rev-list --count "HEAD..origin/$BRANCH")"
	[ "$behind" = "0" ] || die "落后 origin/$BRANCH $behind 个提交，先 git pull --rebase"
}

confirm() {
	[ "$ASSUME_YES" = "1" ] && return 0
	[ -t 0 ] || die "非交互环境请加 --yes"
	local answer
	read -r -p "$1 [y/N] " answer
	case "$answer" in y | Y | yes | YES) return 0 ;; *) return 1 ;; esac
}

create_tag() {
	local version="$1" tag="v$1"
	git tag -a "$tag" -m "$tag"
	ok "已创建 tag $tag"
	git push origin "$BRANCH"
	ok "已推送 $BRANCH"
	git push origin "$tag"
	ok "已推送 tag ${tag}（已触发 release workflow）"
}

fetch_checksums() {
	local version="$1" waited=0
	local url="https://github.com/$SLUG/releases/download/v$version/checksums.txt"
	TMP="$(mktemp -d)"
	while :; do
		if curl -fsSL --max-time 20 "$url" -o "$TMP/checksums.txt" 2>/dev/null; then
			ok "已下载 checksums.txt"
			return 0
		fi
		if [ "$waited" -ge "$WAIT_SECONDS" ]; then
			die "等了 ${WAIT_SECONDS}s 仍没有 v$version 的 checksums.txt。
  Release 可能还在跑或已失败：https://github.com/$SLUG/actions
  稍后补跑：./scripts/release.sh --formula-only $version"
		fi
		info "等待 CI 发布产物…（已等 ${waited}s）"
		sleep 15
		waited=$((waited + 15))
	done
}

rewrite_formula() {
	local version="$1" sums="$2"
	python3 - "$FORMULA" "$version" "$sums" <<'PY'
import re
import sys

formula, version, sums_path = sys.argv[1], sys.argv[2], sys.argv[3]

sums = {}
with open(sums_path, encoding="utf-8") as fh:
    for line in fh:
        parts = line.split()
        if len(parts) == 2 and re.fullmatch(r"[0-9a-f]{64}", parts[0]):
            sums[parts[1]] = parts[0]

with open(formula, encoding="utf-8") as fh:
    lines = fh.readlines()

out, pending, used = [], None, set()
for line in lines:
    match = re.search(
        r"pmock-cli_#\{version\}_(darwin_amd64|darwin_arm64|linux_amd64|linux_arm64)\.tar\.gz",
        line,
    )
    if match:
        pending = "pmock-cli_%s_%s.tar.gz" % (version, match.group(1))
    elif pending and "sha256" in line:
        if pending not in sums:
            sys.exit("checksums.txt 里没有 %s" % pending)
        indent = line[: len(line) - len(line.lstrip())]
        line = '%ssha256 "%s"\n' % (indent, sums[pending])
        used.add(pending)
        pending = None
    elif re.match(r'^\s*version\s+"', line):
        line = re.sub(r'(^\s*version\s+")[^"]+(")', r"\g<1>%s\g<2>" % version, line)
    elif "assert_match" in line:
        line = re.sub(r'"\d+\.\d+\.\d+"', '"%s"' % version, line, count=1)
    out.append(line)

if len(used) != 4:
    sys.exit("只回填了 %d 个平台（期望 4），Formula 结构可能变了" % len(used))

with open(formula, "w", encoding="utf-8") as fh:
    fh.writelines(out)
PY
	ok "已回填 $FORMULA 到 v$version"
}

sync_formula() {
	local version="$1"
	command -v python3 >/dev/null || die "回填 formula 需要 python3"
	fetch_checksums "$version"
	rewrite_formula "$version" "$TMP/checksums.txt"
	if git diff --quiet -- "$FORMULA"; then
		ok "$FORMULA 已是最新，无需提交"
		return 0
	fi
	git add "$FORMULA"
	git commit -m "release: v$version formula and checksums"
	ok "已提交 formula 更新"
	git push origin "$BRANCH"
	ok "已推送 formula 更新"
}

show_links() {
	local version="$1"
	printf '\n'
	info "v$version 已发布"
	printf '  Release  https://github.com/%s/releases/tag/v%s\n' "$SLUG" "$version"
	printf '  npm      https://www.npmjs.com/package/pmock-cli/v/%s\n' "$version"
	printf '  升级     npm install -g pmock-cli@%s\n' "$version"
}

# ---------------------------------------------------------------- 入口

[ "$DRY_RUN" = "1" ] || require_clean_repo

if [ -n "$FORMULA_ONLY" ]; then
	VERSION="$(normalize_version "$FORMULA_ONLY")"
	info "只回填 $FORMULA 到 v$VERSION"
	if [ "$DRY_RUN" = "1" ]; then
		info "dry-run：下载 checksums.txt 并写入 ${FORMULA}（不会提交）"
		exit 0
	fi
	sync_formula "$VERSION"
	exit 0
fi

if [ -n "$VERSION_ARG" ]; then
	VERSION="$(normalize_version "$VERSION_ARG")"
else
	VERSION="$(suggest_version)"
fi
TAG="v$VERSION"

git rev-parse -q --verify "refs/tags/$TAG" >/dev/null 2>&1 &&
	die "tag $TAG 已存在。要重发请先删除：
  git tag -d $TAG && git push origin :refs/tags/$TAG"

LAST="$(latest_tag)"
if [ -n "$LAST" ]; then
	info "版本：$LAST → ${TAG}（$(git rev-list --count "$LAST..HEAD") 个提交）"
	git log --oneline -12 "$LAST..HEAD"
else
	info "版本：${TAG}（首个 tag）"
	git log --oneline -12
fi

if [ "$DRY_RUN" = "1" ]; then
	info "dry-run：将执行 git tag -a $TAG -m $TAG && git push origin $BRANCH && git push origin $TAG"
	[ "$SKIP_FORMULA" = "1" ] || info "dry-run：随后等待 checksums.txt 并回填 $FORMULA"
	exit 0
fi

confirm "推送 $TAG 并触发发布？" || {
	info "已取消"
	exit 0
}

create_tag "$VERSION"

if [ "$SKIP_FORMULA" = "1" ]; then
	info "已跳过 formula 回填；稍后可跑 ./scripts/release.sh --formula-only $VERSION"
else
	sync_formula "$VERSION"
fi

show_links "$VERSION"
