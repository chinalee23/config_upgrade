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
	Content   string
	Rule      string
}

type STOneRule struct {
	Rule string
	Data string
}

func ParseRule(s string) (rule *STOneRule) {
	rule = &STOneRule{
		Rule: "_none_",
		Data: "",
	}

	if s[:1] != "[" {
		return
	}

	idx := strings.Index(s, "]")
	if idx < 0 {
		return
	}

	rule.Rule = strings.Trim(s[1:idx], " ")
	rule.Data = strings.Trim(s[idx+1:], " ")

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
