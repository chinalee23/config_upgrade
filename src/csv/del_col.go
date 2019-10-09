package csv

import (
	"common"
	"env"
	"fmt"
)

type stDelCol struct {
	fpath string

	upg *common.STOneUpgrade
	csv *stCsv

	patterns map[string]string
}

func del_col(upg *common.STOneUpgrade) {
	p := &stDelCol{
		fpath: "",
		upg:   upg,
		csv:   nil,
	}

	p.execute()
}

func (p *stDelCol) execute() {
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
	if _, ok := p.patterns["field"]; ok {
		p.del_col_field()
	}
}

func (p *stDelCol) del_col_field() {
	field := p.patterns["field"]

	flag, idx := p.csv.hasField(field)
	if !flag {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("删除的列不存在"))
		return
	}

	p.csv.head.desc = common.RemoveSlice(p.csv.head.desc, idx)
	p.csv.head.fieldName = common.RemoveSlice(p.csv.head.fieldName, idx)
	p.csv.head.fieldType = common.RemoveSlice(p.csv.head.fieldType, idx)

	for i, line := range p.csv.data {
		sps := splitline(line)
		sps = common.RemoveSlice(sps, idx)
		p.csv.data[i] = joinline(sps)
	}

	p.csv.savefile()

	p.upg.SaveExecuteResult(common.EE_Success, fmt.Sprintf(""))
}
