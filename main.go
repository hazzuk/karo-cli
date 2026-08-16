// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
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
	}

	fmt.Printf(
		"warn: non-standard name for current directory (%s), expected 'karo-custom' or '%s'\n",
		dirName, stackGroupUser,
	)

	// read working directory files
	files, err := os.ReadDir(path)
	check(err)

	// find karo-compose directory
	for _, file := range files {
		if file.Name() == "karo-compose" {
			return
		}
	}

	log.Fatal(
		"error: running from non-standard karo-custom repo, ",
		"create ./karo-compose directory to override this",
	)

}

func getInput() {

	// parse cli flags
	flag.StringVar(&stackGroup, "group", "", "Name of stack group (e.g. '<username>_<scope>')")
	flag.StringVar(&stackName, "stack", "", "Name of stack (e.g. 'jellyfin')")
	flag.StringVar(&stackLicense, "license", "NOASSERTION", "Optional, SPDX license identifier (e.g. 'AGPL-3.0-only')")
	flag.Parse()

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
		log.Fatalf("%s %s", groupFlagError, requiredError)
	}
	if stackName == "" {
		log.Fatalf("%s %s", stackFlagError, requiredError)
	}

	// split stack group
	stackGroupParts := strings.Split(stackGroup, "_")

	// validate underscores
	if len(stackGroupParts) != 2 {
		log.Fatalf(
			"%s expected one underscore",
			groupFlagError,
		)
	}

	stackGroupUser = stackGroupParts[0]
	stackGroupScope = stackGroupParts[1]

	// validate alphanumeric & lowercase
	alphanumeric := regexp.MustCompile(`^[a-z0-9]+$`)

	if !alphanumeric.MatchString(stackGroupUser) {
		log.Fatalf(
			"%s username (%s) %s",
			groupFlagError, stackGroupUser, alphanumericError,
		)
	}
	if !alphanumeric.MatchString(stackGroupScope) {
		log.Fatalf(
			"%s scope (%s) %s",
			groupFlagError, stackGroupScope, alphanumericError,
		)
	}
	if !alphanumeric.MatchString(stackName) {
		log.Fatalf("%s %s", stackFlagError, alphanumericError)
	}

	// validate group scope
	reservedGroupScopes := [10]string{
		"compose", "docker", "git", "nftables", "ssh",
		"system", "stacks", "stack", "custom", "karo",
	}

	for _, scope := range reservedGroupScopes {
		if strings.Contains(stackGroupScope, scope) {
			log.Fatalf(
				"%s scope must not use word '%s'",
				groupFlagError, scope,
			)
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
	groupDirs, err := os.ReadDir(filepath.Join("karo-compose", "templates"))
	check(err)

	for _, dir := range groupDirs {
		if dir.IsDir() {
			// split existing stack group username
			dirParts := strings.Split(dir.Name(), "_")

			// compare usernames
			if stackGroupUser != dirParts[0] {
				log.Fatalf(
					"error: found mismatched stack groups (%s/%s)",
					stackGroupUser, dirParts[0],
				)
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
		log.Fatalf(
			"error: unable to find valid stack group dictionary (%s) in %s",
			stackGroupDict, path,
		)
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
		log.Fatal("unexpected error: ", err)
	}

}
