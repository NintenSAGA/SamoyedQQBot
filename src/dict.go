package main

import (
	"bufio"
	"fmt"
	"github.com/armon/go-radix"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	DICT_API = "https://api.dictionaryapi.dev/api/v2/entries/en/"
)

type Dict struct {
	dictTree *radix.Tree
	Current  string
}

func createDict() *Dict {
	dict := Dict{}
	dict.initDict()
	return &dict
}

func (d *Dict) Reset() {
	d.resetDict()
}

func (d *Dict) resetDict() {
	m := d.dictTree.ToMap()
	for k, _ := range m {
		if len(k) <= 1 {
			m[k] = false
		} else {
			m[k] = true
		}
	}
	d.dictTree = radix.NewFromMap(m)
	d.Current = ""
}

func (d *Dict) initDict() {
	wordlistPath := os.Getenv("WORDLIST_PATH")
	if len(wordlistPath) == 0 {
		log.Panicln("Missing env WORDLIST_PATH")
	}

	fd, err := os.Open(wordlistPath)
	if err != nil {
		log.Panic(err)
	}
	defer fd.Close()

	tree := radix.New()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		word := scanner.Text()
		word = strings.TrimSpace(word)
		word = strings.ToLower(word)
		if len(word) <= 1 {
			tree.Insert(word, false)
		} else {
			tree.Insert(word, true)
		}
	}
	d.dictTree = tree
	d.Current = ""
}

func (d *Dict) CheckWord(word string) (bool, error) {
	reqUrl, _ := url.JoinPath(DICT_API, word)
	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Println(err.Error())
		return false, fmt.Errorf("小狗脑子转不过来了")
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("找不到这个词")
	}
	var unmarked bool
	if val, ex := d.dictTree.Get(word); ex {
		unmarked = val.(bool)
	} else {
		unmarked = true
	}

	if !unmarked {
		return false, fmt.Errorf("之前说过了")
	}

	if len(d.Current) != 0 {
		return strings.EqualFold(word[0:1], d.Current[len(d.Current)-1:]), fmt.Errorf("字母接不上")
	}

	return true, nil
}

func (d *Dict) markWord(word string) {
	d.dictTree.Insert(word, false)
}

func (d *Dict) getMatchingWord(curWord string) (result string, existed bool) {
	existed = false
	d.markWord(curWord)
	l := len(curWord)
	last := curWord[l-1:]
	last = strings.ToLower(last)
	d.dictTree.WalkPrefix(last, func(key string, value interface{}) bool {
		log.Println(key)
		if !(value.(bool)) {
			return false
		}
		result = key
		existed = true
		return true
	})

	if existed {
		d.markWord(result)
		d.Current = result
	}

	return
}
