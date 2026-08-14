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

	assertCustomDir()
	getInput()
	validateInput()
	createFiles()

}

func assertCustomDir() {

	// read current working directory
	path, err := os.Getwd()
	check(err)
	files, err := os.ReadDir(path)
	check(err)

	// find custom directory
	for _, file := range files {
		if file.Name() == "custom" {
			return
		}
	}

	log.Fatal("error: unable to find custom directory")

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
	customComposePath := filepath.Join("custom", stackGroupUser, "karo-compose")

	directories := [2]string{
		filepath.Join(customComposePath, "defaults", "main", stackGroup),
		filepath.Join(customComposePath, "templates", stackGroup, stackName),
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

			err = tmpl.Execute(f, struct {
				StackGroup string
				StackName  string
			}{
				StackGroup: stackGroup,
				StackName:  stackName,
			})
			check(err)
		}

		f.Close()
		fmt.Println("Created:", file.path)
	}

}

func check(err error) {

	if err != nil {
		log.Fatal(err)
	}

}
