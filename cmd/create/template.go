//go:generate go run github.com/valyala/quicktemplate/qtc -dir=./
package create

type TemplateInput struct {
	Package  string
	DontUndo bool
}
