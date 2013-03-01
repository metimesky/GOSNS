package group

import (
	. "app"
	"fmt"
	"html/template"
	. "lib/Util"
	"lib/ubb2html"
	"mgo/bson"
	. "model"
	"net/http"
	"time"
)

func init() {
	Gapps = append(Gapps, Apps{Pkgname: "group", Appname: "小组管理"})
	R.HandleFunc("/group/add", AppHandler(ac_grouppost))
	R.HandleFunc("/group/post", AppHandler(ac_commentpost))
}

//发帖
func ac_grouppost(w http.ResponseWriter, r *http.Request) {
	userid := User_curuserid(w, r)
	if userid == "" {
		ExitMsg(w, "请登录")
	}
	title := r.FormValue("title")
	if len(r.FormValue("title")) > 60 {
		title = r.FormValue("title")[0:60]
	}
	fmt.Println(len(r.FormValue("title")))
	g := &Group_posts{Id_: bson.NewObjectId(), Cateid: bson.ObjectIdHex(r.FormValue("cateid")),
		Groupid: bson.ObjectIdHex(r.FormValue("groupid")), Title: title,
		Userid: bson.ObjectIdHex(userid), Content: template.HTML(r.FormValue("text")),
		Lastuser: bson.ObjectIdHex(userid), Lasttime: time.Now()}
	err := Db.C("go_Group_posts").Insert(&g)
	CheckErr(err)
	http.Redirect(w, r, "/group/g/"+r.FormValue("groupid"), http.StatusFound)
}

//回复
func ac_commentpost(w http.ResponseWriter, r *http.Request) {
	userid := User_curuserid(w, r)
	if userid == "" {
		ExitMsg(w, "请登录")
	}
	text := ubb2html.Ubb2html(r.FormValue("text"))
	p := &Group_posts_comment{Id_: bson.NewObjectId(), Postid: bson.ObjectIdHex(r.FormValue("tid")),
		Userid: bson.ObjectIdHex(userid), Content: template.HTML(text),
		Referid: bson.ObjectIdHex(r.FormValue("rid"))}
	AddComment(p)
	http.Redirect(w, r, "/group/p/"+r.FormValue("tid"), http.StatusFound)
}

//添加回复
func AddComment(gps *Group_posts_comment) {
	err := Db.C("go_Group_posts_comment").Insert(&gps)
	CheckErr(err)
	e2 := Db.C("go_Group_posts").Update(bson.M{"_id": gps.Postid},
		bson.M{"$inc": bson.M{"count_comment": 1},
			"$set": bson.M{"lastuser": gps.Userid, "lasttime": time.Now()}})
	CheckErr(e2)
}
