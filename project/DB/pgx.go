package DB

import (
	"context"
	"sync"

	"github.com/whitewolf185/fs-go-hw/project/addition/add_DB"
)

func StartDB(ctx context.Context, wg *sync.WaitGroup) chan add_DB.Query {
	db := MakeDataBase(ctx)
	queChan := db.QueryHandler(wg)

	return queChan
}
