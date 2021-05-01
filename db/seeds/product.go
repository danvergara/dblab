package seeds

import (
	"math/rand"

	"github.com/bxcodec/faker/v3"
)

// ProductSeed seeds product data.
func (s Seed) ProductSeed() {
	for i := 0; i < 100; i++ {
		// prepare the statement.
		stmt, _ := s.db.Prepare(`INSERT INTO products(name, price) VALUES ($1, $2)`)
		// execute query.
		_, err := stmt.Exec(faker.Word(), rand.Float32())
		if err != nil {
			panic(err)
		}
	}
}
