package seeds

import (
	"log"
	"reflect"

	// mysql driver.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	// postgres driver.
	_ "github.com/lib/pq"
)

// Seed type.
type Seed struct {
	db *sqlx.DB
}

// Execute will executes the given seeder method.
func Execute(db *sqlx.DB, seedMethodNames ...string) {
	s := Seed{
		db: db,
	}

	seedType := reflect.TypeOf(s)

	// Executes all seeders if no method is given.
	if len(seedMethodNames) == 0 {
		log.Println("running all seeder...")
		// We are looping over the method on a Seed struct.
		for i := 0; i < seedType.NumMethod(); i++ {
			// Get the method in the current iteration.
			method := seedType.Method(i)
			// Execute seeder.
			seed(s, method.Name)
		}
	}

	// Execute only the given method names
	for _, item := range seedMethodNames {
		seed(s, item)
	}
}

func seed(s Seed, seedMethodName string) {
	// Get the reflect value of the method.
	m := reflect.ValueOf(s).MethodByName(seedMethodName)
	// Exit if the method diesn't exist.
	if !m.IsValid() {
		log.Fatal("no method called", seedMethodName)
	}

	// Execute the method.
	log.Println("seeding", seedMethodName, "...")
	m.Call(nil)
	log.Println("seed", seedMethodName, "succeded")
}
