OUT=./prombench
PKGS=$(shell go list ./... | grep -v /test/e2e)
BENCHMARK_DATA=https://github.com/timescale/promscale/raw/master/pkg/tests/testdata/real-dataset.sz

$(OUT):
	go build -o $(OUT) cmd/prombench.go 

.PHONY: test-unit
test-unit:
	go test -race -short $(PKGS) -count=1 -timeout 1m

.PHONY: fmt
fmt:
	go fmt ./...

real-dataset.sz:
	curl -s $(BENCHMARK_DATA) -o $@

.PHONY: promscale-up
promscale-up:
	@docker-compose up -d

.PHONY: promscale-down
promscale-down:
	@docker-compose down

.PHONY: setup-benchmark
setup-benchmark: promscale-up real-dataset.sz
	# TODO: Setup a container to ingest data and make it part of docker-compose service.
	@sleep 10
	@curl -v -H "Content-Type: application/x-protobuf" -H "Content-Encoding: snappy" -H "X-Prometheus-Remote-Write-Version: 0.1.0" --data-binary "@real-dataset.sz" "http://localhost:9201/write"

.PHONY: run-benchmark
run-benchmark: $(OUT) setup-benchmark pkg/parser/testdata/obs-queries.csv
	$(OUT) "http://localhost:9201" "pkg/parser/testdata/obs-queries.csv"
