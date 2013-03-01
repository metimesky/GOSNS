package user

import (
	. "app"
	//"image"
	//"fmt"
	"image/jpeg"
	"io/ioutil"
	"lib/resize"
	"model"
	"net/http"
	"os"
)

func init() {
	R.HandleFunc("/user/uptx", AppHandler(ac_uploadtx))
}

func ac_uploadtx(w http.ResponseWriter, r *http.Request) {
	imgdata, _ := ioutil.ReadAll(r.Body)
	userid := model.User_curuserid(w, r)
	imgname := "static\\touxiang\\" + userid + ".jpg"
	f, _ := os.Create(imgname)
	f.Write(imgdata)
	f.Close()

	file, _ := os.Open(imgname)
	img, _ := jpeg.Decode(file)
	file.Close()
	rm := resize.Resize(100, 100, img, resize.Lanczos3)
	out, _ := os.Create(imgname)
	defer out.Close()
	jpeg.Encode(out, rm, nil)
	model.Face(userid)
}
