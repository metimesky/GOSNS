package tag

import (
	. "app"
	"fmt"
	. "lib/Util"
	//_ "mgo"
	"mgo/bson"
	. "model"
	"net/http"
)

func init() {
	Gapps = append(Gapps, Apps{Pkgname: "tag", Appname: "标签管理"})
	R.HandleFunc("/admin/tag/index", AppHandler(ac_taglist))
	R.HandleFunc("/admin/tag/list", AppHandler(ac_taglist))
	R.HandleFunc("/admin/tag/action", AppHandler(ac_tagaction))
}

//标签浏览
func ac_taglist(w http.ResponseWriter, r *http.Request) {
	var tags []Tag
	Db.C("go_tag").Find(nil).All(&tags)
	AdminTemplate(w, r, map[string]interface{}{"tags": tags}, "template/tag/admin_taglist.html")
}

//标签删除修改
func ac_tagaction(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("op") {
	case "update":
		err := Db.C("go_tag").Update(bson.M{"_id": bson.ObjectIdHex(r.FormValue("id"))},
			bson.M{"$set": bson.M{"tagname": r.FormValue("tagname")}})
		CheckErr(err)
		fmt.Fprint(w, "ok")
	case "del":

	}
}
