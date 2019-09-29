package csv

import (
	"common"
	"env"
	"fmt"
	"path/filepath"
)

func add_csv(upg *common.STOneUpgrade) {
	dstfile := filepath.Join(env.CurrRegionRoot, "CSV/"+upg.Item+".csv")

	// 文件已存在就不再拷贝了
	if common.IsPathExist(dstfile) {
		fmt.Println("add csv", "file already exist:", upg.Item)
		return
	}

	// 只支持从其他大区直接拷贝
	rule := common.ParseRule(upg.Rule)
	if rule.Rule != "copy" {
		fmt.Println("add csv", "rule not [copy]")
		return
	}

	rootdir := filepath.Dir(env.CurrRegionRoot)
	srcfile := filepath.Join(rootdir, filepath.Join(rule.Data, "CSV/"+upg.Item+".csv"))

	err := common.CopyFile(srcfile, dstfile)
	if err != nil {
		fmt.Println("copy file", upg.Item, "from region", rule.Data, "fail")
	} else {
		fmt.Println("add csv", upg.Item, "success, from region", rule.Data)
	}
}
