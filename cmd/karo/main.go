// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"github.com/hazzuk/karo-cli/internal/cli"
	"github.com/hazzuk/karo-cli/internal/config"
	"github.com/hazzuk/karo-cli/internal/generate"
	"github.com/hazzuk/karo-cli/internal/lint"
)

func main() {

	cfg := &config.Config{}

	cli.GetUserInput(cfg)
	cli.ValidateInput(cfg)
	lint.AssertCustomRepo(cfg)
	generate.CreateComposeStack(cfg)

}
