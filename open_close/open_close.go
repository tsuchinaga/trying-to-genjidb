package main

import (
	"log"
	"os"
	"time"

	"github.com/genjidb/genji/document"

	"github.com/genjidb/genji/types"

	"github.com/genjidb/genji"
)

func main() {
	intermittent()
	continuous()
}

var path = "open_close.db"
var cnt = 20

func continuous() {
	start := time.Now()

	db, err := genji.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < cnt; i++ {
		if err := execute(db); err != nil {
			log.Fatalln(err)
		}
	}

	if err := db.Close(); err != nil {
		log.Fatalln(err)
	}

	end := time.Now()
	log.Println("time:", end.Sub(start))

	if err := os.Remove(path); err != nil {
		log.Fatalln(err)
	}
}

func intermittent() {
	start := time.Now()

	for i := 0; i < cnt; i++ {
		db, err := genji.Open(path)
		if err != nil {
			log.Fatalln(err)
		}

		if err := execute(db); err != nil {
			log.Fatalln(err)
		}

		if err := db.Close(); err != nil {
			log.Fatalln(err)
		}
	}

	end := time.Now()
	log.Println("time:", end.Sub(start))

	if err := os.Remove(path); err != nil {
		log.Fatalln(err)
	}
}

func execute(db *genji.DB) error {
	if err := db.Exec("create table if not exists hoge"); err != nil {
		return err
	}

	hoge := Hoge{
		Name:   "testing",
		Hoge:   "HogeHoge",
		Age:    100_000,
		Cash:   1.234,
		Active: true,
		Fuga: struct {
			Name string
			Fuga bool
		}{
			Name: "FugaFuga",
			Fuga: false,
		},
		Piyo: []struct {
			Name string
			Piyo bool
		}{
			{Name: "PiyoPiyo1", Piyo: false},
			{Name: "PiyoPiyo2", Piyo: true},
			{Name: "PiyoPiyo3", Piyo: true},
		},
	}

	if err := db.Exec("insert into hoge values ?", hoge); err != nil {
		return err
	}

	res, err := db.Query("select * from hoge")
	if err != nil {
		return err
	}
	if err := res.Iterate(func(d types.Document) error {
		var h Hoge
		if err := document.StructScan(d, &h); err != nil {
			return err
		}
		//log.Printf("%+v\n", h)
		return nil
	}); err != nil {
		return err
	}
	if err := res.Close(); err != nil {
		return err
	}

	if err := db.Exec("drop table hoge"); err != nil {
		return err
	}

	return nil
}

type Hoge struct {
	Name   string
	Hoge   string
	Age    int
	Cash   float64
	Active bool
	Fuga   struct {
		Name string
		Fuga bool
	}
	Piyo []struct {
		Name string
		Piyo bool
	}
}
