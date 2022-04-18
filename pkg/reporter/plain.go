package plain

import (
	"io"
	"text/template"

	"github.com/arajkumar/prombench"
)

type Plain struct {
}

var (
	defaultTmpl = `NumOfQueries: {{ .NumOfQueries }}
TotalDuration:	{{ .TotalDuration }}
Min: {{ .Min }}
Median: {{ .Median }}
Average: {{ .Avg }}
Max: {{ .Max }}
`
)

// Implements prombench.Reporter interface.
func (p Plain) Report(out io.Writer, s prombench.Summary) error {
	tmpl := template.Must(template.New("tmpl").Parse(defaultTmpl))
	tmpl.Execute(out, s)
	return nil
}
