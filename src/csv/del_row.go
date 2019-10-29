package csv

import (
	"common"
	"env"
	"fmt"
)

type stDelRow struct {
	upg   *common.STOneUpgrade
	fpath string
	csv   *stCsv

	patterns map[string]string
}

func del_row(upg *common.STOneUpgrade) {
	p := &stDelRow{
		upg:   upg,
		fpath: "",
		csv:   nil,
	}

	p.execute()
}

func (p *stDelRow) execute() {
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
	key, ok := p.patterns["key"]
	if !ok {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("必须配置key模式"))
		return
	}

	for i, l := range p.csv.data {
		k := getkey(l)
		if k == key {
			p.csv.data = common.RemoveSlice(p.csv.data, i)
			p.csv.savefile()
			p.upg.SaveExecuteResult(common.EE_Success, fmt.Sprintf(""))
			return
		}
	}

	p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("没有找到key"))
}
