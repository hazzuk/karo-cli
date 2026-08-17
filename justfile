# SPDX-FileCopyrightText: © 2026 hazzuk
#
# SPDX-License-Identifier: AGPL-3.0-only

# ---

# Print help
help:
    @{{ just_executable() }} --list --unsorted --list-prefix "  - " --justfile "{{ justfile() }}"

# Compile current platform
build:
    go build -o dist/karo-cli
