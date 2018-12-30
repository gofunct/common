package templates

var TemplateConfig = MustCreateTemplate("config", `package {{.Name}}

type Config struct {
}
`)
