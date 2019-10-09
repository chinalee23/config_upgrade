package csv

import (
	"common"
	"env"
	"fmt"
)

type stAddCol struct {
	fpath string
	upg   *common.STOneUpgrade
	csv   *stCsv

	patterns map[string]string
	rules    []*common.STRule

	insertPos int

	defaultValue string
}

func add_col(upg *common.STOneUpgrade) {
	p := &stAddCol{
		fpath: "",
		upg:   upg,
		csv:   nil,

		defaultValue: "",
	}
	p.execute()
}

func (p *stAddCol) execute() {
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
		p.add_col_field()
	}
}

func (p *stAddCol) add_col_field() {
	field := p.patterns["field"]

	flag, idx := p.csv.hasField(field)
	if flag {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("列已存在"))
		return
	}

	region, ok := p.patterns["copy"]
	if !ok {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("新增列必须配置从其他大区拷贝, 例 [copy]_Dev"))
		return
	}

	copycsv := parseCsv(getCsvPath(region, p.upg.Item))
	if copycsv == nil {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("拷贝大区表格解析失败"))
		return
	}

	flag, idx = copycsv.hasField(field)
	if !flag {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("拷贝大区不存在该列"))
		return
	}

	// p.insertPos = len(p.csv.head.desc)
	p.insertPos = idx
	p.csv.head.desc = common.InsertSlice(p.csv.head.desc, p.insertPos, copycsv.head.desc[idx])
	p.csv.head.fieldName = common.InsertSlice(p.csv.head.fieldName, p.insertPos, copycsv.head.fieldName[idx])
	p.csv.head.fieldType = common.InsertSlice(p.csv.head.fieldType, p.insertPos, copycsv.head.fieldType[idx])

	p.handleDataRule()

	p.upg.SaveExecuteResult(common.EE_Success, fmt.Sprintf(""))
}

func (p *stAddCol) handleDataRule() {
	p.rules = common.ParseRule(p.upg.DataRule)
	for _, rule := range p.rules {
		if rule.R == "default" {
			p.defaultValue = rule.Param
		}
	}

	for i, line := range p.csv.data {
		sps := splitline(line)
		sps = common.InsertSlice(sps, p.insertPos, p.defaultValue)
		p.csv.data[i] = joinline(sps)
	}

	p.csv.savefile()
}
