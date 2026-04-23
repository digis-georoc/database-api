package model

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
