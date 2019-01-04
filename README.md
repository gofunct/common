# ![](https://github.com/gofunct/common/blob/master/logo/dark_logo_transparent_background.png?raw=true)  

## Info

Author: Coleman Word 
Email: coleman.word@gofunct.com
Repo Name: Common
Language: Golang
Download: `go get github.com/gofunct/common/...`

Summary: This is a common library I use for other projects. Some of the libraries may be moved to their own
repositories in the future.

Notable:
- logging-----> logging/
- errors-----> errors/
- ui -----> ui/
- fies -----> files/
- cobra -----> app/
- viper -----> app/
- google cloud -----> google/
- aws -----> aws/
- grpc gateway -----> runtime/

## Makefile

Input: `make help`
```commandline
Colemans-MacBook-Pro:common coleman$ make help
all                            generate binaries to bin/
clean                          remove all binaries in bin/
cover                          run test coverage
format                         go format entire directory
gen                            go generate entire project (wire, vsfgen, mockgen, protoc)
help                           this help
lint                           lint with reviewdog
packages                       generate packages
setup                          setup with gex
test                           run all project tests


```
## File Tree

```commandline
├── LICENSE
├── Makefile
├── README.md
├── app
│   ├── app.go
│   ├── bucket.go
│   ├── flags.go
│   ├── healthcheck.go
│   ├── inject_aws.go
│   ├── inject_gcp.go
│   ├── inject_local.go
│   └── wire_gen.go
├── aws
│   ├── aws.go
│   ├── blob.go
│   ├── runtimevar.go
│   └── user.go
├── cmd
│   ├── root.go
│   └── vfsgen.go
├── config
│   ├── cookie.go
│   └── utils.go
├── errors
│   ├── errors.go
│   └── stack.go
├── executor
│   ├── executor.go
│   ├── interface.go
│   ├── options.go
│   └── options_test.go
├── f.txt
├── files
│   ├── dir.go
│   ├── exec.go
│   ├── httpvfs.go
│   ├── path.go
│   ├── path_test.go
│   └── string.go
├── go.mod
├── go.sum
├── google
│   ├── app.go
│   ├── blob.go
│   ├── db.go
│   ├── gcloud.go
│   ├── kube.go
│   ├── run.go
│   ├── runtime_config.go
│   └── user.go
├── io
│   ├── close.go
│   └── io.go
├── logging
│   ├── logging.go
│   └── logging_test.go
├── logo
│   ├── dark_logo_transparent_background.png
│   ├── dark_logo_white_background.jpg
│   ├── logo_transparent_background.png
│   ├── logo_white_background.jpg
│   ├── white_logo_color_background.jpg
│   ├── white_logo_dark_background.jpg
│   └── white_logo_transparent_background.png
├── main.go
├── runtime
│   ├── config.go
│   ├── engine.go
│   ├── gateway.go
│   ├── grpc.go
│   ├── http_server_middleware.go
│   ├── mux.go
│   ├── options.go
│   ├── passing_header_middleware.go
│   ├── passing_header_middleware_test.go
│   ├── private
│   │   └── server.go
│   └── server.go
├── templates
│   ├── _data
│   │   ├── {{.ProtoDir}}
│   │   │   └── {{.Path}}.proto.tmpl
│   │   ├── {{.RootDir}}
│   │   │   ├── Makefile.tmpl
│   │   │   ├── main.go.tmpl
│   │   │   └── reviewdog.tmpl
│   │   ├── {{.ServerDir}}
│   │   │   ├── {{.Path}}_server.go.tmpl
│   │   │   ├── {{.Path}}_server_register_funcs.go.tmpl
│   │   │   └── {{.Path}}_server_test.go.tmpl
│   │   ├── {{.StaticDir}}
│   │   │   └── guestbook.html.tmpl
│   │   └── {{.TfDir}}
│   │       ├── main.tf.tmpl
│   │       ├── output.tf.tmpl
│   │       └── variable.tf.tmpl
│   ├── gen.go
│   └── vfsgen.go
├── tools.go
├── ui
│   └── ui.go
├── util
    └── exec.go
```
