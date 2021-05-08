package seeds

import (
	"log"
	"math/rand"

	"github.com/bxcodec/faker/v3"
)

// ProductSeed seeds product data.
func (s Seed) ProductSeed() {
	for i := 0; i < 100; i++ {
		// execute query.
		_, err := s.db.Exec(`INSERT INTO products(name, price) VALUES ($1, $2)`, faker.Word(), rand.Float32())
		if err != nil {
			log.Fatalf("error seeding products: %v", err)
		}
	}
}
