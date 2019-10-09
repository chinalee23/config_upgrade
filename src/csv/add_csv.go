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
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("表格已存在"))
		return
	}

	// 只支持从其他大区直接拷贝
	p.patterns = common.ParsePattern(p.upg.Data)
	copyregion, ok := p.patterns["copy"]
	if !ok {
		p.upg.SaveExecuteResult(common.EE_Fail, "新增表必须配置从其他大区拷贝, 例 [copy]_Dev")
		return
	}

	srcfile := getCsvPath(copyregion, p.upg.Item)
	if !common.IsPathExist(srcfile) {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("拷贝大区不存在该表"))
		return
	}

	err := common.CopyFile(srcfile, p.fpath)
	if err != nil {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("拷贝表格失败"))
		return
	}

	p.handleDataRule()

	p.upg.SaveExecuteResult(common.EE_Success, "")
}

func (p *stAddCsv) handleDataRule() {
	p.rules = common.ParseRule(p.upg.DataRule)
	for _, rule := range p.rules {
		switch rule.R {
		case "clear":
			p.clear()
		}
	}
}

func (p *stAddCsv) clear() {
	fmt.Println("clear", p.upg.Item)

	p.csv = parseCsv(p.fpath)
	p.csv.clearData()
	p.csv.savefile()
}
