package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"runtime"

	csvparser "github.com/arajkumar/prombench/pkg/parser"
	plain "github.com/arajkumar/prombench/pkg/reporter"
	promqlworker "github.com/arajkumar/prombench/pkg/worker"
)

// Inspired from https://github.com/rakyll/hey/blob/master/hey.go

const (
	headerRegexp = `^([\w-]+):\s*(.+)`
)

var usage = `Usage: prombench [options...] <url> <query_file>
Options:
  -c  Number of workers to run concurrently. Total number of requests cannot
      be smaller than the concurrency level. Default is 50.
  -H  Custom HTTP header. You can specify as many as needed by repeating the flag.
      For example, -H "Accept: text/html" -H "Content-Type: application/xml" .
  -cpus                 Number of used cpu cores.
                        (default for current machine is %d cores)
`

var (
	concurrency = flag.Int("c", 50, "")
	cpus        = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}
	var hs headerSlice
	flag.Var(&hs, "H", "")

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
	csvParser, err := csvparser.NewCSVParser(inputReader)
	if err != nil {
		errAndExit("Unable to start CSV parser, err %s", err)
	}

	header := make(http.Header)
	// set any other additional repeatable headers
	for _, h := range hs {
		match, err := parseInputWithRegexp(h, headerRegexp)
		if err != nil {
			usageAndExit(err.Error())
		}
		header.Set(match[1], match[2])
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

	w, err := promqlworker.NewPromQLWorker(*hostUrl, promqlworker.WithHeaders(header), promqlworker.WithConcurrency(*concurrency))
	if err != nil {
		errAndExit("Unable to start worker, err %s", err)
	}
	report := w.Run(ctx, csvParser.Queries())

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

func parseInputWithRegexp(input, regx string) ([]string, error) {
	re := regexp.MustCompile(regx)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 1 {
		return nil, fmt.Errorf("could not parse the provided input; input = %v", input)
	}
	return matches, nil
}

type headerSlice []string

func (h *headerSlice) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headerSlice) Set(value string) error {
	*h = append(*h, value)
	return nil
}
