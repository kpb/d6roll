#!/bin/sh

set -eu

repository=kpb/d6roll
install_dir=${INSTALL_DIR:-${HOME:?HOME must be set}/.local/bin}

fail() {
	printf 'd6roll installer: %s\n' "$*" >&2
	exit 1
}

if [ -n "${VERSION:-}" ]; then
	version=$VERSION
else
	latest_url=$(curl --fail --location --silent --show-error --output /dev/null --write-out '%{url_effective}' "https://github.com/$repository/releases/latest") || fail 'finding the latest release'
	version=${latest_url##*/}
fi

# Only stable release tags belong in release asset URLs.
printf '%s\n' "$version" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$' || fail "invalid release version: $version"

case "$(uname -s)" in
	Linux) os=linux ;;
	Darwin) os=darwin ;;
	*) fail "unsupported operating system: $(uname -s)" ;;
esac

case "$(uname -m)" in
	x86_64 | amd64) arch=amd64 ;;
	aarch64 | arm64) arch=arm64 ;;
	*) fail "unsupported architecture: $(uname -m)" ;;
esac

if [ "$os" = linux ] && [ "$arch" = arm64 ]; then
	fail 'Linux arm64 release binaries are not available'
fi

archive="d6roll_${version#v}_${os}_${arch}.tar.gz"
release_url="https://github.com/$repository/releases/download/$version"
temporary_directory=$(mktemp -d "${TMPDIR:-/tmp}/d6roll.XXXXXXXX") || fail 'creating temporary directory'
trap 'rm -rf "$temporary_directory"' 0
trap 'exit 1' HUP INT TERM

curl --fail --location --silent --show-error --output "$temporary_directory/$archive" "$release_url/$archive" || fail "downloading $archive"
curl --fail --location --silent --show-error --output "$temporary_directory/checksums.txt" "$release_url/checksums.txt" || fail 'downloading checksums.txt'

expected_checksum=$(awk -v archive="$archive" '$2 == archive { print $1; exit }' "$temporary_directory/checksums.txt")
[ -n "$expected_checksum" ] || fail "missing checksum for $archive"

if command -v sha256sum >/dev/null 2>&1; then
	actual_checksum=$(sha256sum "$temporary_directory/$archive" | awk '{ print $1 }')
elif command -v shasum >/dev/null 2>&1; then
	actual_checksum=$(shasum -a 256 "$temporary_directory/$archive" | awk '{ print $1 }')
else
	fail 'sha256sum or shasum is required to verify the download'
fi
[ "$actual_checksum" = "$expected_checksum" ] || fail "checksum verification failed for $archive"

tar -xzf "$temporary_directory/$archive" -C "$temporary_directory" d6roll || fail "extracting $archive"
[ -f "$temporary_directory/d6roll" ] && [ ! -L "$temporary_directory/d6roll" ] || fail 'archive does not contain a regular d6roll binary'

mkdir -p "$install_dir"
install -m 0755 "$temporary_directory/d6roll" "$install_dir/d6roll"
printf 'Installed d6roll %s to %s/d6roll\n' "$version" "$install_dir"
