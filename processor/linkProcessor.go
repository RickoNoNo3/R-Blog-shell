package processor

import (
	"encoding/json"
	"fmt"

	"../myerror"
)

type LinkProcessor struct{}

// 把BlogCore发回来的Json数据处理成文字表格.
func (p LinkProcessor) Process(input string) (res string, err error) {
	// 判断是否是link命令
	args := splitArgs(input)
	if args[0] != "link" {
		err = myerror.CannotProcessError
		return
	}
	// 连接BlogCore
	if res, err = runBlogCore(args); err != nil {
		return
	}
	// 处理json数据
	jsonStr := res
	res = ""
	var linkRes _linkResult
	if err = json.Unmarshal([]byte(jsonStr), &linkRes); err != nil {
		return
	}
	for i, v := range linkRes.Link {
		res += fmt.Sprintf("[%v]%v", v.Id, v.Title)
		if i < len(linkRes.Link) - 1 {
			res += " > "
		} else {
			res += "\n"
		}
	}
	return
}

type _linkResult struct {
	Link []_linkResultContent `json:"link"`
}

type _linkResultContent struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}
