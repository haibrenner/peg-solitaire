BINARYDIR := ./pegsol/bin
BINARY := $(BINARYDIR)/pegsol
CMD := ./pegsol/cmd/main.go

.PHONY: build test clean

build:
	mkdir -p $(BINARYDIR)
	go build -o $(BINARY) $(CMD)

test:
	go test ./...

clean:
	rm -f $(BINARY)
