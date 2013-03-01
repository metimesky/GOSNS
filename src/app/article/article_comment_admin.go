package article

import (
	. "app"
	. "lib/Util"
	"mgo/bson"
	. "model"
	"net/http"
)

func init() {
	R.HandleFunc("/admin/article/comment", AppHandler(ac_admin_comment))
}

func ac_admin_comment(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("op") {
	case "":
		var comments []Article_comment
		err := Db.C("go_Article_comment").Find(nil).All(&comments)
		CheckErr(err)
		AdminTemplate(w, r, map[string]interface{}{"comments": comments}, "template/article/admin_comment.html")
	case "del":
		id := bson.ObjectIdHex(r.FormValue("id"))
		err := Db.C("go_Article_comment").Remove(bson.M{"_id": id})
		CheckErr(err)
		http.Redirect(w, r, "/admin/article/comment", http.StatusFound)
	}
}
