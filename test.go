package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println("open error", err)
		return
	}

	err = f.SetCellStr("Sheet2", "C2", "Success")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = f.Save()
	if err != nil {
		fmt.Println(err)
	}
}
