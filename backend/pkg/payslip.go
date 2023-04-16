package pkg

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gingfrederik/docx"
	"github.com/labstack/echo/v4"
)

type Payslip struct {
}

func (p *Payslip) ConvertPayslipCsvToDocxfunc(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer func() {
		dst.Close()
	}()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	lines := readCsvFile(file.Filename)

	newFileName := creatingNewDocx(lines, file.Filename)
	go func() {
		time.Sleep(5 * time.Second)
		os.Remove(newFileName)
		os.Remove(file.Filename)
	}()
	return c.Attachment(newFileName, newFileName)
}

func creatingNewDocx(lines [][]string, fileName string) (newFileName string) {
	f := docx.NewFile()
	for index, row := range lines {
		if index == 0 {
			continue
		}
		// add text
		para := f.AddParagraph()
		para.AddText(row[1]).Size(11)
		para = f.AddParagraph()
		para.AddText(row[0]).Size(11)
		para = f.AddParagraph()
		para.AddText("Siang \t: " + row[2] + "\t x " + row[7] + "\t = " + row[10]).Size(11)
		para = f.AddParagraph()
		para.AddText("Malam \t: " + row[3] + "\t x " + row[7] + "\t = " + row[11]).Size(11)
		para = f.AddParagraph()
		para.AddText("1/2Hari \t: " + row[4] + "\t x " + row[8] + "\t = " + row[12]).Size(11)
		para = f.AddParagraph()
		para.AddText("Overtime \t: " + row[5] + "\t x " + row[9] + "\t = " + row[13]).Size(11)
		para = f.AddParagraph()
		para.AddText("Overload \t: " + row[6] + "\t x " + row[9] + "\t = " + row[14]).Size(11)
		if row[15] != "" {
			para = f.AddParagraph()
			para.AddText("Tunjangan \t\t\t\t = " + row[15]).Size(11)
		} else {
			para = f.AddParagraph()
		}
		para = f.AddParagraph()
		para.AddText("Total \t\t\t\t\t = " + row[16]).Size(11)
		f.AddParagraph()
		f.AddParagraph()
	}

	newFileName = fmt.Sprintf("./%v.docx", fileName)
	f.Save(newFileName)
	return newFileName
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
