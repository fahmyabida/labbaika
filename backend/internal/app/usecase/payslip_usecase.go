package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fahmyabida/labbaika/internal/app/usecase/dto"
	"github.com/fahmyabida/labbaika/internal/logger"
	"github.com/gingfrederik/docx"
	"github.com/gocarina/gocsv"
)

type PayslipUsecase struct {
}

func NewPayslipUsecase() *PayslipUsecase {
	return &PayslipUsecase{}
}

func (u *PayslipUsecase) ConvertPayslip(ctx context.Context, file io.Reader) (filePath string, err error) {
	var rows []dto.PayslipRow
	if err = gocsv.Unmarshal(file, &rows); err != nil {
		return filePath, err
	}

	filePath = u.creatingNewDocx(ctx, rows)

	defer func() {
		go func() {
			time.Sleep(2 * time.Second)
			_ = os.Remove(filePath)
		}()
	}()

	return filePath, nil
}

func (p *PayslipUsecase) creatingNewDocx(ctx context.Context, rows []dto.PayslipRow) (filePath string) {
	doc := docx.NewFile()

	for _, row := range rows {
		if row.IsEmpty() {
			logger.Warn(ctx, nil, "Empty data is skipped", nil)
			continue
		}
		para := doc.AddParagraph()
		para.AddText(row.Name).Size(11)
		para = doc.AddParagraph()
		para.AddText(row.Date).Size(11)

		para = doc.AddParagraph()
		if row.IsCustom(row.RateDayNight) {
			para.AddText("Siang \t: " + row.ShiftDay + "\t\t\t\t = " + row.TotalDay).Size(11)
			para = doc.AddParagraph()
			para.AddText("Malam \t: " + row.ShiftNight + "\t\t\t\t = " + row.TotalNight).Size(11)
		} else {
			para.AddText("Siang \t: " + row.ShiftDay + "\t x " + row.RateDayNight + "\t = " + row.TotalDay).Size(11)
			para = doc.AddParagraph()
			para.AddText("Malam \t: " + row.ShiftNight + "\t x " + row.RateDayNight + "\t = " + row.TotalNight).Size(11)
		}

		para = doc.AddParagraph()
		if row.IsCustom(row.RateHalfDay) {
			para.AddText("1/2Hari \t: " + row.HalfDay + "\t\t\t\t = " + row.TotalHalfDay).Size(11)
		} else {
			para.AddText("1/2Hari \t: " + row.HalfDay + "\t x " + row.RateHalfDay + "\t = " + row.TotalHalfDay).Size(11)
		}

		para = doc.AddParagraph()
		if row.IsCustom(row.RateOvertime) {
			para.AddText("Overtime \t: " + row.Overtime + "\t\t\t\t = " + row.TotalOvertime).Size(11)
			para = doc.AddParagraph()
			para.AddText("Overload \t: " + row.Overload + "\t\t\t\t = " + row.TotalOverload).Size(11)
		} else {
			para.AddText("Overtime \t: " + row.Overtime + "\t x " + row.RateOvertime + "\t = " + row.TotalOvertime).Size(11)
			para = doc.AddParagraph()
			para.AddText("Overload \t: " + row.Overload + "\t x " + row.RateOvertime + "\t = " + row.TotalOverload).Size(11)
		}

		para = doc.AddParagraph()
		if row.Extra != "" {
			para.AddText("Extra  \t\t\t\t\t = " + row.Extra).Size(11)
		}

		para = doc.AddParagraph()
		para.AddText("Total \t\t\t\t\t = " + row.Total).Size(11)
		doc.AddParagraph()
		doc.AddParagraph()
	}

	filePath = fmt.Sprintf("./%v.docx", time.Now().Format("200601021504059999999990700"))
	if err := doc.Save(filePath); err != nil {
		logger.Error(ctx, nil, "Failed to save docx file", err)
		return filePath
	}
	return filePath
}
