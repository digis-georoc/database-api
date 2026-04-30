package model

import (
	"encoding/json"
	"fmt"
)

type SearchIndexPage struct {
	TotalHits    int              `json:"totalHits"`
	Documents    []map[string]any `json:"documents"`
	Aggregations map[string]any   `json:"aggregations"`
}

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type ClusterAggregations struct {
	Clustering ClusteringAgg `json:"clustering"`
}

type ClusteringAgg struct {
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	Key      string      `json:"key"`
	DocCount int         `json:"doc_count"`
	Centroid CentroidAgg `json:"centroid"`
}

type CentroidAgg struct {
	Count    int         `json:"count"`
	Location AggLocation `json:"location"`
}

type AggLocation struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func ParseToSampleByFiltersData(doc map[string]any) (*SampleByFiltersData, error) {
	// marhsal into fullData first
	bytes, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	fullData := FullData{}
	err = json.Unmarshal(bytes, &fullData)
	if err != nil {
		return nil, err
	}
	rc := []*string{}
	for _, t := range fullData.RockClasses {
		rc = append(rc, &t.Value)
	}
	rt := []*string{}
	for _, t := range fullData.RockTypes {
		rc = append(rc, &t.Value)
	}
	authors := []*SampleAuthor{}
	var pubYear *int
	var extID *string
	if len(fullData.References) > 0 {
		pubYear = fullData.References[0].Publicationyear
		extID = fullData.References[0].Externalidentifier
		for _, a := range fullData.References[0].Authors {
			authors = append(authors, &SampleAuthor{
				FirstName: *a.PersonFirstName,
				LastName:  *a.PersonLastName,
				Order:     *a.AuthorOrder,
			})
		}
	}
	mins := []*string{}
	hostMins := []*string{}
	incMins := []*string{}
	for _, b := range fullData.BatchData {
		for _, m := range b.Minerals {
			mins = append(mins, &m.Value)
		}
		for _, hm := range b.HostMinerals {
			hostMins = append(hostMins, &hm.Value)
		}
		for _, im := range b.InclusionMinerals {
			incMins = append(incMins, &im.Value)
		}
	}
	sample := SampleByFiltersData{
		SampleID:           fullData.SampleID,
		GeologicalAges:     []*string{fullData.GeologicalAge},
		RockClasses:        rc,
		RockTypes:          rt,
		PublicationYear:    pubYear,
		ExternalIdentifier: extID,
		Authors:            authors,
		Minerals:           mins,
		HostMinerals:       hostMins,
		InclusionMinerals:  incMins,
		GeologicalSettings: []*string{fullData.TectonicSetting},
	}
	if fullData.AgeMin != nil {
		ageMin := fmt.Sprintf("%f", *fullData.AgeMin)
		sample.GeologicalAgesMin = []*string{&ageMin}
	}
	if fullData.AgeMax != nil {
		ageMax := fmt.Sprintf("%f", *fullData.AgeMax)
		sample.GeologicalAgesMax = []*string{&ageMax}
	}
	if fullData.SampleName != nil {
		sample.SampleName = *fullData.SampleName
	}
	if fullData.Latitude != nil {
		sample.Latitude = float64(*fullData.Latitude)
	}
	if fullData.Longitude != nil {
		sample.Longitude = float64(*fullData.Longitude)
	}
	return &sample, nil
}
