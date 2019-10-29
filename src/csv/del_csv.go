package csv

import (
	"common"
	"env"
	"fmt"
	"os"
)

type stDelCsv struct {
	upg   *common.STOneUpgrade
	fpath string
}

func del_csv(upg *common.STOneUpgrade) {
	p := &stDelCsv{
		upg:   upg,
		fpath: "",
	}

	p.execute()
}

func (p *stDelCsv) execute() {
	p.fpath = getCsvPath(env.CurrRegion, p.upg.Item)

	if !common.IsPathExist(p.fpath) {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("表格不存在"))
		return
	}

	err := os.Remove(p.fpath)
	if err != nil {
		p.upg.SaveExecuteResult(common.EE_Fail, fmt.Sprintf("删除失败"))
		return
	}

	p.upg.SaveExecuteResult(common.EE_Success, "")
}
