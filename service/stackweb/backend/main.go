package main

import (
	"net/http"

	"github.com/stack-labs/stack-rpc/web"
)

func main() {
	s := web.NewService(
		web.Name("stack.stackweb"),
	)

	// favicon.ico
	s.HandleFunc("/favicon.ico", faviconHandler)
	// static dir
	s.Handle(rootPath+"/", http.StripPrefix(rootPath+"/", http.FileServer(http.Dir(StaticDir))))

	if err := s.Init(
		web.Action(
			func(c *cli.Context) {
				// do something
			}),
	); err != nil {
		panic(err)
	}

	if err := s.Run(); err != nil {
		panic(err)
	}
}
