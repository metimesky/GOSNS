package model

import (
	. "lib/Util"
	"lib/session"
	_ "lib/session/memory"
	"mgo"
)

var Db *mgo.Database
var Gsession *session.Manager

func init() {
	sess, err := mgo.Dial("localhost:27017")
	CheckErr(err)
	Db = sess.DB("godb")
	Gsession, _ = session.NewManager("memory", "gosessionid", 3600)
	go Gsession.GC()
}
