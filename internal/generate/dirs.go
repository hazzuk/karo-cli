// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hazzuk/karo-cli/internal/config"
	"github.com/hazzuk/karo-cli/internal/utils"
)

func CreateDirs(cfg *config.Config) {

	const (
		dirPerm  os.FileMode = 0775
		filePerm os.FileMode = 0664
	)

	// directories
	directories := [2]string{
		filepath.Join("karo-compose", "defaults", "main", cfg.StackGroup),
		filepath.Join("karo-compose", "templates", cfg.StackGroup, cfg.StackName),
	}

	// validate against existing stack group dirs
	templatesPath := filepath.Join("karo-compose", "templates")

	// check templates dir exists
	if stat, err := os.Stat(templatesPath); err == nil && stat.IsDir() {
		// read templates dir
		groupDirs, err := os.ReadDir(templatesPath)
		utils.Check(err)

		// iterate templates sub-dirs
		for _, dir := range groupDirs {
			if dir.IsDir() {
				// split username from existing stack group dir
				dirParts := strings.Split(dir.Name(), "_")

				// ensure stack group usernames match
				if cfg.StackGroupUser != dirParts[0] {
					fmt.Printf(
						"error: found mismatched stack groups (%s/%s)\n",
						cfg.StackGroupUser, dirParts[0],
					)
					os.Exit(1)
				}
			}
		}
	}

	// create directories
	for _, dir := range directories {
		err := os.MkdirAll(dir, dirPerm)
		utils.Check(err)
	}

}
