package main

import (
	"context"
	"time"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	so "github.com/digisan/user-mgr/sign-out"
	u "github.com/digisan/user-mgr/user"
	"github.com/wismed-web/wisite-api/server/api/sign"
)

func monitorUser(ctx context.Context, offlineTimeout time.Duration) {
	cInactive := make(chan string, 4096)
	si.MonitorInactive(ctx, cInactive, offlineTimeout, nil)
	go func() {
		for inactive := range cInactive {
			if so.Logout(inactive) == nil {
				sign.MapUserSpace.Delete(inactive)
				if claims, ok := sign.MapUserClaims.Load(inactive); ok {
					lk.Log("delete token: [%v]", inactive)
					claims.(*u.UserClaims).DeleteToken()
					sign.MapUserClaims.Delete(inactive)
				}
				lk.Log("offline: [%v]", inactive)
			}
		}
	}()
}
