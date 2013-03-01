package group

import (
	. "app"
	. "lib/Util"
	"mgo/bson"
	. "model"
	"net/http"
)

func init() {
	R.HandleFunc("/admin/group/group", AppHandler(ac_admin_group, 1))
}

func ac_admin_group(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		switch r.FormValue("op") {
		case "add":
			userid := User_curuserid(w, r)
			cateid := r.FormValue("cateid")
			group := &Group{Id_: bson.NewObjectId(), Userid: bson.ObjectIdHex(userid),
				Cateid: bson.ObjectIdHex(cateid), Count_user: 0, Count_topic: 0,
				Groupname: r.FormValue("groupname"), Role_user: r.FormValue("role_user")}
			err := Db.C("go_Group").Insert(&group)
			Group_cateaddgroup(bson.ObjectIdHex(cateid))
			CheckErr(err)
			ExitMsg(w, "添加成功")
		}
	}
	switch r.FormValue("op") {
	case "":
		list := Group_list()
		AdminTemplate(w, r, map[string]interface{}{"list": list}, "template/group/admin_group.html")
	case "showadd":
		cates := Group_allcate()
		AdminTemplate(w, r, map[string]interface{}{"cates": cates}, "template/group/admin_groupadd.html")
	}
}
