package csv

import (
	"common"
	"env"
	"fmt"
)

type stAddCsv struct {
	upg      *common.STOneUpgrade
	fpath    string
	patterns map[string]string

	copyregion string

	rules []*common.STRule

	csv *stCsv
}

func add_csv(upg *common.STOneUpgrade) {
	p := &stAddCsv{
		upg:        upg,
		fpath:      "",
		copyregion: "",
	}

	p.execute()
}

func (p *stAddCsv) execute() {
	p.fpath = getCsvPath(env.CurrRegion, p.upg.Item)

	// 文件已存在就不再拷贝了
	if common.IsPathExist(p.fpath) {
		fmt.Println("add csv", "file already exist:", p.upg.Item)
		return
	}

	// 只支持从其他大区直接拷贝
	p.patterns = common.ParsePattern(p.upg.Data)
	copyregion, ok := p.patterns["copy"]
	if !ok {
		fmt.Println("add csv [", p.upg.Item, "], pattern not [copy]")
		return
	}

	srcfile := getCsvPath(copyregion, p.upg.Item)
	if !common.IsPathExist(srcfile) {
		fmt.Println("add csv [", p.upg.Item, "] not exist in copy region [", copyregion, "]")
		return
	}

	common.CopyFile(srcfile, p.fpath)

	p.handleDataRule()
}

func (p *stAddCsv) handleDataRule() {
	p.rules = common.ParseRule(p.upg.DataRule)
	for _, rule := range p.rules {
		if rule.R == "clear" {
			p.clear()
		}
	}
}

func (p *stAddCsv) clear() {
	fmt.Println("clear", p.upg.Item)

	p.csv = parseCsv(p.fpath)
	lines := p.csv.lines[:3]

	writeLines(lines, p.fpath)
}
