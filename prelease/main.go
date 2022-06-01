package main

import nt "github.com/digisan/gotk/net-tool"

func LiteralLocIP2PubIP(oldport, newport int, filepaths ...string) error {
	for _, fpath := range filepaths {
		if err := nt.ChangeLocalUrlPort(fpath, oldport, newport, false, true); err != nil {
			return err
		}
		if err := nt.LocIP2PubIP(fpath, false, true); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	LiteralLocIP2PubIP(1323, 1323, "../server/main.go")
}
