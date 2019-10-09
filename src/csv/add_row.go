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
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("表格不存在"))
		return
	}

	p.csv = parseCsv(p.fpath)
	if p.csv == nil {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("解析表格失败"))
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
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("[key]模式必须配置从其他大区拷贝, 例 [copy]_Dev"))
		return
	}

	copycsv := parseCsv(getCsvPath(region, p.upg.Item))
	if copycsv == nil {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("拷贝大区解析表格失败"))
		return
	}
	p.content, ok = copycsv.rows[key]
	if !ok {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("拷贝大区不存在该key"))
		return
	}
}

func (p *stAddRow) handleDataRule() {
	p.rules = common.ParseRule(p.upg.DataRule)
	for _, rule := range p.rules {
		if rule.R == "ascend" || rule.R == "descend" {
			p.sortorder = rule.R
		}
	}
}

func (p *stAddRow) insert() {
	key := getkey(p.content)
	if _, ok := p.csv.rows[key]; ok {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("key已存在"))
		return
	}

	idx := len(p.csv.data)
	if p.sortorder != "" {
		for i, v := range p.csv.data {
			k := getkey(v)
			if (p.sortorder == "ascend" && k > key) ||
				(p.sortorder == "descend" && k < key) {
				idx = i
				break
			}
		}
	}

	p.csv.data = common.InsertSlice(p.csv.data, idx, p.content)

	p.csv.savefile()

	p.upg.SaveExecuteResult(common.EE_Success, fmt.Sprintf(""))
}
