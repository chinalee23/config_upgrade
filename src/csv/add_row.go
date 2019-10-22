package csv

import (
	"common"
	"env"
	"fmt"
	"strings"
)

type stAddRow struct {
	fpath string
	csv   *stCsv
	upg   *common.STOneUpgrade

	key     string
	content string

	sortorder string
	patterns  map[string]string

	flag bool
}

func add_row(upg *common.STOneUpgrade) {
	p := &stAddRow{
		fpath: "",
		csv:   nil,
		upg:   upg,

		content: "",

		sortorder: "",

		flag: true,
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

	if p.flag {
		p.insert()
	}
}

func (p *stAddRow) add_row_key() {
	p.key = p.patterns["key"]
	slice := make([]string, p.csv.colnum)
	slice[0] = p.key
	p.content = strings.Join(slice, ",")
}

func (p *stAddRow) handleDataRule() {
	rules := common.ParseRule(p.upg.DataRule)

	if region, ok := rules["copy"]; ok {
		copycsv := parseCsv(getCsvPath(region, p.upg.Item))
		if copycsv == nil {
			p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("拷贝大区解析表格失败"))
			p.flag = false
			return
		}
		p.content, ok = copycsv.rows[p.key]
		if !ok {
			p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("拷贝大区不存在该key"))
			p.flag = false
			return
		}
	}

	if order, ok := rules["sort"]; ok {
		p.sortorder = order
	}

	p.insert()
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
