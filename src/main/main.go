package main

import (
	"common"
	"csv"
	"env"
	"ini"
)

func main() {
	env.CurrRegionRoot = "E:/svn/Papa2/branch/Resources/External/SM"

	rules := loadUpgrade("tmp.xlsx", "")
	if rules == nil {
		return
	}

	for _, rule := range rules {
		switch rule.Target {
		case common.ET_csv:
			csv.Execute(rule)
		case common.ET_ini:
			ini.Execute(rule)
		}
	}
}
