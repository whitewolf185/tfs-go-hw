package DB

import (
	"context"
	"github.com/whitewolf185/fs-go-hw/project/repository/DB/add_DB"
	"sync"
)

func StartDB(ctx context.Context, wg *sync.WaitGroup) chan add_DB.Query {
	db := MakeDataBase(ctx)
	queChan := db.QueryHandler(wg)

	return queChan
}
