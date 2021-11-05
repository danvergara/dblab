package seeds

import (
	"log"

	"github.com/bxcodec/faker/v3"
)

// CustomerSeed seeds the database with customers.
func (s Seed) CustomerSeed() {
	for i := 0; i < 100; i++ {
		var err error

		// execute query.
		switch s.driver {
		case "postgres":
			_, err = s.db.Exec(`INSERT INTO customers(name, email) VALUES ($1, $2)`, faker.Name(), faker.Email())
		case "mysql":
			_, err = s.db.Exec(`INSERT INTO customers(name, email) VALUES (?, ?)`, faker.Name(), faker.Email())
		case "sqlite3":
			_, err = s.db.Exec(`INSERT INTO customers(name, email) VALUES (?, ?)`, faker.Name(), faker.Email())
		default:
			log.Println("unsupported driver")
		}

		if err != nil {
			log.Fatalf("error seeding customers: %v", err)
		}
	}
}
