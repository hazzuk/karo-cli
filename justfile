# SPDX-FileCopyrightText: © 2026 hazzuk
#
# SPDX-License-Identifier: AGPL-3.0-only

# justfile, for running project-specific commands.
# See https://just.systems/man/en for more information.

# Print help
help:
    @{{ just_executable() }} --list --unsorted --list-prefix "  - " --justfile "{{ justfile() }}"

# Compile current platform
[group('Compile')]
build:
    go build -o dist/karo-cli

# Compile all platforms
[group('Compile')]
build-all:
    #!/usr/bin/env bash
    set -euo pipefail
    # https://go.dev/doc/install/source#environment
    for os in darwin linux windows; do
        for arch in amd64 arm64; do
            ext=""
            if [ "$os" = "windows" ]; then ext=".exe"; fi
            echo "Building $os/$arch"
            GOOS="$os" GOARCH="$arch" go build -o "dist/karo-cli-${os}-${arch}${ext}"
        done
    done
