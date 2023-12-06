package main

import (
	"fmt"
	"log"
	"url-shortener/internal/config"
)

func main() {

	error_msg, cfg := config.LoadCFG()
	if cfg == nil {
		log.Fatal(error_msg)
	}

	// _rnl_
	fmt.Println(cfg)

	// TODO: init logger: slog

	// TODO: init storage: sqlite

	// TODO: init router: chi, "chi render"

	// TODO: run server

}
