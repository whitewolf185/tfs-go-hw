package DB

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/whitewolf185/fs-go-hw/project/addition/MyErrors"
	"github.com/whitewolf185/fs-go-hw/project/addition/add_DB"
)

const OrderQuery = "INSERT INTO operations (ordertime,ticket,type, size, limitprice) VALUES (now(),$1, $2, $3, $4);"

type DataBase struct {
	conn *pgx.Conn

	ctx    context.Context
	cancel context.CancelFunc
}

func MakeDataBase(ctx context.Context) *DataBase {
	var db DataBase
	db.ctx, db.cancel = context.WithCancel(ctx)
	err := db.Connect()
	if err != nil {
		MyErrors.DBConnectionErr(err)
		return nil
	}

	return &db
}

func (db *DataBase) Close() {
	err := db.conn.Close(db.ctx)
	if err != nil {
		MyErrors.DBCloseConnErr(err)
	}
	time.Sleep(time.Second)
	db.cancel()
}

func (db *DataBase) Connect() error {
	// TODO можно потестировать подключение к БД
	var err error
	dsn := "postgres://postgres:kukuruza@localhost:5433/postgres"

	db.conn, err = pgx.Connect(db.ctx, dsn)
	if err != nil {
		return err
	}

	if err = db.conn.Ping(db.ctx); err != nil {
		return err
	}

	return nil
}

func (db *DataBase) QueryHandler(wg *sync.WaitGroup) chan add_DB.Query {
	queChan := make(chan add_DB.Query)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-db.ctx.Done():
				close(queChan)
				return

			case query := <-queChan:
				_, err := db.conn.Exec(db.ctx, OrderQuery, query.Ticket, query.Type, query.Size, query.LimitPrice)
				if err != nil {
					MyErrors.DBExecErr(err)
				}
			}
		}
	}()

	return queChan
}
