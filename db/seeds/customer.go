package seeds

import (
	"log"

	"github.com/bxcodec/faker/v3"
)

// CustomerSeed seeds the database with customers.
func (s Seed) CustomerSeed() {
	for i := 0; i < 100; i++ {
		// execute query.
		_, err := s.db.Exec(`INSERT INTO customers(name, email) VALUES ($1, $2)`, faker.Name(), faker.Email())
		if err != nil {
			log.Fatalf("error seeding customers: %v", err)
		}
	}
}
