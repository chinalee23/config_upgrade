package ini

import (
	"common"
	"fmt"
)

type stAdd struct {
	upg       *common.STOneUpgrade
	value     string
	key       string
	insertIdx int
}

func add(upg *common.STOneUpgrade) {
	p := &stAdd{
		upg:   upg,
		value: "",
	}

	p.execute()
}

func (p *stAdd) execute() {
	p.key = createNewKey(p.upg.Item, p.upg.Data)
	if _, ok := _ini.keyvalue[p.key]; ok {
		fmt.Println("already exist this key")
		return
	}

	idx, ok := _ini.sectionLastIdx[p.upg.Item]
	if !ok {
		_ini.lines = append(_ini.lines, fmt.Sprintf("[%s]", p.upg.Item))
		idx = len(_ini.lines) - 1
		_ini.sectionLastIdx[p.upg.Item] = idx
	}
	p.insertIdx = idx

	p.handleDataRule()
}

func (p *stAdd) handleDataRule() {
	rules := common.ParseRule(p.upg.DataRule)
	for k, v := range rules {
		switch k {
		case "copy":
			p.copy(v)
		case "default":
			p.value = v
			p.insert()
		}
	}
}

func (p *stAdd) copy(region string) {
	copyini := loadIni(getFilePath(region))
	if copyini == nil {
		return
	}

	value, ok := copyini.keyvalue[p.key]
	if !ok {
		fmt.Println("key", p.upg.Data, "not exist in region", region)
		return
	}

	p.value = value
	p.insert()
}

func (p *stAdd) insert() {
	adjustIdx(p.insertIdx, 1)

	str := fmt.Sprintf("%s=%s", p.upg.Data, p.value)
	_ini.lines = common.InsertSlice(_ini.lines, p.insertIdx+1, str)

	save()
}
