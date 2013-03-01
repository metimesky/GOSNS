package group

import (
	. "app"
	"fmt"
	. "lib/Util"
	"mgo/bson"
	. "model"
	"net/http"
)

func init() {
	R.HandleFunc("/admin/group/index", AppHandler(ac_admin_groupcate, 1))
	R.HandleFunc("/admin/group/cate", AppHandler(ac_admin_groupcate, 1))
}

func ac_admin_groupcate(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("op") {
	case "":
		var cates []Group_cate
		err := Db.C("go_Group_cate").Find(nil).All(&cates)
		CheckErr(err)
		AdminTemplate(w, r, map[string]interface{}{"cates": cates}, "template/group/admin_groupcate.html")
	case "add":
		err := Db.C("go_Group_cate").Insert(&Group_cate{Id_: bson.NewObjectId(),
			Catename:    r.FormValue("catename"),
			Count_group: 0})
		CheckErr(err)
		http.Redirect(w, r, "/admin/group/cate", http.StatusFound)
	case "del":
		err := Db.C("go_Group_cate").Remove(bson.M{"_id": bson.ObjectIdHex(r.FormValue("id"))})
		CheckErr(err)
		http.Redirect(w, r, "/admin/group/cate", http.StatusFound)
	case "update":
		id := bson.ObjectIdHex(r.FormValue("id"))
		//fmt.Println(id)
		err := Db.C("go_Group_cate").Update(bson.M{"_id": id},
			bson.M{"$set": bson.M{"catename": r.FormValue("catename")}})
		//fmt.Println(id)
		CheckErr(err)
		fmt.Fprint(w, "ok")
	}
}
