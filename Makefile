test:
	go test ./...
	go test -race ./...

build: test
	go build -v
	
benchmark: test
	go test -race -bench=. -run=none -v -benchmem -benchtime=5s ./...

clean:
	find . -name '*.test' -delete
	find . -name '*.out' -delete

.PHONY: test build benchmark clean