package csv

import (
	"bufio"
	"common"
	"fmt"
	"io"
	"os"
	"strings"
)

func Execute(upg *common.STOneUpgrade) {
	switch upg.Changeway {
	case common.EC_add_csv:
		add_csv(upg)
	case common.EC_add_row:
		add_row(upg)
	}
}

type stCsvHead struct {
	desc      []string
	fieldName []string
	fieldType []string

	filedMap map[string]int
}

type stCsv struct {
	head  stCsvHead
	rows  map[string]string
	cols  [][]string
	lines []string
}

func parseCsv(path string) (rtn *stCsv) {
	rtn = nil

	lines := readCsv(path)
	if lines == nil {
		return
	}

	rtn = &stCsv{
		head: parseHead(lines),

		rows: make(map[string]string),
		cols: make([][]string, 0),
	}

	for _, v := range lines[3:] {
		sps := strings.Split(v, ",")

		rtn.rows[sps[0]] = v
		for ii, vv := range sps {
			if len(rtn.cols) < (ii + 1) {
				rtn.cols = append(rtn.cols, make([]string, 0))
			}
			rtn.cols[ii] = append(rtn.cols[ii], vv)
		}
	}

	return
}

func readCsv(path string) (lines []string) {
	lines = nil

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("readCsv, open file error", path, err)
		return
	}
	defer file.Close()

	bufReader := bufio.NewReader(file)

	lines = make([]string, 0)
	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("readCsv, readline error", path, err)
				return
			}
		}

		lines = append(lines, string(line))
	}

	if len(lines) < 3 {
		fmt.Println("readCsv, csv line error", path, len(lines))
		return nil
	}

	return
}

func splitline(line string) []string {
	return strings.Split(line, ",")
}

func parseHead(lines []string) (head stCsvHead) {
	head = stCsvHead{
		desc:      splitline(lines[0]),
		fieldName: splitline(lines[1]),
		fieldType: splitline(lines[2]),

		filedMap: make(map[string]int),
	}
	for i, v := range head.fieldName {
		if _, ok := head.filedMap[v]; !ok {
			head.filedMap[v] = i
		}
	}

	return
}
