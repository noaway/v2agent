package tpl

import (
	"bytes"
	"text/template"
)

// template var
var (
	tplmap = map[string]Template{}
)

// Template interface
type Template interface {
	Execute(interface{}) (string, error)
}

// Parse fn
func Parse(name, text string, funcMap map[string]interface{}) Template {
	t, ok := tplmap[name]
	if !ok {
		tp := new(tpl)
		tp.Template = template.Must(template.New(name).Parse(text))
		if funcMap != nil {
			tp.Funcs(funcMap)
		}
		t = tp
		tplmap[name] = t
	}
	return t
}

type tpl struct {
	*template.Template
}

func (t *tpl) Execute(data interface{}) (string, error) {
	buf := bytes.Buffer{}
	if t.Template == nil {
		return "", nil
	}
	err := t.Template.Execute(&buf, data)
	return buf.String(), err
}
