package seeds

import (
	"log"

	"github.com/bxcodec/faker/v3"
	"github.com/danvergara/dblab/pkg/drivers"
)

// UserSeed seeds the database with users.
func (s Seed) UserSeed() {
	for i := 0; i < 100; i++ {
		var err error

		// execute query.
		switch s.driver {
		case drivers.POSTGRES:
			_, err = s.db.Exec(`INSERT INTO users(username) VALUES ($1)`, faker.Name())
		case drivers.MYSQL:
			_, err = s.db.Exec(`INSERT INTO users(username) VALUES (?)`, faker.Name())
		case drivers.SQLITE:
			_, err = s.db.Exec(`INSERT INTO users(username) VALUES (?)`, faker.Name())
		default:
			log.Println("unsupported driver")
		}

		if err != nil {
			log.Fatalf("error seeding users: %v", err)
		}
	}
}
