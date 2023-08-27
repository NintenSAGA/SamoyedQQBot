package main

import (
	"fmt"
	"strings"
)

type Solver struct {
	dict *Dict
}

func createSolver() *Solver {
	return &Solver{
		dict: createDict(),
	}
}

func (s *Solver) generateAnswer(username string, content string) string {
	if len(content) != 0 && isPureLetter(content) {
		if existed, err := s.dict.CheckWord(content); !existed {
			return fmt.Sprintf("汪汪，小狗觉得‘%v’这个词不太对呢！好像是因为%v！当前单词：‘%v’", content, err.Error(), s.dict.Current)
		}

		result, existed := s.dict.getMatchingWord(content)
		if !existed {
			s.dict.Current = content
			msg := fmt.Sprintf("小狗找不到了！小狗决定让你一盘！现在的单词是“%v”", content)
			return msg
		}

		return result
	} else {
		switch {
		case strings.HasPrefix(content, "重"):
			s.dict.Reset()
			return "汪汪好的！现在重新开始游戏～～"
		default:
			return fmt.Sprintf("汪汪，怎么了！？\n我是会玩英语单词接龙的小狗！\n现在的单词是“%v”\n"+
				"* 如果你要开始接龙，请输入一个单词！\n"+
				"* 如果要重开，那就告诉我“重来”！\n"+
				"* 对了，找我的时候不要忘记 @我 哦！", s.dict.Current)

		}

	}
}

func isPureLetter(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && r != '-' {
			return false
		}
	}
	return true
}
