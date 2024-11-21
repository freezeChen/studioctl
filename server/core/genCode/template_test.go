package genCode

import (
	"bytes"
	"text/template"

	"testing"
)

func TestHtml(t *testing.T) {

	tt, err := template.New("test").Parse("<>")
	if err != nil {
		t.Error(err)
	}
	var b = bytes.Buffer{}
	tt.Execute(&b, nil)
	t.Log(b.String())
}
