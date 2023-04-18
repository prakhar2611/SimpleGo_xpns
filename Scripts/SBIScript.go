package Scripts

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func ExtractSBI(FileName string) {
	// Reading all files of directory
	// Iterating through all the files
	f, err := excelize.OpenFile("SBIFiles/" + FileName)
	if err != nil {
		log.Fatal(err)
	}
	c1, err := f.GetCellValue("Sheet1", "A1")
	fmt.Println(c1)
}
