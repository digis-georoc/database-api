package download

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
)

const (
	CSV  = "csv"
	XLSX = "xlsx"
)

type Formatter interface {
	FormatData(samples []model.FullData) ([]byte, error)
}

func GetFormatter(targetFormat string) Formatter {
	if targetFormat == CSV {
		return NewCSVFormatter(",")
	}
	if targetFormat == XLSX {
		return NewXLSXFormatter()
	}
	return nil
}

// CSV Formatter for csv files
type CSVFormatter struct {
	separator string
}

func NewCSVFormatter(separator string) Formatter {
	return &CSVFormatter{separator: separator}
}

func (f *CSVFormatter) FormatData(samples []model.FullData) ([]byte, error) {
	rows := [][]string{
		{"YEAR", "CITATION", "SAMPLE NAME", "UNIQUE_ID", "LOCATION", "ELEVATION (MIN.)", "ELEVATION (MAX.)", "SAMPLING TECHNIQUE", "DRILLING DEPTH (MIN.)", "DRILLING DEPTH (MAX.)", "LAND/SEA (SAMPLING)", "ROCK TYPE", "ROCK NAME", "ROCK TEXTURE", "SAMPLE COMMENT", "AGE (MIN.)", "AGE (MAX.)", "GEOLOGICAL AGE", "GEOLOGICAL AGE PREFIX", "ERUPTION DAY", "ERUPTION MONTH", "ERUPTION YEAR", "ALTERATION", "ALTERATION TYPE", "TYPE OF MATERIAL", "MINERAL / COMPONENT", "SPOT (MINERAL)", "CRYSTAL", "RIM / CORE (MINERAL GRAINS)", "PRIMARY/SECONDARY", "INCLUSION TYPE", "MINERAL (INCLUSION)", "RIM / CORE (INCLUSION)", "HOSTMINERAL (INCLUSION)", "LATITUDE (MIN.)", "LONGITUDE (MIN.)", "LATITUDE (MAX.)", "LONGITUDE (MAX.)", "SIO2(WT%)", "TIO2(WT%)", "AL2O3(WT%)", "CR2O3(WT%)", "FE2O3T(WT%)", "FE2O3(WT%)", "FEOT(WT%)", "FEO(WT%)", "CAO(WT%)", "MGO(WT%)", "MNO(WT%)", "BAO(WT%)", "K2O(WT%)", "NA2O(WT%)", "P2O5(WT%)", "H2OP(WT%)", "H2OM(WT%)", "CL(WT%)", "SO3(WT%)", "S(WT%)", "LOI(WT%): altern. values or methods", "", "", "", "SR(WT%)", "LI(PPM): altern. values or methods", "", "", "", "BE(PPM): altern. values or methods", "", "", "", "", "", "B(PPM)", "F(PPM)", "NA(PPM): altern. values or methods", "", "", "", "MG(PPM): altern. values or methods", "", "", "", "AL(PPM): altern. values or methods", "", "", "", "SI(PPM)", "P(PPM): altern. values or methods", "", "", "", "S(PPM): altern. values or methods", "", "", "", "CL(PPM)", "K(PPM): altern. values or methods", "", "", "", "CA(PPM): altern. values or methods", "", "", "", "SC(PPM): altern. values or methods", "", "", "", "", "", "TI(PPM): altern. values or methods", "", "", "", "V(PPM): altern. values or methods", "", "", "", "", "", "CR(PPM): altern. values or methods", "", "", "", "MN(PPM): altern. values or methods", "", "", "", "FE(PPM): altern. values or methods", "", "", "", "CO(PPM): altern. values or methods", "", "", "", "", "", "NI(PPM): altern. values or methods", "", "", "", "", "", "CU(PPM): altern. values or methods", "", "", "", "ZN(PPM): altern. values or methods", "", "", "", "GA(PPM): altern. values or methods", "", "", "", "", "", "GE(PPM)", "AS(PPM): altern. values or methods", "", "", "", "SE(PPM): altern. values or methods", "", "", "", "BR(PPM)", "RB(PPM): altern. values or methods", "", "", "", "", "", "SR(PPM): altern. values or methods", "", "", "", "", "", "Y(PPM): altern. values or methods", "", "", "", "", "", "ZR(PPM): altern. values or methods", "", "", "", "", "", "NB(PPM): altern. values or methods", "", "", "", "", "", "MO(PPM): altern. values or methods", "", "", "", "PD(PPM): altern. values or methods", "", "", "", "PD(PPB)", "AG(PPM): altern. values or methods", "", "", "", "CD(PPM): altern. values or methods", "", "", "", "IN(PPM): altern. values or methods", "", "", "", "SN(PPM): altern. values or methods", "", "", "", "", "", "SB(PPM): altern. values or methods", "", "", "", "TE(PPM): altern. values or methods", "", "", "", "CS(PPM): altern. values or methods", "", "", "", "", "", "BA(PPM): altern. values or methods", "", "", "", "", "", "LA(PPM): altern. values or methods", "", "", "", "", "", "CE(PPM): altern. values or methods", "", "", "", "", "", "PR(PPM): altern. values or methods", "", "", "", "ND(PPM): altern. values or methods", "", "", "", "SM(PPM): altern. values or methods", "", "", "", "EU(PPM): altern. values or methods", "", "", "", "GD(PPM): altern. values or methods", "", "", "", "TB(PPM): altern. values or methods", "", "", "", "DY(PPM): altern. values or methods", "", "", "", "HO(PPM): altern. values or methods", "", "", "", "ER(PPM): altern. values or methods", "", "", "", "TM(PPM): altern. values or methods", "", "", "", "YB(PPM): altern. values or methods", "", "", "", "LU(PPM): altern. values or methods", "", "", "", "HF(PPM): altern. values or methods", "", "", "", "", "", "TA(PPM): altern. values or methods", "", "", "", "", "", "W(PPM): altern. values or methods", "", "", "", "", "", "RE(PPM): altern. values or methods", "", "", "", "PT(PPM): altern. values or methods", "", "", "", "PT(PPB)", "AU(PPM): altern. values or methods", "", "", "", "AU(PPB)", "HG(PPM)", "TL(PPM): altern. values or methods", "", "", "", "PB(PPM): altern. values or methods", "", "", "", "PB204(PPM)", "PB206(PPM)", "PB207(PPM)", "PB208(PPM)", "BI(PPM): altern. values or methods", "", "", "", "TH(PPM): altern. values or methods", "", "", "", "", "", "U(PPM): altern. values or methods", "", "", "", "", "", "U238(PPM)", "ND143_ND144", "EPSILON_ND_INI", "SM147_ND144", "SR87_SR86", "RB87_SR86", "PB206_PB204", "PB207_PB204", "PB207_PB206", "PB208_PB204", "PB208_PB206", "HF176_HF177", "PB206_U238", "PB207_U235", "D13C(VS VPDB)", "D18O(VS VPDB)", "D18O(VS VSMOW)", "D34S(VS VCDT)", "ALBITE(MOL%)", "ALMANDINE(MOL%)", "ANORTHITE(MOL%)", "ENSTATITE(MOL%)", "FERROSILITE(MOL%)", "GROSSULAR(MOL%)", "ORTHOCLASE(MOL%)", "PYROPE(MOL%)", "SPESSARTINE(MOL%)", "WOLLASTONITE(MOL%)", "AGE_PB207_PB206(MA)", "AGE_PB207_U235(MA)", "AGE_PB206_U238(MA)"},
	}
	for _, sample := range samples {
		row := make([]string, len(rows[0]))
		if len(sample.References) > 0 {
			ref := sample.References[0]
			row[0] = getInt(ref.Publicationyear)
			row[1] = fmt.Sprintf("[%v]%v", getString(ref.Externalidentifier), getString(ref.Title))
		}
		row[2] = getString(sample.SampleName)
		row[3] = getString(sample.UniqueID)
		row[4] = join(sample.LocationNames, "/")
		row[5] = getString(sample.ElevationMin)
		row[6] = getString(sample.ElevationMax)

		resultMap := map[string]string{}
		for _, batch := range sample.BatchData {
			for _, result := range batch.Results {
				itemName := getString(result.ItemName)
				value := getFloat64(result.Value)
				if itemName == "" || value == "" {
					continue
				}
				resultMap[itemName] = value
			}
		}
		// TODO: Add all fields
		rows = append(rows, row)
	}
	csv := ""
	for _, row := range rows {
		csv += strings.Join(row, ",")
		csv += "\n"
	}
	data := []byte(csv)
	return data, nil
}

// XSLX Formatter for excel files
type XLSXFormatter struct {
}

func NewXLSXFormatter() Formatter {
	return &XLSXFormatter{}
}

func (f *XLSXFormatter) FormatData(samples []model.FullData) ([]byte, error) {
	data := make([]byte, 0)
	return data, nil
}

// join implemebts strings.Join() for type []*string
func join(stringSlice []*string, delimiter string) string {
	join := ""
	for i, s := range stringSlice {
		if s != nil {
			if i > 0 {
				join += delimiter
			}
			join += *s
		}
	}
	return join
}

// getString returns the string s refers to or empty string
func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// getInt returns the integer s refers to as a string or empty string
func getInt(i *int) string {
	if i == nil {
		return ""
	}
	return strconv.Itoa(*i)
}

// getFloat64 returns the float64 s refers to as a string or empty string
func getFloat64(f *float64) string {
	if f == nil {
		return ""
	}
	return strconv.FormatFloat(*f, 'f', -1, 64)
}
