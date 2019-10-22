package csv

import (
	"common"
	"env"
	"os"
	"path/filepath"
	"strings"
)

func Execute(upg *common.STOneUpgrade) {
	switch upg.Changeway {
	case common.EC_add_csv:
		add_csv(upg)
	case common.EC_add_row:
		add_row(upg)
	case common.EC_add_col:
		add_col(upg)
	case common.EC_del_col:
		del_col(upg)
	}
}

type stCsvHead struct {
	desc      []string
	fieldName []string
	fieldType []string

	filedMap map[string]int
}

type stCsv struct {
	path string

	head stCsvHead
	data []string

	rows map[string]string
	cols [][]string

	colnum int
}

func splitline(line string) []string {
	return strings.Split(line, ",")
}

func joinline(sps []string) string {
	return strings.Join(sps, ",")
}

func (p *stCsv) hasField(field string) (bool, int) {
	for i, v := range p.head.fieldName {
		if v == field {
			return true, i
		}
	}
	return false, -1
}

func (p *stCsv) clearData() {
	p.data = make([]string, 0)
}

func (p *stCsv) savefile() {
	file, _ := os.OpenFile(p.path, os.O_WRONLY|os.O_TRUNC, 0600)
	defer file.Close()

	content := make([]string, len(p.data)+3)
	content[0] = strings.Join(p.head.desc, ",")
	content[1] = strings.Join(p.head.fieldName, ",")
	content[2] = strings.Join(p.head.fieldType, ",")
	for i, v := range p.data {
		content[i+3] = v
	}

	file.WriteString(strings.Join(content, "\n"))
}

func parseCsv(path string) (rtn *stCsv) {
	rtn = nil

	lines := readCsv(path)
	if lines == nil {
		return
	}

	rtn = &stCsv{
		path: path,

		head: parseHead(lines),
		data: lines[3:],

		rows: make(map[string]string),
		cols: make([][]string, 0),
	}

	for _, v := range lines[3:] {
		sps := splitline(v)

		rtn.rows[sps[0]] = v
		for ii, vv := range sps {
			if len(rtn.cols) < (ii + 1) {
				rtn.cols = append(rtn.cols, make([]string, 0))
			}
			rtn.cols[ii] = append(rtn.cols[ii], vv)
		}
	}

	rtn.colnum = len(rtn.head.desc)

	return
}

func readCsv(path string) (lines []string) {
	lines = common.ReadFile(path)
	if lines == nil || len(lines) < 3 {
		return nil
	}
	return
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

func getkey(line string) string {
	return splitline(line)[0]
}

func getCsvPath(region string, name string) string {
	return filepath.Join(env.RegionRoot, filepath.Join(region, "CSV/"+name+".csv"))
}
