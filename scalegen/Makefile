
.PHONY: test
test:
	go mod edit --replace github.com/spacemeshos/go-scale=../
	go test ./... 
	go mod edit --dropreplace github.com/spacemeshos/go-scale