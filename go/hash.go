package main

import (
	"fmt"
	"hash"
	"log"

	"github.com/minio/highwayhash"
)


func highwayhash_hh() (hash.Hash) {

	key := make([]byte, 32)

	hash, err := highwayhash.New(key)

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return hash
}