package model

import (
	. "lib/Util"
	"mgo/bson"
	"net/http"
	"time"
)

type User struct {
	Id_      bson.ObjectId `bson:"_id"`
	Username string
	Pwd      string
	Email    string
	Adminid  int       //管理员判断 0不是 大于0是 
	Face     string    //头像路径
	Regip    string    //注册IP
	Lastip   string    //最后访问Ip
	Lasttime time.Time //最后登录时间
}

func User_ishaveEmail(email string) bool {
	num, _ := Db.C("go_user").Find(bson.M{"email": email}).Count()
	if num == 0 {
		return false
	}
	return true
}

func User_ishaveName(name string) bool {
	num, _ := Db.C("go_user").Find(bson.M{"username": name}).Count()
	if num == 0 {
		return false
	}
	return true
}

func User_curuserid(w http.ResponseWriter, r *http.Request) string {
	sess := Gsession.SessionStart(w, r)
	userid := sess.Get("User")
	if userid == nil {
		return ""
	}
	return userid.(string)
}

func User_curuser(w http.ResponseWriter, r *http.Request) *User {
	id := User_curuserid(w, r)
	if id == "" {
		return nil
	}
	var user User
	err := Db.C("go_user").Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&user)
	CheckErr(err)
	return &user
}

func User_get(uid bson.ObjectId) *User {
	var u User
	err := Db.C("go_user").Find(bson.M{"_id": uid}).One(&u)
	CheckErr(err)
	return &u
}

//修改头像
func Face(userid string) {
	err := Db.C("go_user").Update(bson.M{"_id": bson.ObjectIdHex(userid)}, bson.M{"$set": bson.M{"face": userid}})
	CheckErr(err)
}
