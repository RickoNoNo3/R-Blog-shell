package processor

import (
	"encoding/json"
	"fmt"

	"../myerror"
)

type ReadProcessor struct{}

var LastReadDir = "0"

// 把BlogCore发回来的Json数据处理成文字表格. 再拼接一个当前位置
func (p ReadProcessor) Process(input string) (res string, err error) {
	// 判断是否是read命令
	args := splitArgs(input)
	if args[0] != "read" {
		err = myerror.CannotProcessError
		return
	}
	// 判断是否是回声用法
	if len(args) == 1 {
		args = append(args, "1", LastReadDir)
	} else if len(args) == 2 {
		args = []string{
			"read",
			"1",
			"0",
		}
	}
	// 连接BlogCore
	if res, err = runBlogCore(args); err != nil {
		return
	}
	// 处理json数据
	jsonStr := res
	res = ""
	if args[1] == "0" {
		var article _article
		if err = json.Unmarshal([]byte(jsonStr), &article); err != nil {
			return
		}
		res += fmt.Sprintf("Title: %v\nCreatedT: %v\nModifiedT: %v\n", article.Title, article.CreatedT, article.ModifiedT)
		if len([]rune(article.Html)) >= 10000 {
			article.Html = string([]rune(article.Html)[:10000]) + "\n[and more...]"
		}
		res += fmt.Sprintln("Content:\n" + article.Html)
	} else {
		var dir _dir
		if err = json.Unmarshal([]byte(jsonStr), &dir); err != nil {
			return
		}
		fmtStr := "| %4v | %4v | %v | %19v | %19v |\n"
		res += fmt.Sprintf(fmtStr, "Type", "Id", "             Text             ", "CreatedT", "ModifiedT")
		if len(dir.List) == 0 {
			res += "       该目录暂时没有内容\n"
		}
		for _, dirContent := range dir.List {
			// 处理中文的特殊对其格式
			maxChars, realRunes, englishChars := 0, 0, 0
			for _, ch := range dirContent.Text {
				if maxChars >= 29 {
					break
				}
				if ch <= 0x7F {
					maxChars++
					realRunes++
					englishChars++
				} else {
					maxChars += 2
					realRunes++
				}
			}
			// 如果显示不下, 那直接截取realRunes个就可以
			if maxChars >= 29 {
				dirContent.Text = string([]rune(dirContent.Text)[:realRunes])
			}
			// 中文存在双倍宽度问题, 这里使用手动补空格
			for ; maxChars < 30; maxChars++ {
				dirContent.Text += " "
			}

			res += fmt.Sprintf(fmtStr, dirContent.EntityType, dirContent.EntityId, dirContent.Text, dirContent.CreatedT, dirContent.ModifiedT)
		}
		LastReadDir = args[2]
	}
	if linkStr, err := (LinkProcessor{}.Process("link " + args[1] + " " + args[2])); err == nil {
		res = res + "\n当前位置: " + linkStr
	}
	return
}

type _article struct {
	Title     string
	Html      string
	CreatedT  string
	ModifiedT string
}

type _dir struct {
	List []_dirContent
}

type _dirContent struct {
	EntityId   int `json:"id"`
	EntityType int `json:"type"`
	Text       string
	CreatedT   string
	ModifiedT  string
}
