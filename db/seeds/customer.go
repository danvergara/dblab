package seeds

import "github.com/bxcodec/faker/v3"

// CustomerSeed seeds the database with customers.
func (s Seed) CustomerSeed() {
	for i := 0; i < 100; i++ {
		// prepare the statement.
		stmt, _ := s.db.Prepare(`INSERT INTO customers(name, email) VALUES ($1, $2)`)
		// execute query.
		_, err := stmt.Exec(faker.Name(), faker.Email())
		if err != nil {
			panic(err)
		}
	}
}
