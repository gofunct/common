![](https://github.com/gofunct/common/blob/master/logo/dark_logo_transparent_background.png?raw=true)  

# Common

**This is a common library I use for other projects. Some of the libraries may be moved to their own
repositories in the future.**

* Author: Coleman Word 
* Email: coleman.word@gofunct.com
* Download: `go get github.com/gofunct/common/...`

## Rules 

1. Every package should have an interface named Interface
1. Every package should have an interface implementation called API
1. Root package interfaces have **capital letters**
1. Api implementations are named "API"
1. If you cant make an interface work for a certain type, use viper to read the necessary info from config

## Methodology
1. create package for interface type ---> ex: exec/exec.go
2. create interface definition named Interface ---> exec.Interface
3. create package interface implementation named Api ---> exec.Api
4. create root file for package import ---> exec.go
5. create interface named {package name | upper}} that consumes the previous package interface(Interface) ---> exec.go(Interface)
6. create an initializer for the new interface named New{package name | upper}} ---> exec.go


