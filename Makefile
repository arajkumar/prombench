OUT=prombench
PKGS=$(shell go list ./... | grep -v /test/e2e)

$(OUT):
	go build -o $(OUT) cmd/prombench.go 

.PHONY: test-unit
test-unit:
	go test -race -short $(PKGS) -count=1 -timeout 1m

.PHONY: fmt
fmt:
	go fmt ./...

