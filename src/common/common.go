package common

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type E_Target int

const (
	_ E_Target = iota
	ET_csv
	ET_ini
)

type E_Change int

const (
	_ E_Change = iota
	EC_add_csv
	EC_add_row
	EC_add_column
	EC_change_field
)

type STOneUpgrade struct {
	Target    E_Target
	Changeway E_Change
	Item      string
	Data      string
	DataRule  string
}

type STRule struct {
	R     string
	Param string
}

func ParsePattern(s string) (patterns map[string]string) {
	patterns = make(map[string]string)

	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	sps := strings.Split(s, "\n")
	for _, v := range sps {
		v = strings.TrimSpace(v)
		if v[:1] != "[" {
			continue
		}

		idx := strings.Index(v, "]")
		if idx < 0 {
			continue
		}

		key := strings.TrimSpace(v[1:idx])
		value := strings.TrimSpace(v[idx+1:])
		patterns[key] = value
	}

	return
}

func ParseRule(s string) (rules []*STRule) {
	rules = make([]*STRule, 0)

	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	sps := strings.Split(s, "\n")
	for _, v := range sps {
		v = strings.TrimSpace(v)
		if v[:1] != "[" {
			continue
		}

		idx := strings.Index(v, "]")
		if idx < 0 {
			continue
		}

		rules = append(rules, &STRule{
			R:     strings.TrimSpace(v[1:idx]),
			Param: strings.TrimSpace(v[idx+1:]),
		})
	}

	return
}

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CopyFile(src string, dst string) (err error) {
	fsrc, err := os.Open(src)
	if err != nil {
		fmt.Println("CopyFile", "src err", src, err)
		return
	}
	defer fsrc.Close()

	fdst, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("CopyFile", "dst create err", dst, err)
		return
	}
	defer fdst.Close()

	_, err = io.Copy(fdst, fsrc)

	return
}

func InsertSlice(slice []string, idx int, e string) []string {
	rear := append([]string{}, slice[idx:]...)
	slice = append(slice[:idx], e)
	slice = append(slice, rear...)

	return slice
}
