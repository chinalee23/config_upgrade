package common

import (
	"bufio"
	"fmt"
	"github.com/tealeg/xlsx"
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
	EC_del_row

	EC_add_col
	EC_del_col

	EC_mod_value
	EC_mod_field
)

type E_ExecuteResult string

const (
	_          E_ExecuteResult = ""
	EE_Fail                    = "Fail"
	EE_Success                 = "Success"
)

type STOneUpgrade struct {
	Target    E_Target
	Changeway E_Change
	Item      string
	Data      string
	DataRule  string

	Row *xlsx.Row
}

func (p *STOneUpgrade) SaveExecuteResult(rst string, desc string) {
	sz := len(p.Row.Cells)
	if sz < 8 {
		p.Row.AddCell()
	}
	p.Row.Cells[7].SetString(rst)

	sz = len(p.Row.Cells)
	if sz < 9 {
		p.Row.AddCell()
	}
	p.Row.Cells[8].SetString(desc)
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

func ParseRule(s string) (rules map[string]string) {
	rules = make(map[string]string)

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
		rules[key] = value
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

func RemoveSlice(slice []string, idx int) []string {
	slice = append(slice[:idx], slice[idx+1:]...)
	return slice
}

func ReadFile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("ReadFile, open error", err, path)
		return nil
	}
	defer file.Close()

	bufreader := bufio.NewReader(file)
	lines := make([]string, 0)
	for {
		data, _, err := bufreader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("ReadFile read error", err, path)
				return nil
			}
		}

		lines = append(lines, string(data))
	}
	return lines
}
