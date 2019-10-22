package ini

import (
	"common"
	"fmt"
)

type stDel struct {
	upg *common.STOneUpgrade
	key string
}

func del(upg *common.STOneUpgrade) {
	p := &stDel{
		upg: upg,
	}

	p.execute()
}

func (p *stDel) execute() {
	p.key = createNewKey(p.upg.Item, p.upg.Data)
	idx, ok := _ini.keyidx[p.key]
	if !ok {
		fmt.Println("not exist key", p.key)
		return
	}

	adjustIdx(idx, -1)

	_ini.lines = common.RemoveSlice(_ini.lines, idx)

	save()
}
