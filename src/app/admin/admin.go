package admin

import (
	. "app"
	//"html/template"
	"net/http"
)

func init() {
	R.HandleFunc("/admin", AppHandler(ac_admin))
}

func ac_admin(w http.ResponseWriter, r *http.Request) {
	AdminTemplate(w, r, map[string]interface{}{}, "template/admin/index.html")
}
