package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gingfrederik/docx"
)

type PayslipUsecase struct {
}

func NewPayslipUsecase() *PayslipUsecase {
	return &PayslipUsecase{}
}

func (u *PayslipUsecase) ConvertPayslip(c context.Context, file io.Reader) (filePath string, err error) {
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return filePath, err
	}

	filePath = u.creatingNewDocx(records)

	defer func() {
		go func() {
			time.Sleep(5 * time.Second)
			os.Remove(filePath)
		}()
	}()

	return filePath, nil
}

func (p *PayslipUsecase) creatingNewDocx(lines [][]string) (filePath string) {
	doc := docx.NewFile()

	for index, row := range lines {
		if index == 0 {
			continue
		}
		// add text
		para := doc.AddParagraph()
		para.AddText(row[1]).Size(11)
		para = doc.AddParagraph()
		para.AddText(row[0]).Size(11)
		para = doc.AddParagraph()
		para.AddText("Siang \t: " + row[2] + "\t x " + row[7] + "\t = " + row[10]).Size(11)
		para = doc.AddParagraph()
		para.AddText("Malam \t: " + row[3] + "\t x " + row[7] + "\t = " + row[11]).Size(11)
		para = doc.AddParagraph()
		para.AddText("1/2Hari \t: " + row[4] + "\t x " + row[8] + "\t = " + row[12]).Size(11)
		para = doc.AddParagraph()
		para.AddText("Overtime \t: " + row[5] + "\t x " + row[9] + "\t = " + row[13]).Size(11)
		para = doc.AddParagraph()
		para.AddText("Overload \t: " + row[6] + "\t x " + row[9] + "\t = " + row[14]).Size(11)
		if row[15] != "" {
			para = doc.AddParagraph()
			para.AddText("Extra  \t\t\t\t\t = " + row[15]).Size(11)
		} else {
			doc.AddParagraph()
		}
		para = doc.AddParagraph()
		para.AddText("Total \t\t\t\t\t = " + row[16]).Size(11)
		doc.AddParagraph()
		doc.AddParagraph()
	}

	filePath = fmt.Sprintf("./%v.docx", time.Now().Format("200601021504059999999990700"))
	doc.Save(filePath)
	return filePath
}
