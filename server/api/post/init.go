package post

import (
	"context"

	em "github.com/digisan/event-mgr"
)

var (
	ctx      context.Context
	CancelES context.CancelFunc
)

func init() {
	ctx, CancelES = context.WithCancel(context.Background())
	em.InitDB("./data")
	em.InitEventSpan("MINUTE", ctx)
}
