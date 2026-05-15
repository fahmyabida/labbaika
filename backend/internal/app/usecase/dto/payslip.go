package dto

import "strings"

type PayslipRow struct {
	Date          string `csv:"Tanggal"`
	Name          string `csv:"Nama"`
	ShiftDay      string `csv:"Siang"`
	ShiftNight    string `csv:"Malam"`
	HalfDay       string `csv:"1/2 Hari"`
	Overtime      string `csv:"Overtime"`
	Overload      string `csv:"Overload"`
	RateDayNight  string `csv:"Gaji Pokok"`
	RateHalfDay   string `csv:"Gaji 1/2 Hari"`
	RateOvertime  string `csv:"Gaji OT/OL"`
	TotalDay      string `csv:"Hitung Siang"`
	TotalNight    string `csv:"Hitung Malam"`
	TotalHalfDay  string `csv:"Hitung 1/2 Hari"`
	TotalOvertime string `csv:"Hitung Overtime"`
	TotalOverload string `csv:"Hitung Overload"`
	Extra         string `csv:"Tunjangan / Extra"`
	Total         string `csv:"Total"`
}

func (p *PayslipRow) IsEmpty() bool {
	return (p.ShiftDay == "" || p.ShiftDay == "0") &&
		(p.ShiftNight == "" || p.ShiftNight == "0") &&
		(p.HalfDay == "" || p.HalfDay == "0") &&
		(p.Overtime == "" || p.Overtime == "0") &&
		(p.Overload == "" || p.Overload == "0")
}

func (p *PayslipRow) IsCustom(value string) bool {
	return strings.TrimSpace(strings.ToLower(value)) == "custom"
}
