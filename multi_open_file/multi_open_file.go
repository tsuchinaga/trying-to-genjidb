package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/genjidb/genji"
)

func main() {
	path := "multi_open.db"
	db, err := genji.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.Exec(`create table if not exists users`); err != nil {
		log.Fatalln(err)
	}
	_ = db.Close()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			openClose(path, i)
			wg.Done()
		}()
	}

	wg.Wait()

	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			log.Fatalln(err)
		}
	}
}

func openClose(path string, i int) {
	fmt.Printf("%d番目 start\n", i)
	db, err := genji.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d番目 open\n", i)

	time.Sleep(time.Second)

	_ = db.Close()
	fmt.Printf("%d番目 close\n", i)
	fmt.Printf("%d番目 end\n", i)
}
