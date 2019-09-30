package csv

import (
	"common"
	"env"
	"fmt"
)

type stAddRow struct {
	fpath string
	csv   *stCsv
	upg   *common.STOneUpgrade

	content string

	sortorder string
	patterns  map[string]string
	rules     []*common.STRule
}

func add_row(upg *common.STOneUpgrade) {
	p := &stAddRow{
		fpath: "",
		csv:   nil,
		upg:   upg,

		content: "",

		sortorder: "",
	}

	p.execute()
}

func (p *stAddRow) execute() {
	p.fpath = getCsvPath(env.CurrRegion, p.upg.Item)
	if !common.IsPathExist(p.fpath) {
		fmt.Println("add_row", "file not eixst", p.upg.Item)
		return
	}

	p.csv = parseCsv(p.fpath)
	if p.csv == nil {
		fmt.Println("add_row, parseCsv [", p.upg.Item, "] in [", env.CurrRegion, "] fail")
		return
	}

	p.patterns = common.ParsePattern(p.upg.Data)
	if _, ok := p.patterns["key"]; ok {
		p.add_row_key()
	} else if _, ok := p.patterns["raw"]; ok {
		p.content = p.patterns["raw"]
	}

	p.handleDataRule()

	p.insert()
}

func (p *stAddRow) add_row_key() {
	key := p.patterns["key"]

	region, ok := p.patterns["copy"]
	if !ok {
		fmt.Println("add_row_key, pattern not [copy]")
		return
	}

	copycsv := parseCsv(getCsvPath(region, p.upg.Item))
	if copycsv == nil {
		fmt.Println("add_row_key, copy csv[", p.upg.Item, "] not exist in region[", region, "]")
		return
	}
	p.content, ok = copycsv.rows[key]
	if !ok {
		fmt.Println("add_row_key, key [", key, "] not exist in copy region [", region, "]")
		return
	}
}

func (p *stAddRow) handleDataRule() {
	p.rules = common.ParseRule(p.upg.DataRule)
	for _, rule := range p.rules {
		if rule.R == "ascending" || rule.R == "descending" {
			p.sortorder = rule.R
		}
	}
}

func (p *stAddRow) insert() {
	key := getkey(p.content)
	if _, ok := p.csv.rows[key]; ok {
		fmt.Println("add_row, insert key [", key, "] already exist", p.upg.Item)
		return
	}

	idx := len(p.csv.lines)
	fmt.Println("idx", idx)
	if p.sortorder != "" {
		for i, v := range p.csv.lines[3:] {
			k := getkey(v)
			if (p.sortorder == "ascending" && k > key) ||
				(p.sortorder == "descending" && k < key) {
				idx = i + 3
				break
			}
		}
	}

	fmt.Println("idx", idx)

	rear := append([]string{}, p.csv.lines[idx:]...)
	p.csv.lines = append(p.csv.lines[:idx], p.content)
	p.csv.lines = append(p.csv.lines, rear...)

	writeLines(p.csv.lines, p.fpath)
}

func printLines(lines []string) {
	for _, v := range lines {
		fmt.Println(v)
	}
	fmt.Println("----------")
}
