[![ci](https://github.com/arajkumar/prombench/actions/workflows/ci.yaml/badge.svg)](https://github.com/arajkumar/prombench/actions/workflows/ci.yaml)

# prombench
prombench is a tiny program that sends some load to a [prometheus](https://prometheus.io) implementations in the form of [promql](https://prometheus.io/docs/prometheus/latest/querying/basics).

# Usage
prombench runs provided number of queries in the provided concurrency level and prints stats.

```
Usage: prombench [options...] <url> <query_file>
Options:
  -c  Number of workers to run concurrently. Total number of requests cannot
      be smaller than the concurrency level. Default is 50.
  -H  Custom HTTP header. You can specify as many as needed by repeating the flag.
      For example, -H "Accept: text/html" -H "Content-Type: application/xml" .
  -cpus                 Number of used cpu cores.
                        (default for current machine is 12 cores)
```

# For Developers

To build,

```
make prombench
```

To run tests,

```
make test-unit
```

# To benchmark against Promscale instance
The following make target would setup promscale instance and ingests [sample data](pkg/parser/testdata/obs-queries.csv) to execute benchmark.

```
make run-benchmark
```

# Sample benchmark result
```
$ ./prombench "http://localhost:9201" "pkg/parser/testdata/obs-queries.csv"
Summary:
  NumOfQueries: 11
  TotalDuration: 215.426317ms
  Min: 14.703164ms
  Median: 18.428764ms
  Average: 19.58421ms
  Max: 26.788814ms

Status code distribution:
  [200] 11 responses
```
