package main

import (
	"app"
	_ "app/admin"
	_ "app/article"
	_ "app/group"
	_ "app/home"
	_ "app/tag"
	_ "app/user"

	"fmt"
	"lib/captcha"
	"net/http"
)

func main() {
	fmt.Println("load func main ")
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.Handle("/", app.R)
	http.ListenAndServe(":8080", nil)

}
