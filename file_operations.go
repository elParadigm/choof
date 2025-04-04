package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m *model) checkFile() {
	info, err := os.Stat(m.path)
	if err != nil {
		m.err = err
		return
	}
	if !info.Mode().IsRegular() {
		m.isFile = false
		m.err = fmt.Errorf("the given path is not a regular file")
	}
}
func (m *model) getStat(arg string) {
	statCommand := exec.Command("stat", "--", arg)
	var stdout, stderr bytes.Buffer
	statCommand.Stdout = &stdout
	statCommand.Stderr = &stderr

	err := statCommand.Run()

	if err != nil {
		m.err = err
	}

	result := stdout.String()
	lines := strings.Split(result, "\n")

	resultMatrix := make([][]string, len(lines))

	for i, l := range lines {
		resultMatrix[i] = strings.Fields(l)
	}

	m.filename = filepath.Base(resultMatrix[0][1])

	//check if the file name contain spaces and include all the parts
	if len(resultMatrix[0]) > 2 {
		for i, val := range resultMatrix[0] {
			if i > 1 {
				m.filename = m.filename + " " + val
			}
		}
	}
	m.created = resultMatrix[8][1]
	m.permission = resultMatrix[3][1]
	m.owner = resultMatrix[3][4]
	m.Lastmodified = resultMatrix[6][1]
	m.size = convertSize(resultMatrix[1][1])

}

func (m *model) getPath(arg string) {
	pathCommand := exec.Command("realpath", arg)
	var stdout, stderr bytes.Buffer
	pathCommand.Stdout = &stdout
	pathCommand.Stderr = &stderr

	err := pathCommand.Run()

	if err != nil {
		m.err = err
	}

	m.path = strings.TrimSpace(stdout.String())

}
func (m *model) getFileType(arg string) {

	fileCommand := exec.Command("file", arg)
	var stdout, stderr bytes.Buffer
	fileCommand.Stdout = &stdout
	fileCommand.Stderr = &stderr

	err := fileCommand.Run()

	if err != nil {
		m.err = err
	}
	result := stdout.String()
	resultslice := strings.Split(result, ",")
	var s string
	for _, ch := range resultslice {
		s += "\n" + ch
	}
	s = s[(1 + strings.Index(s, ":")):]
	m.description = s

}
func (m *model) DeleteFile(arg string) {

	absPath := filepath.Clean(arg)

	error := os.Remove(absPath)
	if error != nil {
		m.err = error
	}

}

func (m *model) RenameFile(arg string, newname string) {
	absPath := filepath.Clean(arg)
	absPath = filepath.Join(filepath.Dir(absPath), newname)

	err := os.Rename(arg, absPath)
	if err != nil {
		m.err = err
	}
}

func (m *model) calculateMD5(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		m.err = err
	}
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func (m *model) calculateSHA256(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		m.err = err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
