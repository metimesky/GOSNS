package user

import (
	. "app"
	"fmt"
	"html/template"
	. "lib/Util"
	"lib/captcha"
	"lib/forms"
	"mgo/bson"
	. "model"
	"net/http"
	"time"
)

func init() {
	R.HandleFunc("/user/reg", AppHandler(ac_user))
	R.HandleFunc("/user/login", AppHandler(ac_login))
	R.HandleFunc("/user/center", AppHandler(ac_usercenter))
	fmt.Println("load user")
	Gapps = append(Gapps, Apps{Pkgname: "user", Appname: "会员管理"})

}
func ac_user(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		user := r.FormValue("user")
		pwd := r.FormValue("pwd")
		if !forms.Reg_email(email) {
			ExitMsg(w, "邮箱格式错误")
			return
		}

		if !forms.Reg_user(user) {
			ExitMsg(w, "用户名格式错误")
			return
		}
		if pwd == "" {
			ExitMsg(w, "密码不可为空")
			return
		}
		pwd = forms.Tomd5(pwd)

		if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
			ExitMsg(w, "验证码错误")
			return
		}
		if User_ishaveEmail(email) {
			ExitMsg(w, "Email 已存在")
			return
		}
		if User_ishaveName(user) {
			ExitMsg(w, "用户名 已存在")
			return
		}
		C := Db.C("go_user")
		id := bson.NewObjectId()
		err := C.Insert(&User{Id_: id, Username: user, Pwd: pwd, Email: email, Adminid: 0,
			Regip: forms.Ip(r), Lastip: forms.Ip(r), Lasttime: time.Now()})
		CheckErr(err)
		fmt.Fprint(w, "添加成功")
		return
	}

	t, err := template.ParseFiles("template/user/index.html")
	CheckErr(err)
	t.Execute(w, struct{ CaptchaId string }{captcha.New()})
}

func ac_login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		pwd := r.FormValue("pwd")
		loginuser := User{}
		err := Db.C("go_user").Find(bson.M{"email": email}).One(&loginuser)
		if err != nil {
			ExitMsg(w, "用户名不存在")
			return
		}
		if loginuser.Pwd != forms.Tomd5(pwd) {
			ExitMsg(w, "密码错误")
			return
		}
		Login_info(loginuser.Id_, r)
		sess := Gsession.SessionStart(w, r)
		sess.Set("User", loginuser.Id_.Hex())
		http.Redirect(w, r, "/user/center", http.StatusFound)
		return
	}
	IndexTemplate(w, r, map[string]interface{}{}, "template/user/login.html")
}

func ac_usercenter(w http.ResponseWriter, r *http.Request) {
	sess := Gsession.SessionStart(w, r)
	suid := sess.Get("User")
	if suid == nil {
		ExitMsg(w, "未登录", "/user/login")
		return
	}

	switch r.FormValue("t") {
	case "":
		Userinfo := User{}
		err := Db.C("go_user").Find(bson.M{"_id": bson.ObjectIdHex(suid.(string))}).One(&Userinfo)
		CheckErr(err)
		IndexTemplate(w, r, map[string]interface{}{"info": Userinfo}, "template/user/usercenter.html")
	case "tx":
		IndexTemplate(w, r, map[string]interface{}{}, "template/user/touxiang.html")
	}

}

func Login_info(id bson.ObjectId, r *http.Request) {
	err := Db.C("go_user").Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"lastip": forms.Ip(r),
		"lasttime": time.Now()}})
	CheckErr(err)
}
