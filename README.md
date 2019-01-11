![](https://github.com/gofunct/common/blob/master/logo/dark_logo_transparent_background.png?raw=true)  

# Common

**This is a common library I use for other projects. Some of the libraries may be moved to their own
repositories in the future.**

* Author: Coleman Word 
* Email: coleman.word@gofunct.com
* Download: `go get github.com/gofunct/common/...`

## Overview
* Root directory holds initialization functions that act as providers for the main types in other packages
* Wire is used for dependency injection


## Key Types

    type Runtime interface {
    exec.Interface
    io.ReaderWriter
    }

    type HandleFunc func(context.Context, interface{})   (interface{}, error) 

is the function that is required to initialize the runtime. This handler is then passed to the runtime's wrapper functions:

    type ChainFunc func(Handler) Handler

This Morphs the Handler in a chain until it reaches the required Execute() error method which finally executes the Handler function.

This setup allows you to create programs that support plugable components with an interface that is easily satisfied by types like cobra.Command.Execute, os.Exec.Run, http.Server.Serve, etc...

    type Service struct {
        config.Service
        iio.Service
        exec.Service
        Handle.   HandleFunc
        Chain.    ChainFunc
    }


The service type implements the Runtime interface and additionally provides functionality for reading runtime configuration from remote endpoint, or from a local file. 