package DB

import (
	"context"
	"sync"

	"github.com/whitewolf185/fs-go-hw/project/repository/DB/addDB"
)

func StartDB(ctx context.Context, wg *sync.WaitGroup) chan addDB.Query {
	db := MakeDataBase(ctx)
	queChan := db.QueryHandler(wg)

	return queChan
}
