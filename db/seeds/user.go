package seeds

import (
	"github.com/bxcodec/faker/v3"
)

// UserSeed seeds the database with users.
func (s Seed) UserSeed() {
	for i := 0; i < 100; i++ {

		// prepare the statement.
		stmt, _ := s.db.Prepare(`INSERT INTO users(username) VALUES ($1)`)
		// execute query.
		_, err := stmt.Exec(faker.Name())
		if err != nil {
			panic(err)
		}
	}
}
