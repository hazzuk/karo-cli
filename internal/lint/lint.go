// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package lint

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hazzuk/karo-cli/internal/config"
	"github.com/hazzuk/karo-cli/internal/utils"
)

func AssertCustomRepo(cfg *config.Config) {

	// get working directory path
	path, err := os.Getwd()
	utils.Check(err)

	dirName := filepath.Base(path)

	// check working directory name
	switch dirName {
	case "karo-custom":
		return
	case cfg.StackGroupUser:
		return
	case "karo-cli":
		return
	}

	fmt.Printf(
		"warn: non-standard name for current directory (%s), expected 'karo-custom' or '%s'\n",
		dirName, cfg.StackGroupUser,
	)

	// check karo-compose dir exists
	if stat, err := os.Stat("karo-compose"); err == nil && stat.IsDir() {
		return
	} else {
		fmt.Println(
			"error: running from non-standard custom repo, ",
			"create ./karo-compose directory to override this",
		)
		os.Exit(1)
	}

}
