package file

import fm "github.com/digisan/file-mgr"

func init() {
	// set user file space & file item db space
	fm.InitFileMgr("./data/")
}
