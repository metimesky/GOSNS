package user

import (
	. "app"
	. "model"
	"net/http"
)

func init() {
	R.HandleFunc("/admin/user/index", AppHandler(ac_userlist))
}

func ac_userlist(w http.ResponseWriter, r *http.Request) {
	var ulist []User
	Db.C("go_user").Find(nil).All(&ulist)
	AdminTemplate(w, r, map[string]interface{}{"ulist": ulist}, "template/user/admin_user.html")
}
