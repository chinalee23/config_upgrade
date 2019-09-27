package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"strconv"
	"strings"
)

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

func readxlsx(path string) *xlsx.File {
	xlfile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("open error", err)
		return nil
	}
	return xlfile
}

func loadXlsx() {
	xlfile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Println("open error", err)
		return nil
	}

	sheet := xlfile.Sheets[0]
	target := ""
	changeway := ""
	item := ""
	content := ""
	rule := ""
	for _, row := range sheet.Rows[1:] {
		target = getCell(row, 1, target)
		changeway = getCell(row, 2, changeway)
		item = getCell(row, 3, item)
		content = getCell(row, 4, content)
		rule = getCell(row, 5, rule)

		onerule := stOneRule{
			item:    item,
			content: content,
			rule:    rule,
		}

		ntarget, _ := strconv.Atoi(target[:1])
		onerule.target = E_Target(ntarget)

		nchangeway, _ := strconv.Atoi(changeway[:1])
		onerule.changeway = E_Change(nchangeway)
	}
}

func main() {

}
