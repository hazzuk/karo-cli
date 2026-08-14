package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	stack_group       string // username_scope
	stack_name        string // jellyfin
	stack_group_user  string // username
	stack_group_scope string // scope
)

func main() {

	// check custom dir exists?
	// is git directory?

	getInput()

	validateInput()

	createFiles()

	// debug

	fmt.Println(stack_group)
	fmt.Println(stack_name)

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

	files := [3]string{
		filepath.Join(directories[0], "main.yml"),
		filepath.Join(directories[0], stack_name+".yml"),
		filepath.Join(directories[1], "compose.yml.j2"),
	}

	for _, filename := range files {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, file_perm)
		check(err)
		f.Close()
	}

}

func check(err error) {

	if err != nil {
		log.Fatal(err)
	}

}
