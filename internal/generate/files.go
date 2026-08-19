// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

package generate

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hazzuk/karo-cli/internal/utils"
)

//go:embed templates/*
var templates embed.FS

func CreateFiles() {

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
		utils.Check(err)

		// check file size
		info, err := f.Stat()
		utils.Check(err)

		if info.Size() == 0 {
			// template empty file
			tmpl, err := template.ParseFS(templates, file.template)
			utils.Check(err)

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
			utils.Check(err)

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
	utils.Check(err)

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
	utils.Check(err)

	fmt.Println("info: edited", path)

}
