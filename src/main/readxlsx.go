package main

import (
	"common"
	"fmt"
	"github.com/tealeg/xlsx"
	"strconv"
	"strings"
)

func readxlsx(path string, region string) *xlsx.Sheet {
	xlfile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("open upgrade error", err)
		return nil
	}

	var sheet *xlsx.Sheet
	if region == "" {
		if len(xlfile.Sheets) == 0 {
			fmt.Println("no sheet in upgrade.xlsx")
			return nil
		}
		sheet = xlfile.Sheets[0]
	} else {
		var ok bool
		sheet, ok = xlfile.Sheet[region]
		if !ok {
			fmt.Println("not exist region ", region)
			return nil
		}
	}

	return sheet
}

func getCell(row *xlsx.Row, idx int, default_value string) string {
	if idx >= len(row.Cells) {
		return default_value
	} else {
		cell := row.Cells[idx].String()
		cell = strings.Trim(cell, " ")
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

	var target, changeway, item, content, rule string
	for _, row := range sheet.Rows[1:] {
		target = getCell(row, 1, target)
		changeway = getCell(row, 2, changeway)
		item = getCell(row, 3, item)
		content = getCell(row, 4, content)
		rule = getCell(row, 5, rule)

		onerule := &common.STOneUpgrade{
			Item:    item,
			Content: content,
			Rule:    rule,
		}
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

func loadUpgrade(path string, region string) []*common.STOneUpgrade {
	sheet := readxlsx(path, region)
	if sheet == nil {
		return nil
	}

	return parseRules(sheet)
}
