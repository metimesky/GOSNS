package article

import (
	. "app"
	"fmt"
	"html/template"
	. "lib/Util"
	"lib/forms"
	"mgo/bson"
	. "model"
	"net/http"
)

func init() {
	R.HandleFunc("/admin/article/index", AppHandler(ac_admin_catelist)) //管理首页
	R.HandleFunc("/admin/article/catelist", AppHandler(ac_admin_catelist))
	R.HandleFunc("/admin/article/add", AppHandler(ac_admin_article))
	R.HandleFunc("/admin/article/list", AppHandler(ac_admin_articlelist))
	fmt.Println("load article")
}

//文章分类 增删改
func ac_admin_catelist(w http.ResponseWriter, r *http.Request) {
	C := Db.C("go_Article_cate")

	switch r.FormValue("op") {
	case "add": //添加
		catename := r.FormValue("catename")
		orderid := r.FormValue("sortid")
		if forms.IsEmpty(w, catename, orderid) {
			return
		}
		err := C.Insert(&Article_cate{Id_: bson.NewObjectId(),
			Catename: catename, Orderid: forms.Toint(orderid)})
		CheckErr(err)
		http.Redirect(w, r, "/admin/article/catelist", http.StatusFound)
		return
	case "del": //删除
		err := C.Remove(bson.M{"_id": bson.ObjectIdHex(r.FormValue("delid"))})
		CheckErr(err)
		http.Redirect(w, r, "/admin/article/catelist", http.StatusFound)
	case "update": //修改
		if r.FormValue("catename") == "" || r.FormValue("sortid") == "" {
			fmt.Fprint(w, "no")
			return
		}
		err := C.Update(bson.M{"_id": bson.ObjectIdHex(r.FormValue("id"))},
			bson.M{"$set": bson.M{"catename": r.FormValue("catename"),
				"orderid": forms.Toint(r.FormValue("sortid"))}})
		CheckErr(err)
		fmt.Fprint(w, "ok")
		return
	case "": //显示列表
		cates := Article_allcate(nil)
		AdminTemplate(w, r, map[string]interface{}{"cates": cates}, "template/article/admin_catelist.html")
	}
}

//文章添加编辑
func ac_admin_article(w http.ResponseWriter, r *http.Request) {
	C := Db.C("go_Article")
	if r.Method == "POST" {
		if forms.IsEmpty(w, r.FormValue("title"), r.FormValue("cateid")) {
			return
		}
		switch r.FormValue("op") {
		case "add":
			inid := bson.NewObjectId()
			err := C.Insert(&Article{Id_: inid, Cateid: bson.ObjectIdHex(r.FormValue("cateid")),
				Title: r.FormValue("title"), Content: template.HTML(r.FormValue("elm1")), Count_comment: 0})
			CheckErr(err)
			var nullstr []bson.ObjectId
			Tag_add(r.FormValue("tags"), "Article", inid, nullstr)
			ExitMsg(w, "添加成功")
			return
		case "update":
			err := C.Update(bson.M{"_id": bson.ObjectIdHex(r.FormValue("id"))}, bson.M{"$set": bson.M{
				"cateid": bson.ObjectIdHex(r.FormValue("cateid")),
				"title":  r.FormValue("title"), "content": template.HTML(r.FormValue("elm1"))}})
			var getarticle Article
			CheckErr(err)
			gerr := C.Find(bson.M{"_id": bson.ObjectIdHex(r.FormValue("id"))}).One(&getarticle)
			CheckErr(gerr)
			Tag_add(r.FormValue("tags"), "Article", bson.ObjectIdHex(r.FormValue("id")), getarticle.Tags)
			ExitMsg(w, "修改成功")
			return
		}
	}

	switch r.FormValue("op") {
	case "":
		cates := Article_allcate(nil)
		AdminTemplate(w, r, map[string]interface{}{"cates": cates}, "template/article/admin_addarticle.html")
	case "edit":
		cates := Article_allcate(nil)
		article := Article_Get(bson.ObjectIdHex(r.FormValue("id")))
		tags := TagidToname(article.Tags)
		AdminTemplate(w, r, map[string]interface{}{"cates": cates, "article": article, "tags": tags},
			"template/article/admin_editarticle.html")
	}
}

//文章浏览删除
func ac_admin_articlelist(w http.ResponseWriter, r *http.Request) {
	C := Db.C("go_Article")
	switch r.FormValue("op") {
	case "":
		var articles []Article
		C.Find(nil).Select(bson.M{"content": false}).All(&articles)
		cates := Article_mapcate(bson.M{"orderid": 0})
		AdminTemplate(w, r, map[string]interface{}{"articles": articles,
			"cates": cates}, "template/article/admin_articlelist.html")
	case "del":
		err := C.Remove(bson.M{"_id": bson.ObjectIdHex(r.FormValue("delid"))})
		CheckErr(err)
		http.Redirect(w, r, "/admin/article/list", http.StatusFound)
	}
}
