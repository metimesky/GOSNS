package model

import (
	"fmt"
	"html/template"
	. "lib/Util"
	"lib/forms"
	"mgo/bson"
	"strings"
)

type Dotag interface {
	Objaddtag(Tag, bson.ObjectId)
}

type Tag struct {
	Id_           bson.ObjectId `bson:"_id"`
	Tagname       string
	Count_Article int
	Count         int
}

func Tag_add(strtags, tagtype string, objid bson.ObjectId, oldids []bson.ObjectId) {
	if strtags == "" || strtags == "," || len(strtags) > 100 {
		return
	}
	strtags = template.HTMLEscapeString(strtags)
	strtags = strings.Replace(strtags, "，", ",", -1)
	tags := strings.Split(strtags, ",")
	tags = forms.Unique(tags)
	addtag(tags, tagtype, objid, oldids)
}

func addtag(subtag []string, tagtype string, objid bson.ObjectId, oldids []bson.ObjectId) {
	oldtag := TagidToname(oldids)
	newtag := silcenothave(subtag, oldtag)
	removetag := silcenothave(oldtag, subtag)

	///处理新增的标签
	if len(newtag) > 0 {
		for _, ntv := range newtag { //遍历新增的标签
			tag := Ishavetag(ntv)
			if tag == nil { //如果数据库中没有
				newid := bson.NewObjectId()
				CheckErr(Db.C("go_tag").Insert(&Tag{Id_: newid, Tagname: ntv, Count: 1}))
				CheckErr(Db.C("go_"+tagtype).Update(bson.M{"_id": objid}, bson.M{"$push": bson.M{"tags": newid}}))
				Tagnum(newid, tagtype)
			} else { //如果数据库中已有
				isinnum, err := Db.C("go_" + tagtype).Find(bson.M{"_id": objid, "tags": tag.Id_}).Count()
				CheckErr(err)
				if isinnum == 0 {
					CheckErr(Db.C("go_"+tagtype).Update(bson.M{"_id": objid}, bson.M{"$push": bson.M{"tags": tag.Id_}}))
					Tagnum(tag.Id_, tagtype)
					Tagcount(tag.Id_, 1)
				}
			}
		}
	}
	fmt.Println("newtag")
	fmt.Println(newtag)
	fmt.Println("removetag")
	fmt.Println(removetag)
	//处理删除的标签
	if len(removetag) > 0 {
		for _, delv := range removetag {
			b, tid := Getagid(delv)
			if !b {
				panic("并不存在需要删除的标签 model-tag_model.go 58")
			}
			Db.C("go_"+tagtype).Update(bson.M{"_id": objid}, bson.M{"$pull": bson.M{"tags": tid}})
			Tagnum(tid, tagtype)
			Tagcount(tid, -1)
		}
	}

}

func Ishavetag(tname string) *Tag {
	var tag Tag
	err := Db.C("go_tag").Find(bson.M{"tagname": tname}).One(&tag)
	if err != nil {
		return nil
	}
	return &tag
}

//根据标签id返回标签名字符串数组
func TagidToname(ids []bson.ObjectId) (strtags []string) {
	var tags []Tag
	err := Db.C("go_tag").Find(bson.M{"_id": bson.M{"$in": ids}}).All(&tags)
	CheckErr(err)
	for _, v := range tags {
		strtags = append(strtags, v.Tagname)
	}
	return
}

//返回s1中不存在于s2的字符串切片
func silcenothave(s1, s2 []string) (s []string) {
	for _, v := range s1 {
		if !sliceishavestr(s2, v) {
			s = append(s, v)
		}
	}
	return
}

//判断一个字符串切片(astr)中是否包含某(str)字符串
func sliceishavestr(astr []string, str string) bool {
	for _, v := range astr {
		if v == str {
			return true
		}
	}

	return false
}

func Tagnum(tid bson.ObjectId, tagtype string) {
	num, err := Db.C("go_" + tagtype).Find(bson.M{"tags": tid}).Count()
	if err != nil {
		num = 0
	}

	ierr := Db.C("go_tag").Update(bson.M{"_id": tid},
		bson.M{"$set": bson.M{"count_" + strings.ToLower(tagtype): num}})
	CheckErr(ierr)
}

func Getagid(tname string) (bool, bson.ObjectId) {
	var t Tag
	err := Db.C("go_tag").Find(bson.M{"tagname": tname}).One(&t)
	if err != nil {
		return false, bson.ObjectId(0)
	}
	return true, t.Id_
}

func Tagcount(id bson.ObjectId, num int) {
	var tag Tag
	Db.C("go_tag").Find(bson.M{"_id": id}).One(&tag)
	newnum := tag.Count + num
	if newnum < 1 {
		Db.C("go_tag").Remove(bson.M{"_id": id})
		return
	}
	Db.C("go_tag").Update(bson.M{"_id": id}, bson.M{"$inc": bson.M{"count": num}})
}
