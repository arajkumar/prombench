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
make unit-test
```

# To benchmark against Promscale instance

```
make benchmark
```
