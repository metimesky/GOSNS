package ubb2html

import "regexp"
import "strings"

var emotpath string = "/static/js/xhEditor/xheditor_emot/"

func Ubb2html(ubb string) string {
	ubb = strings.Replace(ubb, "&", "&amp;", -1)
	ubb = strings.Replace(ubb, "<", "&lt;", -1)
	ubb = strings.Replace(ubb, ">", "&gt;", -1)
	ubb = strings.Replace(ubb, "\n", "<br />", -1)
	//加粗 斜体 下划线 上下标
	re, _ := regexp.Compile("\\[(\\/?)(b|u|s|sup|sub)\\]")
	ubb = re.ReplaceAllString(ubb, "<$1$2>")
	//颜色头部
	re, _ = regexp.Compile("\\[color\\s*=\\s*([^\\]\"]+?)(?:\"[^\\]]*?)?\\s*\\]")
	ubb = re.ReplaceAllString(ubb, `<span style="color:$1;">`)
	//大小头部
	re, _ = regexp.Compile(`\[size\s*=\s*([^\]"]+?)(?:"[^\]]*?)?\s*\]`)
	sizes := re.FindStringSubmatch(ubb)
	if len(sizes) > 0 {
		ubb = resizename(sizes, ubb)
	}
	//字体头部
	re, _ = regexp.Compile(`\[font\s*=\s*([^\]"]+?)(?:"[^\]]*?)?\s*\]`)
	ubb = re.ReplaceAllString(ubb, `<span style="font-family:$1;">`)
	//背景色头部
	re, _ = regexp.Compile(`\[back\s*=\s*([^\]"]+?)(?:"[^\]]*?)?\s*\]`)
	ubb = re.ReplaceAllString(ubb, `<span style="background-color:$1;">`)
	//颜色 大小 字体 背景色 尾部	
	re, _ = regexp.Compile("\\[\\/(color|size|font|back)\\]")
	ubb = re.ReplaceAllString(ubb, "</span>")
	//对齐
	for i := 0; i < 2; i++ {
		re, _ := regexp.Compile(`\[align\s*=\s*([^\]"]+?)(?:"[^\]]*?)?\s*\](.*)\[\/align\]`)
		ubb = re.ReplaceAllString(ubb, `<p align="$1">$2</p>`)
	}
	//图片	
	re, _ = regexp.Compile(`\[img\]([^\["]+)\[\/img\]`)
	ubb = re.ReplaceAllString(ubb, `<img src="$1" alt="" />`)

	//表情
	re, _ = regexp.Compile(`\[emot\s*=\s*([^\]"]+?)(?:"[^\]]*?)?\s*\/\]`)
	emot := re.FindStringSubmatch(ubb)
	if len(emot) > 1 {
		ubb = reemot(emot, ubb)
	}
	return ubb
}

func resizename(str []string, ubb string) string {
	fontsizes := []string{"10px", "12px", "13px", "16px", "18px", "24px", "32px", "48px"}
	b := false
	for _, v := range fontsizes {
		if v == str[1] {
			b = true
		}
	}
	if !b {
		str[1] = ""
	}
	ubb = strings.Replace(ubb, str[0], `<span style="font-size:'`+str[1]+`';">`, -1)
	return ubb
}

func reemot(str []string, ubb string) string {
	args := strings.Split(str[1], ",")
	path := emotpath + args[0] + "/" + args[1] + ".gif"
	ubb = strings.Replace(ubb, str[0], `<img src="`+path+`"  />`, -1)
	return ubb
}
