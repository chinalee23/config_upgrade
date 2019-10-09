package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"time"
)

func main() {
	xlfile, _ := xlsx.OpenFile("tmp.xlsx")

	sheet := xlfile.Sheets[0]
	row := sheet.Rows[1]
	sz := len(row.Cells)

	var cell *xlsx.Cell
	if sz > 7 {
		cell = row.Cells[7]
	} else {
		cell = row.AddCell()
	}

	cell.SetString("success " + time.Now().UTC().String())

	err := xlfile.Save("tmp.xlsx")
	if err != nil {
		fmt.Println("save error", err)
		return
	}

	fmt.Println("done.......")
}
