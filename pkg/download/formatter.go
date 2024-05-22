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

	// csv column keys
	KEY_YEAR               = "YEAR"
	KEY_CITATION           = "CITATION"
	KEY_SAMPLENAME         = "SAMPLE NAME"
	KEY_UNIQUE_ID          = "UNIQUE_ID"
	KEY_LOCATION           = "LOCATION"
	KEY_ELEVATION_MIN      = "ELEVATION (MIN.)"
	KEY_ELEVATION_MAX      = "ELEVATION (MAX.)"
	KEY_SAMPLING_TECHNIQUE = "SAMPLING TECHNIQUE"
	KEY_DRILLDEPTH_MIN     = "DRILLING DEPTH (MIN.)"
	KEY_DRILLDEPTH_MAX     = "DRILLING DEPTH (MAX.)"
	KEY_LANDORSEA          = "LAND/SEA (SAMPLING)"
	KEY_ROCKTYPE           = "ROCK TYPE"
	KEY_ROCKNAME           = "ROCK NAME"
	KEY_ROCKTEXTURE        = "ROCK TEXTURE"
	KEY_SAMPLECOMMENT      = "SAMPLE COMMENT"
	KEY_AGE_MIN            = "AGE (MIN.)"
	KEY_AGE_MAX            = "AGE (MAX.)"
	KEY_GEO_AGE            = "GEOLOGICAL AGE"
	KEY_AGE_PREFIX         = "GEOLOGICAL AGE PREFIX"
	KEY_ERUPTION_DATE      = "ERUPTION DATE"
	KEY_ALTERATION         = "ALTERATION"
	KEY_ALTERATION_TYPE    = "ALTERATION TYPE"
	KEY_MATERIAL_TYPE      = "TYPE OF MATERIAL"
	KEY_MINERAL            = "MINERAL / COMPONENT"
	KEY_CRYSTAL            = "CRYSTAL"
	KEY_RIMORCORE          = "RIM / CORE (MINERAL GRAINS)"
	KEY_INCLUSIONTYPE      = "INCLUSION TYPE"
	KEY_INCLUSION_MINERAL  = "MINERAL (INCLUSION)"
	KEY_RIMORCORE_INC      = "RIM / CORE (INCLUSION)"
	KEY_HOST_MINERAL       = "HOSTMINERAL (INCLUSION)"
	KEY_LAT_MIN            = "LATITUDE (MIN.)"
	KEY_LONG_MIN           = "LONGITUDE (MIN.)"
	KEY_LAT_MAX            = "LATITUDE (MAX.)"
	KEY_LONG_MAX           = "LONGITUDE (MAX.)"
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
	rows := make([][]string, 0, len(samples))
	// column headers
	rows = append(rows, []string{"YEAR", "CITATION", "SAMPLE NAME", "UNIQUE_ID", "LOCATION", "ELEVATION (MIN.)", "ELEVATION (MAX.)", "SAMPLING TECHNIQUE", "DRILLING DEPTH (MIN.)", "DRILLING DEPTH (MAX.)", "LAND/SEA (SAMPLING)", "ROCK TYPE", "ROCK NAME", "ROCK TEXTURE", "SAMPLE COMMENT", "AGE (MIN.)", "AGE (MAX.)", "GEOLOGICAL AGE", "GEOLOGICAL AGE PREFIX", "ERUPTION DATE", "ALTERATION", "ALTERATION TYPE", "TYPE OF MATERIAL", "MINERAL / COMPONENT", "CRYSTAL", "RIM / CORE (MINERAL GRAINS)", "INCLUSION TYPE", "MINERAL (INCLUSION)", "RIM / CORE (INCLUSION)", "HOSTMINERAL (INCLUSION)", "LATITUDE (MIN.)", "LONGITUDE (MIN.)", "LATITUDE (MAX.)", "LONGITUDE (MAX.)", "SIO2(WT%)", "TIO2(WT%)", "AL2O3(WT%)", "CR2O3(WT%)", "FE2O3T(WT%)", "FE2O3(WT%)", "FEOT(WT%)", "FEO(WT%)", "CAO(WT%)", "MGO(WT%)", "MNO(WT%)", "BAO(WT%)", "K2O(WT%)", "NA2O(WT%)", "P2O5(WT%)", "H2OP(WT%)", "H2OM(WT%)", "CL(WT%)", "SO3(WT%)", "S(WT%)", "LOI(WT%)", "SR(WT%)", "LI(PPM)", "BE(PPM)", "B(PPM)", "F(PPM)", "NA(PPM)", "MG(PPM)", "AL(PPM)", "SI(PPM)", "P(PPM)", "S(PPM)", "CL(PPM)", "K(PPM)", "CA(PPM)", "SC(PPM)", "TI(PPM)", "V(PPM)", "CR(PPM)", "MN(PPM)", "FE(PPM)", "CO(PPM)", "NI(PPM)", "CU(PPM)", "ZN(PPM)", "GA(PPM)", "GE(PPM)", "AS(PPM)", "SE(PPM)", "BR(PPM)", "RB(PPM)", "SR(PPM)", "Y(PPM)", "ZR(PPM)", "NB(PPM)", "MO(PPM)", "PD(PPM)", "PD(PPB)", "AG(PPM)", "CD(PPM)", "IN(PPM)", "SN(PPM)", "SB(PPM)", "TE(PPM)", "CS(PPM)", "BA(PPM)", "LA(PPM)", "CE(PPM)", "PR(PPM)", "ND(PPM)", "SM(PPM)", "EU(PPM)", "GD(PPM)", "TB(PPM)", "DY(PPM)", "HO(PPM)", "ER(PPM)", "TM(PPM)", "YB(PPM)", "LU(PPM)", "HF(PPM)", "TA(PPM)", "W(PPM)", "RE(PPM)", "PT(PPM)", "PT(PPB)", "AU(PPM)", "AU(PPB)", "HG(PPM)", "TL(PPM)", "PB(PPM)", "PB204(PPM)", "PB206(PPM)", "PB207(PPM)", "PB208(PPM)", "BI(PPM)", "TH(PPM)", "U(PPM)", "U238(PPM)", "ND143_ND144", "EPSILON_ND_INI", "SM147_ND144", "SR87_SR86", "RB87_SR86", "PB206_PB204", "PB207_PB204", "PB207_PB206", "PB208_PB204", "PB208_PB206", "HF176_HF177", "PB206_U238", "PB207_U235", "D13C(VS VPDB)", "D18O(VS VPDB)", "D18O(VS VSMOW)", "D34S(VS VCDT)", "ALBITE(MOL%)", "ALMANDINE(MOL%)", "ANORTHITE(MOL%)", "ENSTATITE(MOL%)", "FERROSILITE(MOL%)", "GROSSULAR(MOL%)", "ORTHOCLASE(MOL%)", "PYROPE(MOL%)", "SPESSARTINE(MOL%)", "WOLLASTONITE(MOL%)", "AGE_PB207_PB206(MA)", "AGE_PB207_U235(MA)", "AGE_PB206_U238(MA)"})
	for _, sample := range samples {
		rowMap := map[string]string{"YEAR": "", "CITATION": "", "SAMPLE NAME": "", "UNIQUE_ID": "", "LOCATION": "", "ELEVATION (MIN.)": "", "ELEVATION (MAX.)": "", "SAMPLING TECHNIQUE": "", "DRILLING DEPTH (MIN.)": "", "DRILLING DEPTH (MAX.)": "", "LAND/SEA (SAMPLING)": "", "ROCK TYPE": "", "ROCK NAME": "", "ROCK TEXTURE": "", "SAMPLE COMMENT": "", "AGE (MIN.)": "", "AGE (MAX.)": "", "GEOLOGICAL AGE": "", "GEOLOGICAL AGE PREFIX": "", "ERUPTION DATE": "", "ALTERATION": "", "ALTERATION TYPE": "", "TYPE OF MATERIAL": "", "MINERAL / COMPONENT": "", "CRYSTAL": "", "RIM / CORE (MINERAL GRAINS)": "", "INCLUSION TYPE": "", "MINERAL (INCLUSION)": "", "RIM / CORE (INCLUSION)": "", "HOSTMINERAL (INCLUSION)": "", "LATITUDE (MIN.)": "", "LONGITUDE (MIN.)": "", "LATITUDE (MAX.)": "", "LONGITUDE (MAX.)": "", "SIO2(WT%)": "", "TIO2(WT%)": "", "AL2O3(WT%)": "", "CR2O3(WT%)": "", "FE2O3T(WT%)": "", "FE2O3(WT%)": "", "FEOT(WT%)": "", "FEO(WT%)": "", "CAO(WT%)": "", "MGO(WT%)": "", "MNO(WT%)": "", "BAO(WT%)": "", "K2O(WT%)": "", "NA2O(WT%)": "", "P2O5(WT%)": "", "H2OP(WT%)": "", "H2OM(WT%)": "", "CL(WT%)": "", "SO3(WT%)": "", "S(WT%)": "", "LOI(WT%)": "", "SR(WT%)": "", "LI(PPM)": "", "BE(PPM)": "", "B(PPM)": "", "F(PPM)": "", "NA(PPM)": "", "MG(PPM)": "", "AL(PPM)": "", "SI(PPM)": "", "P(PPM)": "", "S(PPM)": "", "CL(PPM)": "", "K(PPM)": "", "CA(PPM)": "", "SC(PPM)": "", "TI(PPM)": "", "V(PPM)": "", "CR(PPM)": "", "MN(PPM)": "", "FE(PPM)": "", "CO(PPM)": "", "NI(PPM)": "", "CU(PPM)": "", "ZN(PPM)": "", "GA(PPM)": "", "GE(PPM)": "", "AS(PPM)": "", "SE(PPM)": "", "BR(PPM)": "", "RB(PPM)": "", "SR(PPM)": "", "Y(PPM)": "", "ZR(PPM)": "", "NB(PPM)": "", "MO(PPM)": "", "PD(PPM)": "", "PD(PPB)": "", "AG(PPM)": "", "CD(PPM)": "", "IN(PPM)": "", "SN(PPM)": "", "SB(PPM)": "", "TE(PPM)": "", "CS(PPM)": "", "BA(PPM)": "", "LA(PPM)": "", "CE(PPM)": "", "PR(PPM)": "", "ND(PPM)": "", "SM(PPM)": "", "EU(PPM)": "", "GD(PPM)": "", "TB(PPM)": "", "DY(PPM)": "", "HO(PPM)": "", "ER(PPM)": "", "TM(PPM)": "", "YB(PPM)": "", "LU(PPM)": "", "HF(PPM)": "", "TA(PPM)": "", "W(PPM)": "", "RE(PPM)": "", "PT(PPM)": "", "PT(PPB)": "", "AU(PPM)": "", "AU(PPB)": "", "HG(PPM)": "", "TL(PPM)": "", "PB(PPM)": "", "PB204(PPM)": "", "PB206(PPM)": "", "PB207(PPM)": "", "PB208(PPM)": "", "BI(PPM)": "", "TH(PPM)": "", "U(PPM)": "", "U238(PPM)": "", "ND143_ND144": "", "EPSILON_ND_INI": "", "SM147_ND144": "", "SR87_SR86": "", "RB87_SR86": "", "PB206_PB204": "", "PB207_PB204": "", "PB207_PB206": "", "PB208_PB204": "", "PB208_PB206": "", "HF176_HF177": "", "PB206_U238": "", "PB207_U235": "", "D13C(VS VPDB)": "", "D18O(VS VPDB)": "", "D18O(VS VSMOW)": "", "D34S(VS VCDT)": "", "ALBITE(MOL%)": "", "ALMANDINE(MOL%)": "", "ANORTHITE(MOL%)": "", "ENSTATITE(MOL%)": "", "FERROSILITE(MOL%)": "", "GROSSULAR(MOL%)": "", "ORTHOCLASE(MOL%)": "", "PYROPE(MOL%)": "", "SPESSARTINE(MOL%)": "", "WOLLASTONITE(MOL%)": "", "AGE_PB207_PB206(MA)": "", "AGE_PB207_U235(MA)": "", "AGE_PB206_U238(MA)": ""}
		// citation metadata
		if len(sample.References) > 0 {
			ref := sample.References[0]
			rowMap[KEY_YEAR] = getInt(ref.Publicationyear)
			rowMap[KEY_CITATION] = fmt.Sprintf("[%v]%v", getString(ref.Externalidentifier), getString(ref.Title))
		}
		// sample data
		rowMap[KEY_SAMPLENAME] = getString(sample.SampleName)
		rowMap[KEY_UNIQUE_ID] = getString(sample.UniqueID)
		rowMap[KEY_LOCATION] = join(sample.LocationNames, "/")
		rowMap[KEY_ELEVATION_MIN] = getString(sample.ElevationMin)
		rowMap[KEY_ELEVATION_MAX] = getString(sample.ElevationMax)
		rowMap[KEY_SAMPLING_TECHNIQUE] = getString(sample.SamplingTechnique)
		rowMap[KEY_DRILLDEPTH_MIN] = getString(sample.DrillDepthMin)
		rowMap[KEY_DRILLDEPTH_MAX] = getString(sample.DrillDepthMax)
		rowMap[KEY_LANDORSEA] = getString(sample.LandOrSea)
		rowMap[KEY_ROCKTYPE] = parseTaxonomicclassifier(sample.RockTypes)
		rowMap[KEY_ROCKNAME] = parseTaxonomicclassifier(sample.RockClasses)
		rowMap[KEY_ROCKTEXTURE] = join(sample.RockTextures, ";")
		rowMap[KEY_SAMPLECOMMENT] = join(sample.Comments, ";")
		rowMap[KEY_AGE_MIN] = getInt(sample.AgeMin)
		rowMap[KEY_AGE_MAX] = getInt(sample.AgeMax)
		rowMap[KEY_GEO_AGE] = getString(sample.GeologicalAge)
		rowMap[KEY_AGE_PREFIX] = getString(sample.GeologicalAgePrefix)
		rowMap[KEY_ERUPTION_DATE] = getString(sample.EruptionDate)
		rowMap[KEY_ALTERATION] = getString(sample.Alteration)
		rowMap[KEY_ALTERATION_TYPE] = getString(sample.AlterationType)

		// batch data
		for _, batch := range sample.BatchData {
			rowMap[KEY_MATERIAL_TYPE] = getString(batch.Material)
			rowMap[KEY_MINERAL] = parseTaxonomicclassifier(batch.Minerals)
			rowMap[KEY_CRYSTAL] = getString(batch.Crystal)
			rowMap[KEY_RIMORCORE] = getString(batch.RimOrCoreMineral)
			rowMap[KEY_INCLUSIONTYPE] = join(batch.InclusionTypes, ";")
			rowMap[KEY_INCLUSION_MINERAL] = parseTaxonomicclassifier(batch.InclusionMinerals)
			rowMap[KEY_RIMORCORE_INC] = getString(batch.RimOrCoreInclusion)
			rowMap[KEY_HOST_MINERAL] = parseTaxonomicclassifier(batch.HostMinerals)
			rowMap[KEY_LAT_MIN] = getString(sample.LatitudeMin)
			rowMap[KEY_LONG_MIN] = getString(sample.LongitudeMin)
			rowMap[KEY_LAT_MAX] = getString(sample.LatitudeMax)
			rowMap[KEY_LONG_MAX] = getString(sample.LongitudeMax)
			// add result data
			for _, result := range batch.Results {
				itemName := getString(result.ItemName)
				value := getFloat64(result.Value)
				unit := getString(result.Unit)
				if itemName == "" || value == "" {
					continue
				}
				key := itemName
				if unit != "" {
					key = fmt.Sprintf("%s(%s)", itemName, unit)
				}
				if _, ok := rowMap[key]; !ok {
					fmt.Printf("No column for key: %s\n", key)
				}
				rowMap[key] = value
			}
		}
		// every row must have the same order (as defined by the header row), especially for the chemical items - so we lookup each column name in the map
		row := make([]string, 0, len(rowMap))
		for _, key := range rows[0] {
			row = append(row, fmt.Sprintf("\"%s\"", rowMap[key]))
		}
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

// parseTaxonomicclassifier takes a slice of Taxonomicclassifiers and returns the labels ";"-separated as a string
func parseTaxonomicclassifier(types []*model.FullDataTaxonomicClassifier) string {
	rt := ""
	for i, t := range types {
		if i > 0 {
			rt += ";"
		}
		rt += getString(t.Label)
	}
	return rt
}
