package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func seedAccount(store Storage, firstName, lastName, email, password string) *Account {
	acc, err := NewAccount(firstName, lastName, email, password)
	if err != nil {
		log.Fatal(err)
	}	

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	return acc
}

func seedAccountsData(s Storage) {
	seedAccount(s, "caesar", "sage", "caesarsage@gmail.com", "sage12345")
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.init(); err != nil {
		log.Fatal(err)
	}

	// Seed Accounts
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	if *seed  {
		fmt.Println("Seeding the db...")
		seedAccountsData(store)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
