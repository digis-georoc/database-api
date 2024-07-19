package download

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
)

const (
	CSV  = "csv"
	XLSX = "xlsx"

	// csv column keys
	KEY_YEAR               = "YEAR"
	KEY_DOI                = "DOI"
	KEY_CITATION           = "CITATION"
	KEY_AUTHORS            = "AUTHORS"
	KEY_CITATION_METADATA  = "CITATION METADATA"
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
	rows := makeRows(samples)
	csv := ""
	for i, row := range rows {
		csv += strings.Join(row, ",")
		if i < len(rows)-1 {
			csv += "\n"
		}
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
	file := excelize.NewFile()
	defer file.Close()
	rows := makeRows(samples)
	for i, row := range rows {
		for j, val := range row {
			cell, err := excelize.CoordinatesToCellName(j+1, i+1) // cells are 1-based
			if err != nil {
				return nil, fmt.Errorf("Can not convert coordinates (%d, %d) to cellName: %s", i+1, j+1, err.Error())
			}
			err = file.SetCellValue("Sheet1", cell, val)
			if err != nil {
				return nil, fmt.Errorf("Can not set cell value (%s) to %+v", cell, val)
			}
		}
	}
	buf, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("Can not write xlsx data to buffer: %s", err.Error())
	}
	return buf.Bytes(), nil
}

// join implements strings.Join() for type []*string
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

// makeRows formats a slice of FullData models as table rows
func makeRows(samples []model.FullData) [][]string {
	rows := make([][]string, 0, len(samples))
	// column headers
	headerRow := []string{KEY_YEAR, KEY_DOI, KEY_CITATION, KEY_CITATION_METADATA, KEY_AUTHORS, KEY_SAMPLENAME, KEY_UNIQUE_ID, KEY_LOCATION, KEY_ELEVATION_MIN, KEY_ELEVATION_MAX, KEY_SAMPLING_TECHNIQUE, KEY_DRILLDEPTH_MIN, KEY_DRILLDEPTH_MAX, KEY_LANDORSEA, KEY_ROCKTYPE, KEY_ROCKNAME, KEY_ROCKTEXTURE, KEY_SAMPLECOMMENT, KEY_AGE_MIN, KEY_AGE_MAX, KEY_GEO_AGE, KEY_AGE_PREFIX, KEY_ERUPTION_DATE, KEY_ALTERATION, KEY_ALTERATION_TYPE, KEY_MATERIAL_TYPE, KEY_MINERAL, KEY_CRYSTAL, KEY_RIMORCORE, KEY_INCLUSIONTYPE, KEY_INCLUSION_MINERAL, KEY_RIMORCORE_INC, KEY_HOST_MINERAL, KEY_LAT_MIN, KEY_LONG_MIN, KEY_LAT_MAX, KEY_LONG_MAX}
	itemsMap := map[string]map[string]bool{}
	rowMaps := make([]map[string]string, 0, len(samples))
	for _, sample := range samples {
		rowMap := map[string]string{}
		// citation metadata
		if len(sample.References) > 0 {
			ref := sample.References[0]
			rowMap[KEY_YEAR] = getInt(ref.Publicationyear)
			rowMap[KEY_DOI] = getString(ref.Externalidentifier)
			rowMap[KEY_CITATION] = getString(ref.Title)
			authors := ""
			for i, author := range ref.Authors {
				if i > 0 {
					authors += ";"
				}
				authors += fmt.Sprintf("%s %s", getString(author.PersonLastName), getString(author.PersonFirstName))
			}
			rowMap[KEY_AUTHORS] = authors
			rowMap[KEY_CITATION_METADATA] = fmt.Sprintf("Journal:%s;Volume:%s;Issue:%s;BookTitle:%s;FirstPage:%s;LastPage:%s", getString(ref.Journal), getString(ref.Volume), getString(ref.Issue), getString(ref.BookTitle), getString(ref.FirstPage), getString(ref.LastPage))
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
		rowMap[KEY_AGE_MIN] = getFloat64(sample.AgeMin)
		rowMap[KEY_AGE_MAX] = getFloat64(sample.AgeMax)
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
				itemType := getString(result.ItemGroup)
				value := getFloat64(result.Value)
				unit := getString(result.Unit)
				method := getString(result.Method)
				if itemName == "" || value == "" {
					continue
				}
				key := itemName
				if unit != "" {
					key += fmt.Sprintf("(%s)", unit)
				}
				if method != "" {
					key += fmt.Sprintf("[%s]", method)
				}
				if typeMap := itemsMap[itemType]; typeMap == nil {
					itemsMap[itemType] = map[string]bool{}
				}
				itemsMap[itemType][key] = true
				rowMap[key] = value
			}
		}
		rowMaps = append(rowMaps, rowMap)
	}
	// append sorted items to headerRow
	for _, itemType := range []string{"mj", "ree", "te", "rg", "ir", "is", "us", "em", "age"} {
		items := getKeySlice(itemsMap[itemType])
		sort.SliceStable(items, func(i, j int) bool { return items[i] < items[j] })
		headerRow = append(headerRow, items...)
	}
	rows = append(rows, headerRow)
	for _, rowMap := range rowMaps {
		// every row must have the same order (as defined by the header row), especially for the chemical items - so we lookup each column name in the map
		row := make([]string, 0, len(headerRow))
		for _, key := range headerRow {
			metaData := rowMap[key]
			row = append(row, fmt.Sprintf("\"%s\"", metaData))
		}
		rows = append(rows, row)
	}
	return rows
}

func getKeySlice(m map[string]bool) []string {
	s := make([]string, len(m))
	i := 0
	for k := range m {
		s[i] = k
		i++
	}
	return s
}
