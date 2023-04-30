package seeds

import (
	"log"

	"github.com/bxcodec/faker/v3"
	"github.com/danvergara/dblab/pkg/drivers"
)

// CustomerSeed seeds the database with customers.
func (s Seed) CustomerSeed() {
	for i := 0; i < 100; i++ {
		var err error

		// execute query.
		switch s.driver {
		case drivers.Postgres:
			_, err = s.db.Exec(`INSERT INTO customers(name, email) VALUES ($1, $2)`, faker.Name(), faker.Email())
		case drivers.MySQL:
			_, err = s.db.Exec(`INSERT INTO customers(name, email) VALUES (?, ?)`, faker.Name(), faker.Email())
		case drivers.SQLite:
			_, err = s.db.Exec(`INSERT INTO customers(name, email) VALUES (?, ?)`, faker.Name(), faker.Email())
		default:
			log.Println("unsupported driver")
		}

		if err != nil {
			log.Fatalf("error seeding customers: %v", err)
		}
	}
}
