package article

import (
	. "app"
	"html/template"
	. "lib/Util"
	"lib/ubb2html"
	"mgo/bson"
	. "model"
	"mux"
	"net/http"
)

func init() {
	Gapps = append(Gapps, Apps{Pkgname: "article", Appname: "文章管理"})
	R.HandleFunc("/article/list", AppHandler(ac_articlelist))
	R.HandleFunc("/article/show/{id}.html", AppHandler(ac_articleshow))
	R.HandleFunc("/article/comment", AppHandler(ac_articlecomment))
}

func ac_articlelist(w http.ResponseWriter, r *http.Request) {
	var alist []Article
	Db.C("go_Article").Find(nil).All(&alist)
	IndexTemplate(w, r, map[string]interface{}{"alist": alist}, "template/article/articlelist.html")
}

func ac_articleshow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article := Article_Get(bson.ObjectIdHex(vars["id"]))
	tags := TagidToname(article.Tags)
	comment := Article_Getcomment(bson.ObjectIdHex(vars["id"]))
	IndexTemplate(w, r,
		map[string]interface{}{"article": article, "tags": tags, "comment": comment},
		"template/article/articleshow.html")

}

func ac_articlecomment(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := User_curuserid(w, r)
		if id == "" {
			ExitMsg(w, "请先登录")
			return
		}
		ac := Article_comment{Id_: bson.NewObjectId(),
			Articleid: bson.ObjectIdHex(r.FormValue("articleid")),
			Userid:    bson.ObjectIdHex(id),
			Content:   template.HTML(ubb2html.Ubb2html(r.FormValue("text"))),
		}
		err := Db.C("go_Article_comment").Insert(&ac)
		CheckErr(err)
		ExitMsg(w, "评论成功")
	}
}
