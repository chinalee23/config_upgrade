package csv

import (
	"common"
	"env"
	"fmt"
	// "path/filepath"
)

type stAddCol struct {
	fpath string
	upg   *common.STOneUpgrade
	csv   *stCsv

	patterns map[string]string
}

func add_col(upg *common.STOneUpgrade) {
	p := &stAddCol{
		fpath: "",
		upg:   upg,
		csv:   nil,
	}
	p.execute()
}

func (p *stAddCol) execute() {
	p.fpath = getCsvPath(env.CurrRegion, p.upg.Item)
	if !common.IsPathExist(p.fpath) {
		fmt.Println("add_col, file not exist", p.upg.Item)
		return
	}

	p.csv = parseCsv(p.fpath)
	if p.csv == nil {
		fmt.Println("add_col, parseCsv [", p.upg.Item, "] in [", env.CurrRegion, "] fail")
		return
	}

	p.patterns = common.ParsePattern(p.upg.Data)
	if _, ok := p.patterns["field"]; ok {
		p.add_col_field()
	}
}

func (p *stAddCol) add_col_field() {
	field := p.patterns["field"]

	region, ok := p.patterns["copy"]
	if !ok {
		fmt.Println("add_col_field, pattern not [copy]")
		return
	}

	copycsv := parseCsv(getCsvPath(region, p.upg.Item))
	if copycsv == nil {
		fmt.Println("add_col_field, copy csv [", p.upg.Item, "] not exist in region [", region, "]")
		return
	}

	flag, idx := copycsv.hasField(field)
	if !flag {
		fmt.Println("add_col_field, field [", field, "]not exist in [", p.upg.Item, "] of region [", region, "]")
		return
	}

	p.csv.head.desc = append(p.csv.head.desc, copycsv.head.desc[idx])
	p.csv.head.fieldName = append(p.csv.head.fieldName, copycsv.head.fieldName[idx])
	p.csv.head.fieldType = append(p.csv.head.fieldType, copycsv.head.fieldType[idx])

}
