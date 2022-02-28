package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	fm "github.com/digisan/file-mgr"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
	su "github.com/digisan/user-mgr/sign-up"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	vf "github.com/digisan/user-mgr/user/valfield"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/postfinance/single"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/wismed-web/wisite-api/server/api"
	_ "github.com/wismed-web/wisite-api/server/docs" // once `swag init`, comment it out
	"github.com/wismed-web/wisite-api/server/ws"
)

var (
	fHttp2 = false
)

func init() {
	lk.WarnDetail(false)
}

// @title WISMED WISITE API
// @version 1.0
// @description This is wismed wisite-api server.
// @termsOfService
// @contact.name API Support
// @contact.url
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:1323
// @BasePath
func main() {

	http2Ptr := flag.Bool("http2", false, "http2 mode?")
	flag.Parse()
	fHttp2 = *http2Ptr

	// only one instance
	const dir = "./tmp-locker"
	gio.MustCreateDir(dir)
	one, err := single.New("echo-service", single.WithLockPath(dir))
	lk.FailOnErr("%v", err)
	lk.FailOnErr("%v", one.Lock())
	defer func() {
		lk.FailOnErr("%v", one.Unlock())
		os.RemoveAll(dir)
		lk.Log("Server Exited Successfully")
	}()

	// other init actions
	{
		// set user db dir, activate ***[udb.UserDB]***
		udb.OpenUserStorage("./data/db-user")

		// set user validator
		su.SetValidator(map[string]func(interface{}) bool{
			vf.AvatarType: func(i interface{}) bool {
				return i == "" || strings.HasPrefix(i.(string), "image/")
			},
		})

		// set user file space & file item db space
		fm.SetFileMgrRoot("./data/user-space", "./data/db-fileitem")
	}

	// start Service
	done := make(chan string)
	echoHost(done)
	lk.Log(<-done)
}

func waitShutdown(e *echo.Echo) {
	go func() {
		defer udb.CloseUserStorage() // after closing echo, then close db, deactivate ***[udb.UserDB]***

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		lk.Log("Got Ctrl+C")

		// other clean-up before closing echo
		{
			ws.BroadCast("backend service shutting down...") // testing
			ws.CloseAllMsg()
		}

		// shutdown echo
		lk.FailOnErr("%v", e.Shutdown(ctx)) // close echo at e.Shutdown
	}()
}

func echoHost(done chan<- string) {
	go func() {
		defer func() { done <- "Echo Shutdown Successfully" }()

		e := echo.New()
		defer e.Close()

		// Middleware
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.BodyLimit("2G"))
		// CORS
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))

		// waiting for shutdown
		waitShutdown(e)

		// host static file/folder
		hookStatic(e)

		// web socket
		e.GET("/ws/msg", ws.WSMsg)

		// host swagger http://localhost:1323/swagger/index.html
		e.GET("/swagger/*", echoSwagger.WrapHandler)

		// sign group without JWT
		{
			r := e.Group("/api/sign")
			api.SignHandler(r)
		}

		// other groups with JWT
		groups := []string{"/api/sign-out", "/api/admin", "api/file", "api/user"}
		handlers := []func(*echo.Group){
			api.SignoutHandler,
			api.AdminHandler,
			api.FileHandler,
			api.UserHandler,
		}
		for i, group := range groups {
			r := e.Group(group)
			r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
				Claims:     &usr.UserClaims{},
				SigningKey: []byte(usr.TokenKey()),
			}))
			r.Use(ValidateToken)
			handlers[i](r)
		}

		// running...
		var err error
		if fHttp2 {
			err = e.StartTLS(":1323", "./cert/public.pem", "./cert/private.pem")
		} else {
			err = e.Start(":1323")
		}
		lk.FailOnErrWhen(err != http.ErrServerClosed, "%v", err)
	}()
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userTkn := c.Get("user").(*jwt.Token)
		claims := userTkn.Claims.(*usr.UserClaims)
		if claims.ValidateToken(userTkn.Raw) {
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "invalid or expired jwt",
		})
	}
}
