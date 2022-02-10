package module

import (
	"fmt"
	"path/filepath"

	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	"github.com/digisan/user-mgr/udb"
)

const (
	root = "../data/user-space"
)

var (
	mMemPriv = map[string]Membership{
		"0": {Space: 30, NCasePerMon: 1},
		"1": {Space: 300, NCasePerMon: 4},
		"2": {Space: 3000, NCasePerMon: 16},
		"3": {Space: 30000, NCasePerMon: 64},
	}
)

func AllocDisk(uname string) {
	dir := filepath.Join(root, uname)
	gio.MustCreateDir(dir)
}

// unit: Megabyte
func CheckSpace(uname string) (all, used, available int, err error) {
	dir := filepath.Join(root, uname)
	if !fd.DirExists(dir) {
		return 0, 0, 0, fmt.Errorf("no space available for user [%s]", uname)
	}
	sz, err := fd.DirSize(dir, "m")
	if err != nil {
		return 0, 0, 0, err
	}
	user, ok, err := udb.UserDB.LoadActiveUser(uname)
	if err != nil || !ok {
		return 0, 0, 0, err
	}
	lvl := user.MemLevel
	all = mMemPriv[lvl].Space
	return all, int(sz), all - int(sz), nil
}
