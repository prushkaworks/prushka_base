GC=go
UP=docker start 
IMAGE_NAME=some-postgres

test:
	$(GC) test -v ./internal/server

run: postgres
	$(GC) run ./cmd/prushka/main.go

postgres:
	$(UP) $(IMAGE_NAME)