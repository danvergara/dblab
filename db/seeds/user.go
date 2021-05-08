package seeds

import (
	"log"

	"github.com/bxcodec/faker/v3"
)

// UserSeed seeds the database with users.
func (s Seed) UserSeed() {
	for i := 0; i < 100; i++ {
		// execute query.
		_, err := s.db.Exec(`INSERT INTO users(username) VALUES ($1)`, faker.Name())
		if err != nil {
			log.Fatalf("error seeding users: %v", err)
		}
	}
}
