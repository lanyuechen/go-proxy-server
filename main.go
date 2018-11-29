package main

import (
	"flag"
	"github.com/go-martini/martini"
	"log"
	"net/http"
)

func main() {
	log.Println("start", martini.Env)

	addr := flag.String("p", ":3000", "address where the server listen on")
	war := flag.String("d", "./", "directory of web files")
	proxyConfig := flag.String("proxy", "./proxy.config.json", "proxy config directory")
	flag.Parse()

	mount(*war, *proxyConfig)

	log.Printf("start server on %s", *addr)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func mount(war string, proxyConfig string) {
	m := martini.Classic()
	m.Use(martini.Static(war, martini.StaticOptions{SkipLogging: true}))

	//请求注入web参数
	m.Use(func(w http.ResponseWriter, c martini.Context) {
		web := &Web{w: w}
		c.Map(web)
	})

	m.Use(midProxy(proxyConfig))

	/*------------------------ API --------------------------*/
	m.Group("", func(r martini.Router) {
		r.Get("/test", testHandler)
	})

	http.Handle("/", m)
}

func testHandler(web *Web) (int, string) {
	return web.Json(200, J{"test": "success"})
}
