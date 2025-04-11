package fuzzy_test

import (
	"testing"

	"github.com/windyy10/fuzzy-go"

	"github.com/stretchr/testify/require"
)

func TestFuzzyBuildEmpty(t *testing.T) {
	match := fuzzy.NewFuzzy().Build()
	require.Equal(t, false, match.CommonPrefixSearch([]rune("测试"), nil))
}

func TestFuzzyExactMatchSearch(t *testing.T) {
	fu := fuzzy.NewFuzzy()
	match := fu.Insert("中华人民共和国", "原神", "星穹铁道").Build()
	require.Equal(t, true, match.ExactMatchSearch([]rune("原神")))
	require.Equal(t, false, match.ExactMatchSearch([]rune("崩坏·星穹铁道")))
}

func TestFuzzyMatchSearch(t *testing.T) {
	fu := fuzzy.NewFuzzy()
	match := fu.Insert("中华人民共和国", "原神", "星穹铁道").Build()
	require.Equal(t, true, match.FuzzyMatchSearch([]rune("星穹铁道"), nil))
	require.Equal(t, true, match.FuzzyMatchSearch([]rune("崩坏·星穹铁道"), nil))
}

func getResWords(res *fuzzy.MatchResult) []string {
	words := make([]string, len(res.Positions))
	for i, r := range res.Positions {
		words[i] = string(res.OriginWord[r.Begin : r.Begin+r.Length])
	}
	return words
}

func TestFuzzyMatchSearchSubWords(t *testing.T) {
	fu := fuzzy.NewFuzzy()
	match := fu.Insert("中华人民共和国", "原神", "崩坏·星穹铁道", "崩铁", "崩坏", "原", "礼包", "周年庆", "星穹铁道", "福利").Build()
	res := &fuzzy.MatchResult{}
	require.Equal(t, true, match.FuzzyMatchSearch([]rune("原神"), res))
	require.Equal(t, []string{"原", "原神"}, getResWords(res))

	require.Equal(t, true, match.FuzzyMatchSearch([]rune("崩坏·星穹铁道周年庆福利礼包"), res))
	require.Equal(t, []string{"崩坏", "崩坏·星穹铁道", "星穹铁道", "周年庆", "福利", "礼包"}, getResWords(res))
}
