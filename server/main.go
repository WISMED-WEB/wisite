package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
	su "github.com/digisan/user-mgr/sign-up"
	"github.com/digisan/user-mgr/udb"
	usr "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/postfinance/single"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/wismed-web/wisite/server/docs" // once `swag init`, comment it out
	"github.com/wismed-web/wisite/server/ws"
)

var (
	fHttp2 = false
)

func init() {
	lk.WarnDetail(false)
}

// @title WISMED WISITE API
// @version 1.0
// @description This is wismed wisite server.
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
		udb.OpenUserStorage("../data/user")
		su.SetValidator(nil)
	}

	// start Service
	done := make(chan string)
	echoHost(done)
	lk.Log(<-done)
}

func waitShutdown(e *echo.Echo) {
	go func() {
		defer udb.CloseUserStorage() // after closing echo, then close db

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

		// sign group
		{
			r := e.Group("/api/sign")
			hookSignHandler(r)
		}

		// admin group
		{
			r := e.Group("/api/admin")
			r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
				Claims:     &usr.UserClaims{},
				SigningKey: []byte(usr.TokenKey()),
			}))
			r.Use(ValidateToken)
			hookAdminHandler(r)
		}

		hookStatic(e)   // host static file/folder
		waitShutdown(e) // waiting for shutdown

		// web socket for message
		e.GET("/ws/msg", ws.WSMsg)

		// host swagger
		// http://localhost:1323/swagger/index.html
		e.GET("/swagger/*", echoSwagger.WrapHandler)

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
