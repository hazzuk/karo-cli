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

	fmt.Println("warn: current directory ("+dirName+") has an unexpected name, \n",
		"expected 'karo-custom' or '"+stackGroupUser+"'")

	// read working directory files
	files, err := os.ReadDir(path)
	check(err)

	// find karo-compose directory
	for _, file := range files {
		if file.Name() == "karo-compose" {
			return
		}
	}

	log.Fatal("error: not inside a karo-custom repo, \n",
		"create a ./karo-compose directory to override this")

}

func getInput() {

	// parse cli flags
	flag.StringVar(&stackGroup, "group", "", "Name of stack group (e.g. '<username>_<scope>')")
	flag.StringVar(&stackName, "stack", "", "Name of stack (e.g. 'jellyfin')")
	flag.Parse()

}

func validateInput() {

	const (
		groupFlagError = "error: \033[33m" + "-group" + "\033[0m "
		stackFlagError = "error: \033[33m" + "-stack" + "\033[0m "
	)

	// split stack group
	stackGroupParts := strings.Split(stackGroup, "_")

	// validate underscores
	if len(stackGroupParts) != 2 {
		log.Fatal(groupFlagError, "expected one underscore")
	}

	stackGroupUser = stackGroupParts[0]
	stackGroupScope = stackGroupParts[1]

	// validate alphanumeric & lowercase
	alphanumeric := regexp.MustCompile(`^[a-z0-9]+$`)

	if !alphanumeric.MatchString(stackGroupUser) {
		log.Fatal(groupFlagError, "username (", stackGroupUser, ") must only contain lowercase alphanumeric characters")
	}
	if !alphanumeric.MatchString(stackGroupScope) {
		log.Fatal(groupFlagError, "scope (", stackGroupScope, ") must only contain lowercase alphanumeric characters")
	}
	if !alphanumeric.MatchString(stackName) {
		log.Fatal(stackFlagError, "must only contain lowercase alphanumeric characters")
	}

	// validate group scope
	const groupScopeError = groupFlagError + "scope must not use word "

	reservedGroupScopes := [10]string{
		"compose", "docker", "git", "nftables", "ssh",
		"system", "stacks", "stack", "custom", "karo",
	}

	for _, scope := range reservedGroupScopes {
		if strings.Contains(stackGroupScope, scope) {
			log.Fatal(groupScopeError, scope)
		}
	}

}

func createFiles() {

	const (
		dirPerm  = 0775
		filePerm = 0664
	)

	// directories
	directories := [2]string{
		filepath.Join("karo-compose", "defaults", "main", stackGroup),
		filepath.Join("karo-compose", "templates", stackGroup, stackName),
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
			}{
				StackGroup:     stackGroup,
				StackGroupUser: stackGroupUser,
				StackName:      stackName,
				Year:           time.Now().Year(),
			}

			err = tmpl.Execute(f, data)
			check(err)
		}

		f.Close()
		fmt.Println("created:", file.path)
	}

}

func check(err error) {

	if err != nil {
		log.Fatal(err)
	}

}
