package rel

import (
	r "github.com/digisan/user-mgr/relation"
)

func init() {
	// set user relation db
	r.InitDB("./data/db-user")
}
