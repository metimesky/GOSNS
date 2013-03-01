package home

import (
	. "app"
	. "model"
	"net/http"
)

func init() {
	R.HandleFunc("/", AppHandler(ac_home))
}

func ac_home(w http.ResponseWriter, r *http.Request) {

	var alist []Article
	Db.C("go_Article").Find(nil).Limit(10).All(&alist)
	cate := Group_allcate()
	IndexTemplate(w, r, map[string]interface{}{"alist": alist, "cate": cate}, "template/home/index.html")
}
