all: test benchmark clean

test:
	go test ./...
	go test -race ./...

benchmark: test
	GOMAXPROCS=4 go test -race -bench=. -run=none -v -benchmem -benchtime=5s ./...

clean:
	find . -name '*.test' -delete
	find . -name '*.out' -delete


.PHONY: all test benchmark clean