// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"fmt"
	"os"

	"github.com/hazzuk/karo-cli/internal/cli"
	"github.com/hazzuk/karo-cli/internal/generate"
	"github.com/hazzuk/karo-cli/internal/lint"
)

var (
	stackGroup      string // username_scope
	stackName       string // jellyfin
	stackLicense    string // AGPL-3.0-only
	stackGroupUser  string // username
	stackGroupScope string // scope
)

func main() {

	cli.GetUserInput()
	cli.ValidateInput()
	lint.AssertCustomRepo()
	generate.CreateDirs()
	generate.CreateFiles()

}

func check(err error) {

	if err != nil {
		fmt.Println("unexpected error: ", err)
		os.Exit(1)
	}

}
