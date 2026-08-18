// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
)

//go:embed templates/*
var templates embed.FS

var (
	stackGroup      string // username_scope
	stackName       string // jellyfin
	stackLicense    string // AGPL-3.0-only
	stackGroupUser  string // username
	stackGroupScope string // scope
)

func main() {

	getInput()
	validateInput()
	assertCustomRepo()
	createFiles()

}

func assertCustomRepo() {

	// get working directory path
	path, err := os.Getwd()
	check(err)

	dirName := filepath.Base(path)

	// check working directory name
	switch dirName {
	case "karo-custom":
		return
	case stackGroupUser:
		return
	case "karo-cli":
		return
	}

	fmt.Printf(
		"warn: non-standard name for current directory (%s), expected 'karo-custom' or '%s'\n",
		dirName, stackGroupUser,
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

func getInput() {

	usageMsg :=
		"Usage:  karo <command>\n\n" +
			"Commands:\n" +
			"  compose add    Create a custom stack for karo-compose\n"

	// create subcommand
	composeAddCmd := flag.NewFlagSet("compose add", flag.ExitOnError)

	// set flags
	composeAddCmd.StringVar(&stackGroup, "group", "", "Name of stack group (e.g. '<username>_<scope>')")
	composeAddCmd.StringVar(&stackName, "stack", "", "Name of stack (e.g. 'jellyfin')")
	composeAddCmd.StringVar(&stackLicense, "license", "NOASSERTION", "Optional, SPDX license identifier (e.g. 'AGPL-3.0-only')")

	// check args provided
	if len(os.Args) < 3 {
		fmt.Println(usageMsg)
		os.Exit(1)
	}

	// parse flags for subcommands
	switch os.Args[1] {
	case "compose":
		switch os.Args[2] {
		case "add":
			composeAddCmd.Parse(os.Args[3:])
			return
		}
	}

	fmt.Println(usageMsg)
	os.Exit(1)

}

func validateInput() {

	const (
		groupFlagError    = "error: \033[33m" + "-group" + "\033[0m"
		stackFlagError    = "error: \033[33m" + "-stack" + "\033[0m"
		requiredError     = "is required, see -help"
		alphanumericError = "must only use lowercase alphanumeric characters"
	)

	// stack group/name required
	if stackGroup == "" {
		fmt.Printf("%s %s\n", groupFlagError, requiredError)
		os.Exit(1)
	}
	if stackName == "" {
		fmt.Printf("%s %s\n", stackFlagError, requiredError)
		os.Exit(1)
	}

	// split stack group
	stackGroupParts := strings.Split(stackGroup, "_")

	// validate underscores
	if len(stackGroupParts) != 2 {
		fmt.Printf(
			"%s expected one underscore\n",
			groupFlagError,
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
			groupFlagError, stackGroupUser, alphanumericError,
		)
		os.Exit(1)
	}
	if !alphanumeric.MatchString(stackGroupScope) {
		fmt.Printf(
			"%s scope (%s) %s\n",
			groupFlagError, stackGroupScope, alphanumericError,
		)
		os.Exit(1)
	}
	if !alphanumeric.MatchString(stackName) {
		fmt.Printf("%s %s\n", stackFlagError, alphanumericError)
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
				groupFlagError, scope,
			)
			os.Exit(1)
		}
	}

}

func createFiles() {

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
		check(err)

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
		check(err)
	}

	// files
	files := [3]struct {
		path     string
		template string
	}{
		{
			path:     filepath.Join(directories[0], "main.yml"),
			template: "templates/main.yml.tmpl",
		},
		{
			path:     filepath.Join(directories[0], stackName+".yml"),
			template: "templates/stack.yml.tmpl",
		},
		{
			path:     filepath.Join(directories[1], "compose.yml.j2"),
			template: "templates/compose.yml.j2.tmpl",
		},
	}

	for _, file := range files {
		// create file
		f, err := os.OpenFile(file.path, os.O_CREATE|os.O_WRONLY, filePerm)
		check(err)

		// check file size
		info, err := f.Stat()
		check(err)

		if info.Size() == 0 {
			// template empty file
			tmpl, err := template.ParseFS(templates, file.template)
			check(err)

			data := struct {
				StackGroup     string
				StackGroupUser string
				StackName      string
				Year           int
				StackLicense   string
			}{
				StackGroup:     stackGroup,
				StackGroupUser: stackGroupUser,
				StackName:      stackName,
				Year:           time.Now().Year(),
				StackLicense:   stackLicense,
			}

			err = tmpl.Execute(f, data)
			check(err)

			f.Close()
			fmt.Println("info: created", file.path)

		} else if file.path == files[0].path {
			f.Close()
			// edit existing main.yml file
			editStackGroup(file.path, filePerm)
		}
	}

}

func editStackGroup(path string, filePerm os.FileMode) {

	// read main.yml file
	raw, err := os.ReadFile(path)
	check(err)

	content := string(raw)

	stackGroupDict := stackGroup + "_stacks:"

	// check for valid stack group dictionary
	if !strings.Contains(content, stackGroupDict) {
		fmt.Printf(
			"error: unable to find valid stack group dictionary (%s) in %s\n",
			stackGroupDict, path,
		)
		os.Exit(1)
	}

	// check for existing stack entry
	if strings.Contains(content, "- "+stackName+"\n") {
		return
	}

	// replace content with new stack entry
	content = strings.Replace(
		content,
		stackGroupDict,
		stackGroupDict+"\n  # - "+stackName,
		1,
	)

	// write modified content to file
	err = os.WriteFile(path, []byte(content), filePerm)
	check(err)

	fmt.Println("info: edited", path)

}

func check(err error) {

	if err != nil {
		fmt.Println("unexpected error: ", err)
		os.Exit(1)
	}

}
