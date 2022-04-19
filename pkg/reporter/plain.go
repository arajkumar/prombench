package plain

import (
	"io"
	"text/template"

	"github.com/arajkumar/prombench"
)

type Plain struct {
}

var (
	defaultTmpl = `Summary:
  NumOfQueries: {{ .NumOfQueries }}
  TotalDuration: {{ .TotalDuration }}
  Min: {{ .Min }}
  Median: {{ .Median }}
  Average: {{ .Avg }}
  Max: {{ .Max }}

Status code distribution:{{ range $code, $num := .StatusCodeDist }}
  [{{ $code }}]	{{ $num }} responses{{ end }}

{{ if gt (len .ErrDist) 0 }}Error distribution:{{ range $err, $num := .ErrDist }}
  [{{ $num }}]	{{ $err }}{{ end }}{{ end }}
`
)

// Implements prombench.Reporter interface.
func (p Plain) Report(out io.Writer, s prombench.Summary) error {
	tmpl := template.Must(template.New("tmpl").Parse(defaultTmpl))
	tmpl.Execute(out, s)
	return nil
}
