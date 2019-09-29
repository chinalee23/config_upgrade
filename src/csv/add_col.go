package csv

import (
	"common"
	"env"
)

func add_col(upg *common.STOneUpgrade) {
	fpath := filepath.Join(env.CurrRegionRoot, "CSV/"+upg.Item+".csv")
	if !common.IsPathExist(fpath) {
		fmt.Println("add_col", "file not eixst", upg.Item)
		return
	}

	ruleContent := common.ParseRule(upg.Content)
	switch ruleContent.Rule {
	case "field":

	}
}
