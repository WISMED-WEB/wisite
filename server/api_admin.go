package main

import (
	"github.com/labstack/echo/v4"
	ad "github.com/wismed-web/wisite/server/api/admin"
)

func hookAdminHandler(r *echo.Group) {

	var mGET = map[string]echo.HandlerFunc{
		"/users":       ad.ListUser,
		"/onlineusers": ad.ListOnlineUser,
	}

	var mPOST = map[string]echo.HandlerFunc{
		"/activate-user": ad.ActivateUser,
	}

	var mPUT = map[string]echo.HandlerFunc{}

	var mDELETE = map[string]echo.HandlerFunc{}

	var mPATCH = map[string]echo.HandlerFunc{}

	// ------------------------------------------------------- //

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	mRegAPIs := map[string]map[string]echo.HandlerFunc{
		"GET":    mGET,
		"POST":   mPOST,
		"PUT":    mPUT,
		"DELETE": mDELETE,
		"PATCH":  mPATCH,
		// others...
	}

	mRegMethod := map[string]func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route{
		"GET":    r.GET,
		"POST":   r.POST,
		"PUT":    r.PUT,
		"DELETE": r.DELETE,
		"PATCH":  r.PATCH,
		// others...
	}

	for _, m := range methods {
		mAPI, method := mRegAPIs[m], mRegMethod[m]
		for path, handler := range mAPI {
			if handler == nil {
				continue
			}
			method(path, handler)
		}
	}
}
