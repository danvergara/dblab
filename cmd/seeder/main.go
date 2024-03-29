package main

import (
	"flag"
	"log"
	"os"

	// mysql driver.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/danvergara/dblab/db/seeds"
	"github.com/danvergara/dblab/pkg/config"

	// postgres driver.
	_ "github.com/lib/pq"

	// sqlite driver.
	_ "modernc.org/sqlite"
)

func main() {
	handleArgs()
}

func handleArgs() {
	flag.Parse()
	args := flag.Args()
	cfg := config.Get()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			connString := cfg.GetSQLXDBConnStr()
			db, err := sqlx.Open(cfg.Driver, connString)
			if err != nil {
				log.Fatalf("Error opening DB: %v", err)
			}
			seeds.Execute(db, cfg.Driver, args[1:]...)
			os.Exit(0)
		}
	}
}
