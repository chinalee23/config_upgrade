package main

import (
	"common"
	"csv"
	"env"
	"fmt"
	"ini"
)

func main() {
	env.RegionRoot = "E:/svn/Papa2/branch/Resources/External"
	env.CurrRegion = "HMT"

	file := "upgrade.xlsx"
	xlfile, upgs := loadUpgrade(file, "")
	if upgs == nil {
		return
	}

	for _, rule := range upgs {
		switch rule.Target {
		case common.ET_csv:
			csv.Execute(rule)
		case common.ET_ini:
			ini.Execute(rule)
		}
	}

	saveResultXlsx(xlfile)

	fmt.Println("done.......")
}
