package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/djherbis/times"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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
	fileInfo, err := os.Stat(arg)
	if err != nil {
		m.err = err
	}
	m.filename = fileInfo.Name()

	m.permission = fileInfo.Mode().String()

	m.Lastmodified = fileInfo.ModTime().String()
	m.size = fileInfo.Size()

	val := fileInfo.Sys().(*syscall.Stat_t)
	usr, err := user.LookupId(strconv.Itoa(int(val.Uid)))
	if err != nil {
		m.err = err
	}
	m.owner = usr.Username

	creationTime, err := times.Stat(arg)
	m.created = creationTime.BirthTime().Format("2006-01-02 15:04:05")

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
