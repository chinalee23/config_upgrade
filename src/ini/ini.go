package ini

import (
	"common"
	"env"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type stIni struct {
	path string

	lines          []string
	keyidx         map[string]int
	keyvalue       map[string]string
	sectionLastIdx map[string]int
}

func isNewSection(line string) (bool, string) {
	if !strings.HasPrefix(line, "[") || !strings.HasSuffix(line, "]") {
		return false, ""
	}
	return true, line[1 : len(line)-1]
}

func createNewKey(section string, key string) string {
	return section + "|" + key
}

func loadIni(path string) *stIni {
	lines := common.ReadFile(path)
	if lines == nil {
		fmt.Println("load ini error", path)
		return nil
	}

	p := &stIni{
		path:           path,
		lines:          lines,
		keyidx:         make(map[string]int),
		keyvalue:       make(map[string]string),
		sectionLastIdx: make(map[string]int),
	}

	var currsection string = ""
	var lastidx int
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		flag, section := isNewSection(line)
		if flag {
			if currsection != "" {
				p.sectionLastIdx[currsection] = lastidx
			}
			currsection = section
			lastidx = i
			p.sectionLastIdx[currsection] = lastidx
			continue
		}

		if currsection == "" {
			continue
		}

		lastidx = i

		sps := strings.Split(line, "=")
		if len(sps) != 2 {
			fmt.Println("ini error config", line)
			return nil
		}

		key := createNewKey(currsection, sps[0])
		p.keyidx[key] = i
		p.keyvalue[key] = sps[1]
	}
	if currsection != "" {
		p.sectionLastIdx[currsection] = lastidx
	}

	return p
}

func getFilePath(region string) string {
	return filepath.Join(env.RegionRoot, filepath.Join(region, "INI/config.ini"))
}

func save() {
	path := getFilePath(env.CurrRegion)
	file, _ := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)
	defer file.Close()

	file.WriteString(strings.Join(_ini.lines, "\n"))
}

func adjustIdx(idx int, offset int) {
	for k, v := range _ini.keyidx {
		if v > idx {
			_ini.keyidx[k] = v + offset
		}
	}
	for k, v := range _ini.sectionLastIdx {
		if v > idx {
			_ini.sectionLastIdx[k] = v + offset
		}
	}
}

var _initialized bool = false
var _ini *stIni

func Execute(upg *common.STOneUpgrade) {
	if !_initialized {
		_ini = loadIni(getFilePath(env.CurrRegion))
		_initialized = true
	}

	if _ini == nil {
		upg.SaveExecuteResult(common.EE_Fail, "ini加载失败")
		return
	}

	switch upg.Changeway {
	case common.EC_add_row:
		add(upg)
	case common.EC_del_row:
		del(upg)
	}
}
