package app

import (
	"fmt"
	"html/template"
	. "lib/Util"
	"mgo/bson"
	"model"
	"mux"
	"net/http"
	"time"
)

var R *mux.Router
var Gapps []Apps

func init() {
	R = mux.NewRouter()
	fmt.Println("load app")
}

type Apps struct {
	Pkgname string
	Appname string
}

//错误句柄
func AppHandler(fn http.HandlerFunc, role ...int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(role) > 0 {
			if role[0] == 1 {
				id := model.User_curuserid(w, r)
				if id == "" {
					ExitMsg(w, "请先登录", "/user/login")
					return
				}
			}
		}

		defer func() {
			if err, ok := recover().(error); ok {
				fmt.Fprint(w, err)
			}
		}()
		fn(w, r)
	}
}

func AdminTemplate(w http.ResponseWriter, r *http.Request, data map[string]interface{}, file string) {
	t, err := template.ParseFiles("template/adminbase.html", file)
	CheckErr(err)
	data["apps"] = Gapps
	data["view"] = &View{}
	data["user"] = model.User_curuser(w, r)
	t.Execute(w, data)
}

func IndexTemplate(w http.ResponseWriter, r *http.Request, data map[string]interface{}, file string) {
	t, err := template.ParseFiles("template/indexbase.html", file)
	CheckErr(err)
	data["view"] = &View{}
	data["user"] = model.User_curuser(w, r)
	t.Execute(w, data)
}

type View struct {
}

func (g *View) FormatTime(bid bson.ObjectId) string {
	return bid.Time().Format("2006-01-02 15:04")
}

func (v *View) FormatTime_t(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)
	if duration.Seconds() < 60 {
		return fmt.Sprintf("%.0f 秒前", duration.Seconds())
	} else if duration.Minutes() < 60 {
		return fmt.Sprintf("%.0f 分钟前", duration.Minutes())
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%.0f 小时前", duration.Hours())
	}

	return t.Format("2006-01-02 15:04")
}

func (v *View) FormatTime_2(bid bson.ObjectId) string {
	now := time.Now()
	t := bid.Time()
	duration := now.Sub(t)
	if duration.Hours() < 24 {
		return t.Format("15:04")
	}
	return t.Format("01-02")
}

func (v *View) To(i interface{}) *model.Group_posts {
	var a model.Group_posts
	a = i.(model.Group_posts)
	return &a
}

func (v *View) Num(i int, add int) int {
	return i + add
}
