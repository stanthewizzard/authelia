package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	cmdscripts "github.com/authelia/authelia/v4/cmd/authelia-scripts/cmd"
	"github.com/authelia/authelia/v4/internal/commands"
)

func newDocsCLICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cmdUseDocsCLI,
		Short: "Generate CLI docs",
		RunE:  docsCLIRunE,

		DisableAutoGenTag: true,
	}

	return cmd
}

func docsCLIRunE(cmd *cobra.Command, args []string) (err error) {
	var root, pathDocsCLIReference string

	if root, err = cmd.Flags().GetString(cmdFlagRoot); err != nil {
		return err
	}

	if pathDocsCLIReference, err = cmd.Flags().GetString(cmdFlagDocsCLIReference); err != nil {
		return err
	}

	fullPathDocsCLIReference := filepath.Join(root, pathDocsCLIReference)

	if err = os.MkdirAll(fullPathDocsCLIReference, 0775); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	if err = genCLIDoc(commands.NewRootCmd(), filepath.Join(fullPathDocsCLIReference, "authelia")); err != nil {
		return err
	}

	if err = genCLIDocWriteIndex(fullPathDocsCLIReference, "authelia"); err != nil {
		return err
	}

	if err = genCLIDoc(cmdscripts.NewRootCmd(), filepath.Join(fullPathDocsCLIReference, "authelia-scripts")); err != nil {
		return err
	}

	if err = genCLIDocWriteIndex(fullPathDocsCLIReference, "authelia-scripts"); err != nil {
		return err
	}

	if err = genCLIDoc(newRootCmd(), filepath.Join(fullPathDocsCLIReference, cmdUseRoot)); err != nil {
		return err
	}

	if err = genCLIDocWriteIndex(fullPathDocsCLIReference, cmdUseRoot); err != nil {
		return err
	}

	return nil
}

func genCLIDoc(cmd *cobra.Command, path string) (err error) {
	if _, err = os.Stat(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil || !os.IsNotExist(err) {
		if err = os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove docs: %w", err)
		}
	}

	if err = os.Mkdir(path, 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	if err = doc.GenMarkdownTreeCustom(cmd, path, prepend, linker); err != nil {
		return err
	}

	return nil
}

func genCLIDocWriteIndex(path, name string) (err error) {
	now := time.Now()

	f, err := os.Create(filepath.Join(path, name, "_index.md"))
	if err != nil {
		return err
	}

	weight := 900

	if name == "authelia" {
		weight = 320
	}

	_, err = fmt.Fprintf(f, indexDocs, name, now.Format(dateFmtYAML), "cli-"+name, weight)

	return err
}

func prepend(input string) string {
	now := time.Now()

	pathz := strings.Split(strings.Replace(input, ".md", "", 1), "\\")
	parts := strings.Split(pathz[len(pathz)-1], "_")

	cmd := parts[0]

	args := strings.Join(parts, " ")

	weight := 330
	if len(parts) == 1 {
		weight = 320
	}

	return fmt.Sprintf(prefixDocs, args, fmt.Sprintf("Reference for the %s command.", args), "", now.Format(dateFmtYAML), "cli-"+cmd, weight)
}

func linker(input string) string {
	return input
}

const indexDocs = `---
title: "%s"
description: ""
lead: ""
date: %s
draft: false
images: []
menu:
  reference:
    parent: "cli"
    identifier: "%s"
weight: %d
toc: true
---
`

const prefixDocs = `---
title: "%s"
description: "%s"
lead: "%s"
date: %s
draft: false
images: []
menu:
  reference:
    parent: "%s"
weight: %d
toc: true
---

`
