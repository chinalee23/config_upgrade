package ini

import (
	"bufio"
	"common"
	"env"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var _initialized bool = false
var _ini_ map[string](map[string]string)

func _init() {
	fpath := filepath.Join(env.RegionRoot, filepath.Join(env.CurrRegion, "INI/config.ini"))

	file, err := os.Open(fpath)
	if err != nil {
		fmt.Println("open config.ini error", err)
		return
	}
	defer file.Close()

	bufreader := bufio.NewReader(file)
	lines := make([]string, 0)
	for {
		data, _, err := bufreader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("read config.ini error", err)
				return
			}
		}

		line := string(data)
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}

	load(lines)

	_initialized = true
}

func isNewSection(line string) (bool, string) {
	if !strings.HasPrefix(line, "[") || !strings.HasSuffix(line, "]") {
		return false, ""
	}
	return true, line[1 : len(line)-1]
}

func load(lines []string) {
	var currsection string
	for _, line := range lines {
		flag, section := isNewSection(line)
		if flag {
			_ini_[section] = make(map[string]string)
			currsection = section
			continue
		}

		sps := strings.Split(line, "=")
		if len(sps) != 2 {
			continue
		}
		key := sps[0]
		value := sps[1]

		_ini_[currsection][key] = value
	}
}

func Execute(upg *common.STOneUpgrade) {
	if !_initialized {
		_ini_ = make(map[string](map[string]string))
		_init()
	}

	switch upg.Changeway {
	case common.EC_add_row:
	case common.EC_del_row:
	}
}
