#!/bin/sh

set -eu

project_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
test_dir=$(mktemp -d "${TMPDIR:-/tmp}/d6roll-install-test.XXXXXXXX")
trap 'rm -rf "$test_dir"' 0
trap 'exit 1' HUP INT TERM

mkdir -p "$test_dir/fixtures" "$test_dir/mock-bin" "$test_dir/stage"

for target in linux_amd64 darwin_amd64 darwin_arm64; do
	printf '%s\n' "$target" > "$test_dir/stage/d6roll"
	tar -czf "$test_dir/fixtures/d6roll_0.1.0_${target}.tar.gz" -C "$test_dir/stage" d6roll
done

(
	cd "$test_dir/fixtures"
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum d6roll_*.tar.gz > checksums.txt
	else
		shasum -a 256 d6roll_*.tar.gz > checksums.txt
	fi
)

cat > "$test_dir/mock-bin/curl" <<'EOF'
#!/bin/sh
set -eu
output=
url=
while [ "$#" -gt 0 ]; do
	case "$1" in
		-o | --output) output=$2; shift 2 ;;
		--write-out) shift 2 ;;
		--fail | --location | --silent | --show-error) shift ;;
		*) url=$1; shift ;;
	esac
done
case "$url" in
	"https://raw.githubusercontent.com/kpb/d6roll/main/install.sh")
		if [ -n "$output" ]; then
			cp "$TEST_PROJECT_DIR/install.sh" "$output"
		else
			cat "$TEST_PROJECT_DIR/install.sh"
		fi
		;;
	"https://github.com/kpb/d6roll/releases/latest")
		printf '%s' 'https://github.com/kpb/d6roll/releases/tag/v0.1.0'
		;;
	"https://github.com/kpb/d6roll/releases/download/v0.1.0/"*)
		asset=${url##*/}
		[ "${TEST_FAIL_ASSET:-}" != "$asset" ] || exit 22
		cp "$TEST_FIXTURES/$asset" "$output"
		;;
	*) exit 22 ;;
esac
EOF

cat > "$test_dir/mock-bin/uname" <<'EOF'
#!/bin/sh
case "$1" in
	-s) printf '%s\n' "$TEST_OS" ;;
	-m) printf '%s\n' "$TEST_ARCH" ;;
	*) exit 1 ;;
esac
EOF

chmod +x "$test_dir/mock-bin/curl" "$test_dir/mock-bin/uname"
PATH="$test_dir/mock-bin:$PATH"
TEST_FIXTURES="$test_dir/fixtures"
TEST_PROJECT_DIR=$project_dir
export PATH TEST_FIXTURES TEST_PROJECT_DIR

fail() {
	printf 'installer test: %s\n' "$*" >&2
	exit 1
}

expect_binary() {
	label=$1
	install_dir=$2
	want=$3
	[ -x "$install_dir/d6roll" ] || fail "$label: missing executable"
	[ "$(cat "$install_dir/d6roll")" = "$want" ] || fail "$label: wrong binary"
}

expect_failure() {
	label=$1
	install_dir=$2
	if sh "$project_dir/install.sh" > "$test_dir/stdout" 2> "$test_dir/stderr"; then
		fail "$label: installer succeeded unexpectedly"
	fi
	[ ! -e "$install_dir/d6roll" ] || fail "$label: installed binary on failure"
}

# Exercise both documented curl | sh invocations with an offline curl stub.
curl -fsSL https://raw.githubusercontent.com/kpb/d6roll/main/install.sh |
	TEST_OS=Linux TEST_ARCH=amd64 HOME="$test_dir/documented-home" sh > "$test_dir/stdout"
expect_binary documented-command "$test_dir/documented-home/.local/bin" linux_amd64

curl -fsSL https://raw.githubusercontent.com/kpb/d6roll/main/install.sh |
	TEST_OS=Darwin TEST_ARCH=arm64 VERSION=v0.1.0 INSTALL_DIR="$test_dir/documented-override" sh > "$test_dir/stdout"
expect_binary documented-override "$test_dir/documented-override" darwin_arm64

TEST_OS=Linux TEST_ARCH=x86_64 HOME="$test_dir/home" sh "$project_dir/install.sh" > "$test_dir/stdout"
expect_binary latest "$test_dir/home/.local/bin" linux_amd64

TEST_OS=Darwin TEST_ARCH=x86_64 VERSION=v0.1.0 INSTALL_DIR="$test_dir/mac intel/bin" sh "$project_dir/install.sh" > "$test_dir/stdout"
expect_binary mac-intel "$test_dir/mac intel/bin" darwin_amd64

TEST_OS=Darwin TEST_ARCH=arm64 VERSION=v0.1.0 INSTALL_DIR="$test_dir/mac arm/bin" sh "$project_dir/install.sh" > "$test_dir/stdout"
expect_binary mac-arm "$test_dir/mac arm/bin" darwin_arm64

TEST_OS=Linux TEST_ARCH=arm64 VERSION=v0.1.0 INSTALL_DIR="$test_dir/unsupported" expect_failure linux-arm64 "$test_dir/unsupported"
TEST_OS=FreeBSD TEST_ARCH=amd64 VERSION=v0.1.0 INSTALL_DIR="$test_dir/unsupported" expect_failure freebsd "$test_dir/unsupported"
TEST_OS=Linux TEST_ARCH=amd64 VERSION=../../wrong INSTALL_DIR="$test_dir/invalid" expect_failure invalid-version "$test_dir/invalid"

mkdir "$test_dir/bad-checksum" "$test_dir/missing-checksum"
cp "$test_dir/fixtures/d6roll_0.1.0_linux_amd64.tar.gz" "$test_dir/bad-checksum/"
cp "$test_dir/fixtures/d6roll_0.1.0_linux_amd64.tar.gz" "$test_dir/missing-checksum/"
printf '%064d  %s\n' 0 d6roll_0.1.0_linux_amd64.tar.gz > "$test_dir/bad-checksum/checksums.txt"
: > "$test_dir/missing-checksum/checksums.txt"

TEST_OS=Linux TEST_ARCH=amd64 VERSION=v0.1.0 TEST_FIXTURES="$test_dir/bad-checksum" INSTALL_DIR="$test_dir/bad" expect_failure checksum-mismatch "$test_dir/bad"
TEST_OS=Linux TEST_ARCH=amd64 VERSION=v0.1.0 TEST_FIXTURES="$test_dir/missing-checksum" INSTALL_DIR="$test_dir/missing" expect_failure missing-checksum "$test_dir/missing"
TEST_OS=Linux TEST_ARCH=amd64 VERSION=v0.1.0 TEST_FAIL_ASSET=d6roll_0.1.0_linux_amd64.tar.gz INSTALL_DIR="$test_dir/download" expect_failure download-failure "$test_dir/download"

printf '%s\n' 'Installer tests passed'
