package csv

import (
	"common"
	"env"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func add_row(upg *common.STOneUpgrade) {
	fpath := filepath.Join(env.CurrRegionRoot, "CSV/"+upg.Item+".csv")
	if !common.IsPathExist(fpath) {
		fmt.Println("add_row", "file not eixst", upg.Item)
		return
	}

	ruleContent := common.ParseRule(upg.Content)
	switch ruleContent.Rule {
	case "key":
		add_row_key(ruleContent.Data, upg)
	case "raw":
		add_row_raw(ruleContent.Data, upg)
	}
}

// key模式：从指定大区下相应的表，拷贝指定key数据。可额外配置排序方式
func add_row_key(key string, upg *common.STOneUpgrade) {
	srcregion := ""
	sortorder := ""

	sps := strings.Split(upg.Rule, "\n")
	for _, v := range sps {
		rule := common.ParseRule(v)
		switch rule.Rule {
		case "copy":
			srcregion = rule.Data
		case "ascending":
			sortorder = "ascending"
		case "descending":
			sortorder = "descending"
		}
	}

	if srcregion == "" {
		fmt.Println("add_row_key, must copy from other region")
		return
	}

	dstfile := filepath.Join(env.CurrRegionRoot, "CSV/"+upg.Item+".csv")
	dstcsv := parseCsv(dstfile)
	if _, ok := dstcsv.rows[key]; ok {
		fmt.Println("key [", key, "already exist", dstfile)
		return
	}

	fmt.Println(sortorder)

	srcfile := filepath.Join(filepath.Dir(env.CurrRegionRoot), filepath.Join(srcregion, "CSV/"+upg.Item+".csv"))
	srccsv := parseCsv(srcfile)

	content, ok := srccsv.rows[key]
	if !ok {
		fmt.Println("add_row_key, key [", key, "] not exist in region [", srcregion, "]")
		return
	}

	if sortorder == "" {
		fd, _ := os.OpenFile(dstfile, os.O_RDWR|os.O_APPEND, 0644)
		fd.Write([]byte(content))
		fd.Close()

		fmt.Println("add_row_key, add key [", key, "] success, in the end", dstfile)
	} else {

	}
}

func add_row_raw(data string, upg *common.STOneUpgrade) {
	dstfile := filepath.Join(env.CurrRegionRoot, "CSV/"+upg.Item+".csv")
	dstcsv := parseCsv(dstfile)
	if _, ok := dstcsv.rows[key]; ok {
		fmt.Println("key [", key, "already exist", dstfile)
		return
	}
}
