package test

import (
	"go-learning/pool"
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	//要使用的goroutine数量
	maxGoroutines = 25
	//池中的资源数量
	pooledResourdes = 2
)

type dbConnection struct {
	ID int32
}

func (dbConn *dbConnection) Close() error {
	log.Println("Close: Connection", dbConn.ID)
	return nil
}

var idCounter int32

func createConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Println("Create: New Connection", id)
	return &dbConnection{id}, nil
}

func TestPool(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(maxGoroutines)
	//创建用来管理连接的池
	p, err := pool.New(createConnection, pooledResourdes)
	if err != nil {
		log.Println(err)
	}

	for query := 0; query < maxGoroutines; query++ {
		go func(q int) {
			performQueries(q, p)
			wg.Done()
		}(query)
	}

	wg.Wait()
	log.Println("Shutdown Program.")
	p.Close()
}

func performQueries(query int, p *pool.Pool) {
	conn, err := p.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	defer p.Release(conn)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("QID[%d] CID[%d]\n", query, conn.(*dbConnection).ID)
}
