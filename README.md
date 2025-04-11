# fuzzy-match
基于darts的模糊匹配库

## 示例 参考[单测代码](./fuzzy_test.go)

```go
package main

import (
	"fmt"

	"github.com/windyy10/fuzzy-go"
)

func main() {
	// 构建匹配字典
	trie := fuzzy.NewFuzzy().Insert("中华人民共和国", "原神", "王者荣耀").Build()

	// 进行精准匹配, 词典内存在词语与匹配词完全一致时成功
	trie.ExactMatchSearch([]rune("王者荣耀"))

	// 进行前缀匹配, 词典内存在词语与匹配词前缀一致时成功, 并且返回所有匹配上的词; res可以为nil, 仅判断是否匹配成功
	res := &fuzzy.MatchResult{}
	trie.CommonPrefixSearch([]rune("原神启动"), res)

	// 进行模糊匹配, 匹配词内存在词典内的词时成功, 并且返回所有匹配上的词; res可以为nil, 仅判断是否匹配成功
	trie.FuzzyMatchSearch([]rune("下载王者荣耀游戏"), res)
}
```