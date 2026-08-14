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

var (
	stack_group       string // username_scope
	stack_name        string // jellyfin
	stack_group_user  string // username
	stack_group_scope string // scope
)

//go:embed templates/*
var templates embed.FS

func main() {

	assertCustomDir()

	getInput()

	validateInput()


	createFiles()

}

func assertCustomDir() {

	path, err := os.Getwd()
	check(err)
	files, err := os.ReadDir(path)
	check(err)

	for _, file := range files {
		if file.Name() == "custom" {
			return
		}
	}

	log.Fatal("error: unable to find custom directory")

}

func getInput() {

	// cli flags

	flag.StringVar(&stack_group, "group", "", "Name of stack group (e.g. '<username>_<scope>')")
	flag.StringVar(&stack_name, "stack", "", "Name of stack (e.g. 'jellyfin')")
	flag.Parse()

}

func validateInput() {

	const (
		group_flag_error = "error: \033[33m" + "-group" + "\033[0m "
		stack_flag_error = "error: \033[33m" + "-stack" + "\033[0m "
	)

	// group underscore

	stack_group_parts := strings.Split(stack_group, "_")

	if len(stack_group_parts) != 2 {
		log.Fatal(group_flag_error, "expected one underscore")
	}

	stack_group_user = stack_group_parts[0]
	stack_group_scope = stack_group_parts[1]

	// alphanumeric & lowercase

	alphanumeric := regexp.MustCompile(`^[a-z0-9]+$`)

	if !alphanumeric.MatchString(stack_group_user) {
		log.Fatal(group_flag_error, "username (", stack_group_user, ") must only contain lowercase alphanumeric characters")
	}
	if !alphanumeric.MatchString(stack_group_scope) {
		log.Fatal(group_flag_error, "scope (", stack_group_scope, ") must only contain lowercase alphanumeric characters")
	}
	if !alphanumeric.MatchString(stack_name) {
		log.Fatal(stack_flag_error, "must only contain lowercase alphanumeric characters")
	}

	// group scope

	const group_scope_error = group_flag_error + "scope must not use word "

	reserved_group_scopes := [10]string{
		"compose", "docker", "git", "nftables", "ssh",
		"system", "stacks", "stack", "custom", "karo",
	}

	for _, scope := range reserved_group_scopes {
		if strings.Contains(stack_group_scope, scope) {
			log.Fatal(group_scope_error, scope)
		}
	}

}

func createFiles() {

	const (
		dir_perm  = 0775
		file_perm = 0664
	)

	// directories

	custom_compose_path := filepath.Join("custom", stack_group_user, "karo-compose")

	directories := [2]string{
		filepath.Join(custom_compose_path, "defaults", "main", stack_group),
		filepath.Join(custom_compose_path, "templates", stack_group, stack_name),
	}

	for _, dir := range directories {
		err := os.MkdirAll(dir, dir_perm)
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
			path:     filepath.Join(directories[0], stack_name+".yml"),
			template: "templates/stack.yml.tmpl",
		},
		{
			path:     filepath.Join(directories[1], "compose.yml.j2"),
			template: "templates/compose.yml.j2.tmpl",
		},
	}

	for _, file := range files {
		// create file
		f, err := os.OpenFile(file.path, os.O_CREATE|os.O_WRONLY, file_perm)
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
				StackGroup: stack_group,
				StackName:  stack_name,
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
