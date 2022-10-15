package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	fm "github.com/digisan/file-mgr"
	cfg "github.com/digisan/go-config"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
	r "github.com/digisan/user-mgr/relation"
	u "github.com/digisan/user-mgr/user"
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
	port   = 1323
)

func init() {
	lk.WarnDetail(false)

	cfg.Init("main", false, "./config.json")

	fHttp2 = cfg.Val[bool]("http2")
	port = cfg.Val[int]("port")
}

// @title WISMED WISITE API
// @version 1.0
// @description WISMED Wisite-API Server. Updated@ 2022-09-03T11:33:31+10:00
// @termsOfService
// @contact.name API Support
// @contact.url
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:1323
// @BasePath
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name authorization
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

	// start Service
	done := make(chan string)
	echoHost(done)
	lk.Log(<-done)
}

func waitShutdown(e *echo.Echo) {
	go func() {
		defer u.CloseDB()         // after closing echo, close user db, i.e. deactivate ***[UserDB]***
		defer r.CloseDB()         // after closing echo, close relation db, i.e. deactivate ***[RelDB]***
		defer fm.DisposeFileMgr() // close file db

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

		// host static file/folder | only for testing
		hookStatic(e)

		// web socket
		e.GET("/ws/msg", ws.WSMsg)

		// host swagger http://localhost:1323/swagger/index.html
		e.GET("/swagger/*", echoSwagger.WrapHandler)

		// sign group without JWT
		{
			api.SignHandler(e.Group("/api/sign"))
			api.SystemHandler(e.Group("/api/system"))
			api.DebugHandler(e.Group("/api/debug"))
		}

		// other groups with JWT
		groups := []string{
			"/api/sign-out",
			"/api/admin",
			"/api/file",
			"/api/post",
			"/api/user",
			"/api/rel",
			"/api/client",
		}
		handlers := []func(*echo.Group){
			api.SignoutHandler,
			api.AdminHandler,
			api.FileHandler,
			api.PostHandler,
			api.UserHandler,
			api.RelHandler,
			api.ClientHandler,
		}
		for i, group := range groups {
			r := e.Group(group)
			r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
				Claims:     &u.UserClaims{},
				SigningKey: []byte(u.TokenKey()),
			}))
			r.Use(ValidateToken)
			handlers[i](r)
		}

		// running...
		portstr := fmt.Sprintf(":%d", port)
		var err error
		if fHttp2 {
			err = e.StartTLS(portstr, "./cert/public.pem", "./cert/private.pem")
		} else {
			err = e.Start(portstr)
		}
		lk.FailOnErrWhen(err != http.ErrServerClosed, "%v", err)
	}()
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userTkn := c.Get("user").(*jwt.Token)
		claims := userTkn.Claims.(*u.UserClaims)
		if claims.ValidateToken(userTkn.Raw) {
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"message": "invalid or expired jwt",
		})
	}
}
