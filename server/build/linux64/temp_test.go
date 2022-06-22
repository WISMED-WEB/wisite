package linux64

import (
	"context"
	"fmt"
	"testing"
	"time"

	em "github.com/digisan/event-mgr"
	lk "github.com/digisan/logkit"
)

func TestGetAllEvtIds(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	em.InitDB("./data")
	defer em.CloseDB()

	em.InitEventSpan("MINUTE", ctx)

	ids, err := em.FetchAllEvtIDs() // GetEvtIdAllDB()
	if err != nil {
		panic(err)
	}

	for i, eid := range ids {
		fmt.Println(i, eid)
	}
	fmt.Println("------> total:", len(ids))
}

func TestTemp(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	em.InitDB("./data")
	defer em.CloseDB()

	em.InitEventSpan("MINUTE", ctx)

	ids, err := em.FetchOwn("cdutwhu", "202206")
	lk.WarnOnErr("%v", err)

	for i, eid := range ids {
		fmt.Println(i, eid)
	}
	fmt.Println("------> total:", len(ids))

	time.Sleep(1 * time.Second)
}
