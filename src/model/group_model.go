package model

import (
	"html/template"
	. "lib/Util"
	"mgo/bson"
	"time"
)

//小组分类
type Group_cate struct {
	Id_         bson.ObjectId `bson:"_id"`
	Catename    string
	Count_group int
}

//小组
type Group struct {
	Id_         bson.ObjectId `bson:"_id"`
	Groupname   string        //组名
	Userid      bson.ObjectId //组长ID
	Cateid      bson.ObjectId //分类ID
	Count_user  int           //会员数
	Count_topic int           //帖子数
	Role_user   string        //会员名称
}

//会员加入了哪个组
type Group_collects struct {
	Id_     bson.ObjectId `bson:"_id"`
	Userid  bson.ObjectId //用户id
	Groupid bson.ObjectId //小组ID
}

//帖子列表
type Group_posts struct {
	Id_           bson.ObjectId `bson:"_id"`
	Cateid        bson.ObjectId //分类ID
	Groupid       bson.ObjectId //小组ID
	Userid        bson.ObjectId //用户ID
	Title         string        //标题
	Content       template.HTML //正文
	Isposts       bool          //是否为精华帖
	Istop         bool          //是否置顶
	Count_comment int           //回复统计
	Count_click   int           //点击数
	Lastuser      bson.ObjectId //最后回复会员
	Lasttime      time.Time     //最后回复时间
}

//帖子回复列表
type Group_posts_comment struct {
	Id_     bson.ObjectId `bson:"_id"`
	Postid  bson.ObjectId //帖子ID
	Referid bson.ObjectId //回复楼ID
	Userid  bson.ObjectId //用户ID
	Content template.HTML //回复内容
}

//所有分类
func Group_allcate() []Group_cate {
	var groupcates []Group_cate
	err := Db.C("go_Group_cate").Find(nil).All(&groupcates)
	CheckErr(err)
	return groupcates
}

//返回所有小组 （待完善）
func Group_list() []Group {
	var g []Group
	err := Db.C("go_Group").Find(nil).All(&g)
	CheckErr(err)
	return g
}

//添加小组 分类所属小组加一
func Group_cateaddgroup(id bson.ObjectId) {
	err := Db.C("go_Group_cate").Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"count_group": 1}})
	CheckErr(err)
}

func Group_get(id bson.ObjectId) *Group {
	var g Group
	err := Db.C("go_Group").Find(bson.M{"_id": id}).One(&g)
	CheckErr(err)
	return &g
}

//帖子获取楼主信息
func (gp *Group_posts) User() *User {
	return User_get(gp.Userid)
}

//帖子获取小组信息
func (gp *Group_posts) Group() *Group {

	return Group_get(gp.Groupid)
}

//小组分类获取所含小组
func (gc *Group_cate) Groups() (r []Group) {
	Db.C("go_Group").Find(bson.M{"cateid": gc.Id_}).All(&r)
	return
}

//评论获取会员信息
func (gpc *Group_posts_comment) User() *User {
	return User_get(gpc.Userid)
}

//获取最后回复用户信息
func (g *Group_posts) LastUser() *User {
	return User_get(g.Lastuser)
}
