package main

import (
	"fmt"
	"github.com/my-dev-workspace/go-pkg/database"
)

type Product struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func main() {
	db := database.NewDB(database.NewDBConfigWith("127.0.0.1", 3306, "clients_b2b", "dev", "1111"))

	// find one row in the database and load it
	// into a struct variable
	var row Product
	err := db.
		Select("id,name").
		From("products").
		Where(database.Eq("id", 70465)).
		GetRow(&row)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", row)

}
