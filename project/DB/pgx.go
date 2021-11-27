package DB

import (
	"context"
	"main.go/project/addition"
	"sync"
)

func StartDB(ctx context.Context, wg *sync.WaitGroup) chan addition.Query {
	db := MakeDataBase(ctx)
	queChan := db.QueryHandler(wg)

	return queChan
}
