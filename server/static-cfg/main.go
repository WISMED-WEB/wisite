package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/digisan/gotk/io"
	"github.com/digisan/gotk/project"
	lk "github.com/digisan/logkit"
)

const (
	config = "../api/system/config.go"
)

func writeln(ln string) {
	io.MustAppendFile(config, []byte(ln), true)
}

func recordVer() bool {
	ver, ok := project.GitVer("v0.0.0")
	ln := fmt.Sprintf(`    version = "%s"`, ver)
	writeln(ln)
	return ok
}

func recordTag() error {
	tag, err := project.GitTag()
	ln := fmt.Sprintf(`    tag     = "%s"`, tag)
	writeln(ln)
	return err
}

func main() {

	os.Remove(config)

	writeln(`package system`)
	writeln(``)
	writeln(`const (`)
	defer writeln(`)`)

	lk.WarnOnErrWhen(!recordVer(), "%v", errors.New("Couldn't get version"))
	lk.WarnOnErrWhen(recordTag() != nil, "%v", errors.New("Couldn't get tag"))
}
