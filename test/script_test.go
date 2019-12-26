package test

import (
	"fmt"
	"strings"
	"testing"
	//"fmt"
	//"encoding/json"

	"log"

	"github.com/qiniu/qlang"
	_ "github.com/qiniu/qlang/lib/builtin" // 导入 builtin 包
)

var strings_Exports = map[string]interface{}{
	"replacer": strings.NewReplacer,
}

func Test_qlang(t *testing.T) {

	qlang.Import("strings", strings_Exports) // 导入一个自定义的包，叫 strings（和标准库同名）
	ql := qlang.New()

	err := ql.SafeEval(`x = strings.replacer("?", "!").replace("hello, world???")`)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("x:", ql.Var("x")) // 输出 x: hello, world!!!

	expr := `ret = strings.replacer("年", "-", "月", "-", "日", "").replace("2019年12月26日 17:10")`
	err = ql.SafeEval(expr)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("ret:", ql.Var("ret")) // 输出 x: hello, world!!!
}
