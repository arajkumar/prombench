package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"runtime"

	csvparser "github.com/arajkumar/prombench/pkg/parser"
	plain "github.com/arajkumar/prombench/pkg/reporter"
	promqlworker "github.com/arajkumar/prombench/pkg/worker"
)

// Inspired from https://github.com/rakyll/hey/blob/master/hey.go
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

	// copied from https://medium.com/@matryer/make-ctrl-c-cancel-the-context-context-bd006a8ad6ff
	ctx := context.Background()
	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	// start parsing loop
	go csvParser.Parse(ctx)

	w, err := promqlworker.NewPromQLWorker()
	if err != nil {
		errAndExit("Unable to start worker, err %s", err)
	}
	report, err := w.Run(ctx, *hostUrl, csvParser.Queries())
	if err != nil {
		errAndExit("Worker failed with err %s", err)
	}

	plain := plain.Plain{}
	plain.Report(os.Stdout, report.ToSummary())
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
