package model

import (
	"html/template"
	. "lib/Util"
	"mgo/bson"
)

//文章分类
type Article_cate struct {
	Id_      bson.ObjectId `bson:"_id"`
	Catename string
	Orderid  int
}

//文章
type Article struct {
	Id_           bson.ObjectId `bson:"_id"`
	Cateid        bson.ObjectId
	Title         string
	Content       template.HTML
	Count_comment int             //评论数
	Tags          []bson.ObjectId //标签
}

//文章评论
type Article_comment struct {
	Id_       bson.ObjectId `bson:"_id"`
	Articleid bson.ObjectId
	Userid    bson.ObjectId
	Content   template.HTML
}

//返回所有文章分类 sel 为过滤字段
func Article_allcate(sel bson.M) []Article_cate {
	var cates []Article_cate
	q := Db.C("go_Article_cate").Find(nil)
	if sel != nil {
		q.Select(sel)
	}
	q.Sort("orderid").All(&cates)
	return cates
}

//返回所有文章分类 id为键 分类结构体为value的map
func Article_mapcate(sel bson.M) map[bson.ObjectId]Article_cate {
	cates := Article_allcate(sel)
	mapcates := make(map[bson.ObjectId]Article_cate)
	for _, v := range cates {
		mapcates[v.Id_] = v
	}
	return mapcates
}

//返回单个文章
func Article_Get(bid bson.ObjectId) Article {
	var ar Article
	Db.C("go_Article").Find(bson.M{"_id": bid}).One(&ar)
	return ar
}

func Article_addcomment() {

}

//返回某个文章的分类
func Article_Getcomment(articleid bson.ObjectId) []Article_comment {
	var ac []Article_comment
	err := Db.C("go_Article_comment").Find(bson.M{"articleid": articleid}).All(&ac)
	CheckErr(err)
	return ac
}

func (ac *Article_comment) User() *User {
	return User_get(ac.Userid)
}
