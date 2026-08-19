// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hazzuk/karo-cli/internal/config"
)

func ValidateInput(cfg *config.Config) {

	const (
		groupArgError     = "error: \033[33m" + "group" + "\033[0m"
		stackArgError     = "error: \033[33m" + "stack" + "\033[0m name"
		requiredError     = "is required, see -help"
		alphanumericError = "must only use lowercase alphanumeric characters"
	)

	// stack group/name required
	if cfg.StackGroup == "" {
		fmt.Printf("%s %s\n", groupArgError, requiredError)
		os.Exit(1)
	}
	if cfg.StackName == "" {
		fmt.Printf("%s %s\n", stackArgError, requiredError)
		os.Exit(1)
	}

	// split stack group
	stackGroupParts := strings.Split(cfg.StackGroup, "_")

	// validate underscores
	if len(stackGroupParts) != 2 {
		fmt.Printf(
			"%s name expected 1 underscore, found %s\n",
			groupArgError,
			fmt.Sprint(len(stackGroupParts)-1),
		)
		os.Exit(1)
	}

	cfg.StackGroupUser = stackGroupParts[0]
	cfg.StackGroupScope = stackGroupParts[1]

	// validate alphanumeric & lowercase
	alphanumeric := regexp.MustCompile(`^[a-z0-9]+$`)

	if !alphanumeric.MatchString(cfg.StackGroupUser) {
		fmt.Printf(
			"%s username (%s) %s\n",
			groupArgError, cfg.StackGroupUser, alphanumericError,
		)
		os.Exit(1)
	}
	if !alphanumeric.MatchString(cfg.StackGroupScope) {
		fmt.Printf(
			"%s scope (%s) %s\n",
			groupArgError, cfg.StackGroupScope, alphanumericError,
		)
		os.Exit(1)
	}
	if !alphanumeric.MatchString(cfg.StackName) {
		fmt.Printf("%s %s\n", stackArgError, alphanumericError)
		os.Exit(1)
	}

	// validate group scope
	reservedGroupScopes := [10]string{
		"compose", "docker", "git", "nftables", "ssh",
		"system", "stacks", "stack", "custom", "karo",
	}

	for _, scope := range reservedGroupScopes {
		if strings.Contains(cfg.StackGroupScope, scope) {
			fmt.Printf(
				"%s scope must not use word '%s'\n",
				groupArgError, scope,
			)
			os.Exit(1)
		}
	}

}
