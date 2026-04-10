BINARYDIR := ./pegsol/bin
BINARY_DFS := $(BINARYDIR)/pegsol-dfs
BINARY_BFS := $(BINARYDIR)/pegsol-bfs
CMD_DFS := ./pegsol/cmd/main.go
CMD_BFS := ./pegsol/cmd_bfs/main.go

.PHONY: build test clean

build:
	mkdir -p $(BINARYDIR)
	go build -o $(BINARY_DFS) $(CMD_DFS)
	go build -o $(BINARY_BFS) $(CMD_BFS)

test:
	go test ./...

clean:
	rm -f $(BINARY_DFS) $(BINARY_BFS)
