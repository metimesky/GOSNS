package group

import (
	. "app"
	. "lib/Util"
	"mgo/bson"
	. "model"
	"mux"
	"net/http"
)

func init() {
	R.HandleFunc("/group/list", AppHandler(ac_grouplist))
	R.HandleFunc("/group/g/{group}", AppHandler(ac_group))
}

func ac_grouplist(w http.ResponseWriter, r *http.Request) {
	cates := Group_allcate()
	IndexTemplate(w, r, map[string]interface{}{"cates": cates}, "template/group/grouplist.html")
}

func ac_group(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var group Group
	err := Db.C("go_Group").Find(bson.M{"_id": bson.ObjectIdHex(vars["group"])}).One(&group)
	CheckErr(err)
	var posts []Group_posts
	Db.C("go_Group_posts").Find(bson.M{"groupid": bson.ObjectIdHex(vars["group"])}).
		Sort("-lasttime").All(&posts)
	IndexTemplate(w, r, map[string]interface{}{"group": group, "posts": posts,
		"user": User_curuser(w, r)}, "template/group/group.html")
}
