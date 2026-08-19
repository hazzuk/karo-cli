// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hazzuk/karo-cli/internal/utils"
)

func CreateDirs() {

	const (
		dirPerm  os.FileMode = 0775
		filePerm os.FileMode = 0664
	)

	// directories
	directories := [2]string{
		filepath.Join("karo-compose", "defaults", "main", stackGroup),
		filepath.Join("karo-compose", "templates", stackGroup, stackName),
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
				if stackGroupUser != dirParts[0] {
					fmt.Printf(
						"error: found mismatched stack groups (%s/%s)\n",
						stackGroupUser, dirParts[0],
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
