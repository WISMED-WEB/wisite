package post

import (
	em "github.com/digisan/event-mgr"
)

var (
	edb = em.GetDB("./data")
	es  = em.NewEventSpan("MINUTE", edb.SaveEvtSpan)
)
