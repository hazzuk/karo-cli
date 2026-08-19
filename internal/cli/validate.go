package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ValidateInput() {

	const (
		groupArgError     = "error: \033[33m" + "group" + "\033[0m"
		stackArgError     = "error: \033[33m" + "stack" + "\033[0m name"
		requiredError     = "is required, see -help"
		alphanumericError = "must only use lowercase alphanumeric characters"
	)

	// stack group/name required
	if stackGroup == "" {
		fmt.Printf("%s %s\n", groupArgError, requiredError)
		os.Exit(1)
	}
	if stackName == "" {
		fmt.Printf("%s %s\n", stackArgError, requiredError)
		os.Exit(1)
	}

	// split stack group
	stackGroupParts := strings.Split(stackGroup, "_")

	// validate underscores
	if len(stackGroupParts) != 2 {
		fmt.Printf(
			"%s name expected 1 underscore, found %s\n",
			groupArgError,
			fmt.Sprint(len(stackGroupParts)-1),
		)
		os.Exit(1)
	}

	stackGroupUser = stackGroupParts[0]
	stackGroupScope = stackGroupParts[1]

	// validate alphanumeric & lowercase
	alphanumeric := regexp.MustCompile(`^[a-z0-9]+$`)

	if !alphanumeric.MatchString(stackGroupUser) {
		fmt.Printf(
			"%s username (%s) %s\n",
			groupArgError, stackGroupUser, alphanumericError,
		)
		os.Exit(1)
	}
	if !alphanumeric.MatchString(stackGroupScope) {
		fmt.Printf(
			"%s scope (%s) %s\n",
			groupArgError, stackGroupScope, alphanumericError,
		)
		os.Exit(1)
	}
	if !alphanumeric.MatchString(stackName) {
		fmt.Printf("%s %s\n", stackArgError, alphanumericError)
		os.Exit(1)
	}

	// validate group scope
	reservedGroupScopes := [10]string{
		"compose", "docker", "git", "nftables", "ssh",
		"system", "stacks", "stack", "custom", "karo",
	}

	for _, scope := range reservedGroupScopes {
		if strings.Contains(stackGroupScope, scope) {
			fmt.Printf(
				"%s scope must not use word '%s'\n",
				groupArgError, scope,
			)
			os.Exit(1)
		}
	}

}
