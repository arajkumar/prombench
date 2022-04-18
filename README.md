[![ci](https://github.com/arajkumar/prombench/actions/workflows/ci.yaml/badge.svg)](https://github.com/arajkumar/prombench/actions/workflows/ci.yaml)

# prombench
prombench is a tiny program that sends some load to a prometheus implementations in the form of promql.

# Usage
prombench runs provided number of queries in the provided concurrency level and prints stats.

```
Usage: prombench [options...] <url> <query_file>
Options:
  -c  Number of workers to run concurrently. Total number of requests cannot
      be smaller than the concurrency level. Default is 50.
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
NumOfQueries: 11
TotalDuration:  257.912407ms
Min: 15.728243ms
Median: 22.763106ms
Average: 23.446582ms
Max: 31.884119ms

NumOfQueries: 11
TotalDuration:  257.912407ms
Min: 15.728243ms
Median: 22.763106ms
Average: 23.446582ms
Max: 31.884119ms
```
