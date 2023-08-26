package main

import (
	"fmt"
	"github.com/huichen/sego"
	"log"
	"math/rand"
	"os"
)

type Solver struct {
	seg *sego.Segmenter
}

func createSolver() *Solver {
	seg := sego.Segmenter{}
	dictPath := os.Getenv("DICT_PATH")
	if len(dictPath) == 0 {
		log.Panicln("Missing env DICT_PATH")
	}
	seg.LoadDictionary(dictPath)
	return &Solver{
		seg: &seg,
	}
}

func (s *Solver) generateAnswer(username string, content string) string {
	words := s.seg.Segment([]byte(content))
	log.Println(sego.SegmentsToString(words, true))
	idx := rand.Intn(len(words))
	result := fmt.Sprintf("%v 你好呀，我觉得你说的这句话中，\"%v\"最有道理", username, words[idx].Token().Text())
	return result
}
