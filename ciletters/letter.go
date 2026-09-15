//go:build !solution

package ciletters

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed letter.tmpl
var letterTemplate string

func MakeLetter(n *Notification) (string, error) {
	tmpl, err := template.New("ci").Funcs(
		template.FuncMap{
			"isPipelineOK": func(status PipelineStatus) bool {
				return status == PipelineStatusOK
			},
			"trimHash": func(hash string) string {
				return hash[:8]
			},
			"last10lines": func(log string) []string {
				lines := strings.Split(log, "\n")
				var start int
				if len(lines) > 10 {
					start = len(lines) - 10
				} else {
					start = 0
				}
				return lines[start:]
			},
		}).Parse(letterTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		Notification
		PipelineStatusOK PipelineStatus
	}{
		Notification:     *n,
		PipelineStatusOK: PipelineStatusOK,
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, data)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
