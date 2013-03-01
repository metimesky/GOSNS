package group

import (
	. "app"
	//"fmt"
	. "lib/Util"
	"mgo/bson"
	. "model"
	"mux"
	"net/http"
)

func init() {
	R.HandleFunc("/group/p/{tid}", AppHandler(ac_tshow))
}

func ac_tshow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var post Group_posts
	err := Db.C("go_Group_posts").Find(bson.M{"_id": bson.ObjectIdHex(vars["tid"])}).One(&post)
	CheckErr(err)
	var comm []Group_posts_comment
	Db.C("go_Group_posts_comment").Find(bson.M{"postid": post.Id_}).All(&comm)
	Db.C("go_Group_posts").Update(bson.M{"_id": post.Id_}, bson.M{"$inc": bson.M{"count_click": 1}})
	IndexTemplate(w, r, map[string]interface{}{"post": post, "comm": comm}, "template/group/tlist.html")
}
