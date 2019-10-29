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

	flag bool
}

func add(upg *common.STOneUpgrade) {
	p := &stAdd{
		upg:   upg,
		value: "",
		flag:  true,
	}

	p.execute()
}

func (p *stAdd) execute() {
	p.key = createNewKey(p.upg.Item, p.upg.Data)
	if _, ok := _ini.keyvalue[p.key]; ok {
		p.upg.SaveExecuteResult(common.EE_Fail, "key已存在")
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

	if p.flag {
		p.insert()
	}
}

func (p *stAdd) handleDataRule() {
	rules := common.ParseRule(p.upg.DataRule)
	for k, v := range rules {
		switch k {
		case "copy":
			p.copy(v)
		case "default":
			p.value = v
		}
	}
}

func (p *stAdd) copy(region string) {
	copyini := loadIni(getFilePath(region))
	if copyini == nil {
		p.flag = false
		return
	}

	value, ok := copyini.keyvalue[p.key]
	if !ok {
		p.flag = false
		p.upg.SaveExecuteResult(common.EE_Fail, "拷贝大区不存在该key")
		return
	}

	p.value = value
}

func (p *stAdd) insert() {
	adjustIdx(p.insertIdx, 1)

	str := fmt.Sprintf("%s=%s", p.upg.Data, p.value)
	_ini.lines = common.InsertSlice(_ini.lines, p.insertIdx+1, str)

	save()
}
