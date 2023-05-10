package main

import "prushka/internal/server"

func main() {
	s := server.New()
	s.Run()
}
