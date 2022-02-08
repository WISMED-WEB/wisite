package ws

import (
	"context"
	"fmt"
	"sync"

	"github.com/labstack/echo/v4"
	lk "github.com/digisan/logkit"
	"golang.org/x/net/websocket"
)

// after implementing, register with path in 'api_reg.go'

var (
	mIdMsg      = &sync.Map{}
	mIdWSCancel = &sync.Map{}
)

func SendMsg(id string, msg interface{}) bool {
	chMsg, ok := mIdMsg.Load(id)
	if !ok {
		return false
	}
	chMsg.(chan interface{}) <- msg
	return true
}

func BroadCast(msg interface{}) {
	mIdMsg.Range(func(id, msg interface{}) bool {
		go SendMsg(id.(string), msg)
		return true
	})
}

func CloseMsg(id string) bool {
	chCancel, ok := mIdWSCancel.Load(id)
	if !ok {
		return false
	}
	chCancel.(context.CancelFunc)()
	return true
}

func CloseAllMsg() {
	mIdWSCancel.Range(func(id, chCancel interface{}) bool {
		go chCancel.(context.CancelFunc)()
		return true
	})
}

// Activate WS Msg by GET
func WSMsg(c echo.Context) error {

	id := c.Request().Header.Get("id")
	id = "id" // just for testing ***********************************

	// reg a new message channel
	mIdMsg.Store(id, make(chan interface{}, 1024))

	// reg message channel closing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mIdWSCancel.Store(id, cancel)

	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// Read
		msg := ""
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			c.Logger().Error(err)
			return
		}
		fmt.Printf("%s\n", msg)

		done := make(chan struct{})
		go func(ctx context.Context, done chan<- struct{}) {
			defer func() { done <- struct{}{} }()
			if IChMsg, ok := mIdMsg.Load(id); ok {
				chMsg := IChMsg.(chan interface{})
				for {
					select {
					case msg := <-chMsg:
						lk.WarnOnErr("%v", websocket.Message.Send(ws, fmt.Sprintf("From WS Server! --- %v", msg)))
					case <-ctx.Done():
						return
					}
				}
			}
		}(ctx, done)
		<-done

	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
