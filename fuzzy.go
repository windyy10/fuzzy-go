package fuzzy

import (
	"sort"

	"github.com/windyy10/go-darts"
)

type MatchPosition struct {
	Begin  int
	Length int
}

type MatchResult struct {
	OriginWord []rune
	Positions  []MatchPosition
}

type FuzzyMatch interface {
	// ExactMatchSearch 精准匹配, trie树节点完全匹配返回成功
	ExactMatchSearch(word []rune) bool

	// CommonPrefixSearch 前缀匹配, trie树匹配到节点即返回成功, 会返回匹配上的所有节点
	// res 可以为nil, 只返回是否匹配成功
	CommonPrefixSearch(word []rune, res *MatchResult) bool

	// FuzzyMatchSearch 模糊匹配, 匹配整个word，是否命中trie树, 会返回匹配上的所有节点
	// res 可以为nil, 只返回是否匹配成功
	FuzzyMatchSearch(word []rune, res *MatchResult) bool
}

type MatchCreator interface {
	Insert(words ...string) MatchCreator
	Build() FuzzyMatch
}

func NewFuzzy() MatchCreator {
	return &dartsMatchCreator{}
}

type dartsMatchCreator struct {
	keys [][]rune
}

func (d *dartsMatchCreator) Insert(words ...string) MatchCreator {
	for _, word := range words {
		key := []rune(word)
		d.keys = append(d.keys, key)
	}
	return d
}

func less(left []rune, right []rune) bool {
	iLen := len(left)
	jLen := len(right)
	if iLen < jLen {
		for p := range iLen {
			if left[p] < right[p] {
				return true
			}
			if left[p] > right[p] {
				return false
			}
		}
		return true
	}
	if iLen > jLen {
		for p := 0; p < jLen; p++ {
			if left[p] < right[p] {
				return true
			}
			if left[p] > right[p] {
				return false
			}
		}
		return false
	}
	for p := 0; p < iLen; p++ {
		if left[p] < right[p] {
			return true
		}
		if left[p] > right[p] {
			return false
		}
	}
	return false
}

func (d *dartsMatchCreator) Build() FuzzyMatch {
	// 空列表BuildFromDAWG会报错, 直接返回一个全false的对象
	if len(d.keys) == 0 {
		return &emptyMatch{}
	}
	// 建树的词语必须按长度排序, 不然重复前缀的无法存储节点, 不会被匹配出来
	sort.Slice(d.keys, func(i, j int) bool {
		return less(d.keys[i], d.keys[j])
	})
	freq := make([]int, len(d.keys))
	for i := range freq {
		freq[i] = 1
	}
	darts := darts.BuildFromDAWG(d.keys, freq)
	return &dartsFuzzyMatch{
		darts: darts,
	}
}

type dartsFuzzyMatch struct {
	darts darts.Darts
}

func (d *dartsFuzzyMatch) ExactMatchSearch(word []rune) bool {
	return d.darts.ExactMatchSearch(word, 0)
}

func (d *dartsFuzzyMatch) CommonPrefixSearch(word []rune, res *MatchResult) bool {
	resPair := d.darts.CommonPrefixSearch(word, 0)
	if res != nil {
		res.OriginWord = word
		res.Positions = make([]MatchPosition, len(resPair))
		for i, r := range resPair {
			res.Positions[i].Begin = 0
			res.Positions[i].Length = r.PrefixLen
		}
	}
	return len(resPair) != 0
}

func (d *dartsFuzzyMatch) FuzzyMatchSearch(word []rune, res *MatchResult) bool {
	if res != nil {
		res.OriginWord = word
		res.Positions = make([]MatchPosition, 0)
	}
	ret := false
	for wordIndex := range word {
		resPair := d.darts.CommonPrefixSearch(word[wordIndex:], 0)
		if res != nil {
			for _, r := range resPair {
				res.Positions = append(res.Positions, MatchPosition{
					Begin:  wordIndex,
					Length: r.PrefixLen,
				})
				ret = true
			}
		} else {
			// 有匹配到结构时, 只关心是否模糊匹配上时, 可以直接返回
			if len(resPair) != 0 {
				return true
			}
		}
	}
	return ret
}

type emptyMatch struct{}

func (e *emptyMatch) FuzzyMatchSearch(word []rune, res *MatchResult) bool {
	return false
}

func (e *emptyMatch) ExactMatchSearch(word []rune) bool {
	return false
}

func (e *emptyMatch) CommonPrefixSearch(word []rune, res *MatchResult) bool {
	return false
}
