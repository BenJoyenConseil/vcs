package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"ugit/data"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Use one of the following commands : init | commit | hash_object")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		if len(os.Args) > 2 {
			path, _ := filepath.Abs(os.Args[2])
			data.UInit(path)
		} else {
			path, _ := os.Getwd()
			data.UInit(path)
		}
		os.Exit(0)
	case "cat":
		log.Println("cat object file ", os.Args[2])
		log.Println(data.GetObject(os.Args[2]))
		os.Exit(0)
	case "hash_object":
		if len(os.Args) > 2 {
			log.Println("hashing the content", os.Args[2])
			log.Println(data.PutObject(os.Args[2], data.BLOB))
			os.Exit(0)
		}
	case "commit":
		log.Println(data.Commit("./", os.Args[2]))
		os.Exit(0)
	default:
		log.Println("Usage Ugit")
	}
	os.Exit(1)
}
