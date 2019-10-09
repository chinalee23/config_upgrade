package main

import (
	"common"
	"fmt"
	"github.com/tealeg/xlsx"
	"strconv"
	"strings"
)

func readxlsx(path string, region string) (xlfile *xlsx.File, sheet *xlsx.Sheet) {
	xlfile = nil
	sheet = nil

	xlfile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("open upgrade error", err)
		return
	}

	if region == "" {
		if len(xlfile.Sheets) == 0 {
			fmt.Println("no sheet in upgrade.xlsx")
			return
		}
		sheet = xlfile.Sheets[0]
	} else {
		var ok bool
		sheet, ok = xlfile.Sheet[region]
		if !ok {
			fmt.Println("not exist region ", region)
			return
		}
	}

	return
}

func getCell(row *xlsx.Row, idx int, default_value string) string {
	if idx >= len(row.Cells) {
		return default_value
	} else {
		cell := row.Cells[idx].String()
		cell = strings.TrimSpace(cell)
		if cell == "" {
			return default_value
		} else {
			return cell
		}
	}
}

func parseRules(sheet *xlsx.Sheet) (upgs []*common.STOneUpgrade) {
	upgs = make([]*common.STOneUpgrade, 0)

	if len(sheet.Rows) < 2 {
		return
	}

	var target, changeway, item, data, datarule string
	for _, row := range sheet.Rows[1:] {
		target = getCell(row, 1, target)
		changeway = getCell(row, 2, changeway)
		item = getCell(row, 3, item)
		data = getCell(row, 4, data)
		datarule = getCell(row, 5, datarule)

		onerule := &common.STOneUpgrade{
			Item:     item,
			Data:     data,
			DataRule: datarule,

			Row: row,
		}
		onerule.ClearExecuteResult()

		ntarget, err := strconv.Atoi(target[:1])
		if err != nil {
			continue
		}
		nchangeway, err := strconv.Atoi(changeway[:1])
		if err != nil {
			continue
		}

		onerule.Target = common.E_Target(ntarget)
		onerule.Changeway = common.E_Change(nchangeway)

		upgs = append(upgs, onerule)
	}

	return
}

func loadUpgrade(path string, region string) (xlfile *xlsx.File, upgs []*common.STOneUpgrade) {
	xlfile = nil
	upgs = nil

	xlfile, sheet := readxlsx(path, region)
	if xlfile == nil {
		return
	}

	return xlfile, parseRules(sheet)
}
