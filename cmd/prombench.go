package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"runtime"

	csvparser "github.com/arajkumar/prombench/pkg/parser"
	promqlworker "github.com/arajkumar/prombench/pkg/worker"
)

var usage = `Usage: prombench [options...] <url> <query_file>
Options:
  -c  Number of workers to run concurrently. Total number of requests cannot
      be smaller than the concurrency level. Default is 50.
  -cpus                 Number of used cpu cores.
                        (default for current machine is %d cores)
`

var (
	c    = flag.Int("c", 50, "")
	cpus = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}
	flag.Parse()
	if flag.NArg() < 2 {
		usageAndExit("")
	}

	runtime.GOMAXPROCS(*cpus)
	host := flag.Args()[0]
	hostUrl, err := url.Parse(host)
	if err != nil {
		errAndExit("Invalid host format %s, err %s", host, err)
	}

	input := flag.Args()[1]

	inputReader, err := os.Open(input)
	if err != nil {
		errAndExit("Unable to open file %s, err %s", input, err)
	}
	csvParser, err := csvparser.NewCSVParser(inputReader, csvparser.WithConcurrency(*c))
	if err != nil {
		errAndExit("Unable to start CSV parser, err %s", err)
	}

	// start parsing loop
	go csvParser.Parse()

	w, err := promqlworker.NewPromQLWorker()
	if err != nil {
		errAndExit("Unable to start worker, err %s", err)
	}
	ctx := context.Background()
	report, err := w.Run(ctx, hostUrl, csvParser.Queries())
	if err != nil {
		errAndExit("Worker failed with err %s", err)
	}
	fmt.Println(report.ToSummary())
}

func errAndExit(format string, msg ...interface{}) {
	fmt.Fprintf(os.Stderr, format, msg...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
