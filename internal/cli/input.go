package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func GetUserInput() {

	// create flag sets
	rootCmd := flag.NewFlagSet("karo", flag.ExitOnError)
	composeAddCmd := flag.NewFlagSet("compose add", flag.ExitOnError)

	// set flags
	composeAddCmd.StringVar(
		&stackLicense,
		"license",
		"NOASSERTION",
		"SPDX license identifier",
	)

	// create usage messages
	rootCmd.Usage = func() {
		fmt.Printf("Usage:  karo COMMAND\n\n")

		fmt.Printf("Commands:\n")
		fmt.Printf("  compose add    Create a custom stack for karo-compose\n\n")
	}

	composeAddCmd.Usage = func() {
		fmt.Printf("Usage:  karo compose add [OPTIONS] STACK_ID\n\n")

		fmt.Printf("Create a custom stack for karo-compose.\n\n")

		fmt.Printf("Options:\n")
		composeAddCmd.PrintDefaults()
		fmt.Printf("\n")

		fmt.Printf("Arguments:\n")
		fmt.Printf("  STACK_ID    Stack identifier (format \"<username>_<scope>/<stack>\")\n\n")

		fmt.Printf("Examples:\n")
		fmt.Printf("  karo compose add hazzuk_core/traefik\n")
		fmt.Printf("  karo compose add -license AGPL-3.0-only mosslocker_photos/immich\n\n")
	}

	// check args were provided
	if len(os.Args) < 3 {
		// handle -help
		rootCmd.Parse(os.Args[1:])
		// handle no args
		rootCmd.Usage()
		os.Exit(1)
	}

	// parse args
	switch os.Args[1] {
	case "compose":

		switch os.Args[2] {
		case "add":
			// parse args after subcommand
			composeAddCmd.Parse(os.Args[3:])

			// store parsed args
			args := composeAddCmd.Args()

			// check args for subcommand
			if len(args) != 1 {
				// no args
				composeAddCmd.Usage()
				os.Exit(1)
			}

			// split stack_id
			parts := strings.Split(args[0], "/")

			// check stack_id format
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				fmt.Printf(
					"error: invalid stack identifier (%s), expected <group>/<stack>\n",
					args[0],
				)
				os.Exit(1)
			}

			// use args
			stackGroup = parts[0]
			stackName = parts[1]
			return
		}

	}

	// handle subcommand -help
	rootCmd.Parse(os.Args[2:])

	// handle no args after subcommand, or wrong args
	rootCmd.Usage()
	os.Exit(1)

}
