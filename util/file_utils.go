/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/
//nolint:golint,stylecheck
package util

import (
	"os"
	"path/filepath"
)

func init() { //nolint:gochecknoinits
	AcceptedExtensions[".spec"] = true
	AcceptedExtensions[".md"] = true
	AcceptedExtensions[".cpt"] = true
}

var AcceptedExtensions = make(map[string]bool) //nolint:golint,gochecknoglobals

func IsConceptFile(file string) bool { //nolint:golint
	return filepath.Ext(file) == ".cpt"
}

func IsValidSpecExtension(path string) bool {
	return AcceptedExtensions[filepath.Ext(path)]
}

func findFilesInDir(dirPath string, isValidFile func(path string) bool) []string {
	var files []string

	filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error { //nolint:errcheck,gosec
		if err == nil && !f.IsDir() && isValidFile(path) {
			files = append(files, path)
		}
		return err
	})

	return files
}

func findFilesIn(dirRoot string, isValidFile func(path string) bool) []string {
	absRoot := getAbsPath(dirRoot)
	files := findFilesInDir(absRoot, isValidFile)

	return files
}

func dirExists(dirPath string) bool {
	stat, err := os.Stat(dirPath)
	if err == nil && stat.IsDir() {
		return true
	}

	return false
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return !os.IsNotExist(err)
}

func GetFiles(path string) []string { //nolint:golint
	var specFiles []string
	if dirExists(path) {
		specFiles = append(specFiles, findFilesIn(path, IsValidSpecExtension)...)
	} else if FileExists(path) && IsValidSpecExtension(path) {
		specFiles = append(specFiles, getAbsPath(path))
	}

	return specFiles
}

func getAbsPath(path string) string {
	f, err := filepath.Abs(path)
	Fatal("Cannot get absolute path", err)

	return f
}
