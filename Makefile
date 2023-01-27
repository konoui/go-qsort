BASE_DIR := $(shell pwd)

test:
	GODEBUG=cgocheck=0 \
	go test -v ./...

run:
	go run cmd/main.go > $(BASE_DIR)/go-snapshot.txt
	(cd c && gcc -o main main.c psort.c qsort.c && ./main > $(BASE_DIR)/c-snapshot.txt)

diff: run
	@(echo "---------------------------------------")
	@(diff $(BASE_DIR)/go-snapshot.txt $(BASE_DIR)/c-snapshot.txt)
	@(echo "---------------------------------------")

cover:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html
