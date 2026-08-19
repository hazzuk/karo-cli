// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package generate

import (
	"os"
	"path/filepath"

	"github.com/hazzuk/karo-cli/internal/config"
)

const (
	dirPerm  os.FileMode = 0775
	filePerm os.FileMode = 0664
)

func CreateComposeStack(cfg *config.Config) {

	directories := [2]string{
		filepath.Join("karo-compose", "defaults", "main", cfg.StackGroup),
		filepath.Join("karo-compose", "templates", cfg.StackGroup, cfg.StackName),
	}

	createDirs(cfg, directories)
	createFiles(cfg, directories)

}
