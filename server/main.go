package main

import "github.com/simt/dtacc"

func main() {
	db, err := dtacc.NewDB()
	if err != nil {
		panic(err)
	}

	_ = db
}
