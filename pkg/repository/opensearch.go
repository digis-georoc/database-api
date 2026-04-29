package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/defensestation/osquery/v2"
	log "github.com/sirupsen/logrus"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/geometry"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/model"
	"gitlab.gwdg.de/fe/digis/database-api/pkg/secretstore"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

const (
	INDEX_NAME = "digis-test-index-multishard"

	MAX_OS_PAGESIZE = 10000

	KEY_CLUSTERING = "clustering"
	KEY_CENTROID   = "centroid"

	FILTER_CHEMISTRY = "chemistry"
	FILTER_POLYGON   = "polygon"
	FILTER_BBOX      = "bbox"

	FIELD_GEOPOINT  = "geo_point"
	FIELD_VALUE     = "value"
	FIELD_ITEMGROUP = "itemGroup"
	FIELD_ITEMNAME  = "itemName"

	PREFIX_IN = "IN"
	PREFIX_EQ = "EQ"
	PREFIX_LT = "LT"
	PREFIXGT  = "GT"

	DELIM = ","
)

var (
	MINIMALFIELDS = []string{"sampleID", "latitude", "longitude"}
	SEARCH_FIELDS = []string{"sampleID", "sampleName", "latitude", "longitude", "batchData.batchID", "references.publicationYear", "references.externalIdentifier", "references.authors", "batchData.minerals", "batchData.hostMinerals", "batchData.inclusionMinerals", "rockClasses", "rockTypes", "batchData.inclusionTypes", "geologicalSettings", "geologicalAge", "ageMin", "ageMax"}
)

type OSClient struct {
	client *opensearchapi.Client
}

func NewOSClient(secStore secretstore.SecretStore) (*OSClient, error) {
	client := OSClient{}
	params, err := getOSConnectionParams(secStore)
	if err != nil {
		return nil, err
	}
	err = client.connect(params.Host, params.Username, params.Password)
	return &client, err
}

func getOSConnectionParams(secStore secretstore.SecretStore) (*ConnectionParams, error) {
	host, err := secStore.GetSecret("OS_HOST")
	if err != nil {
		return nil, fmt.Errorf("can not load secret OS_HOST")
	}
	username, err := secStore.GetSecret("OS_USER")
	if err != nil {
		return nil, fmt.Errorf("can not load secret OS_USER")
	}
	password, err := secStore.GetSecret("OS_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("can not load secret OS_PASSWORD")
	}
	return &ConnectionParams{
		Host:     host,
		Username: username,
		Password: password,
	}, nil
}

func (os *OSClient) connect(host string, user string, password string) error {
	// init opensearch client
	client, err := opensearchapi.NewClient(
		opensearchapi.Config{
			Client: opensearch.Config{
				// InsecureSkipVerify: true, // For testing only. Use certificate for validation.
				Addresses:         []string{host},
				Username:          user,
				Password:          password,
				EnableDebugLogger: true,
			},
		},
	)
	if err != nil {
		return err
	}

	os.client = client

	infoResp, err := client.Info(context.Background(), nil)
	if err != nil {
		return err
	}
	log.Infof("Connected to Cluster:\n  Cluster Name: %s\n  Cluster UUID: %s\n  Version Number: %s\n", infoResp.ClusterName, infoResp.ClusterUUID, infoResp.Version.Number)
	return nil
}

func (os *OSClient) QueryClustered(includeFields []string, filters map[string]string, zoomLevel int) (model.ClusterResponse, error) {
	clusterResp := model.ClusterResponse{}
	query, err := buildQuery(filters)
	if err != nil {
		return clusterResp, fmt.Errorf("can not build query from filters: %w", err)
	}
	params := &opensearchapi.SearchParams{
		TrackTotalHits: true,
		Source:         false,
	}
	baseQuery := osquery.Search().Size(0).Query(query).Sort(osquery.FieldSort("sampleID").Order(osquery.OrderAsc)).Aggs(osquery.CustomAgg(KEY_CLUSTERING, map[string]any{"geotile_grid": map[string]any{"field": FIELD_GEOPOINT, "precision": zoomLevel}, "aggs": map[string]any{KEY_CENTROID: map[string]any{"geo_centroid": map[string]any{"field": FIELD_GEOPOINT}}}}))
	q, _ := baseQuery.MarshalJSON()
	fmt.Println(string(q))
	searchResponse, err := runQuery(os.client.Client, *baseQuery, INDEX_NAME, params)
	if err != nil {
		return clusterResp, fmt.Errorf("can not run query: %w", err)
	}

	clusterResp, err = parseClusterResponse(searchResponse)
	if err != nil {
		return clusterResp, fmt.Errorf("can not parse results: %w", err)
	}
	return clusterResp, nil
}

func (os *OSClient) QuerySortSearchAfterStream(ctx context.Context, includeFields []string, filters map[string]string, size int, resultChan chan model.SearchIndexPage) {
	defer close(resultChan)
	// TODO: use PIT or not? If yes, need to remove index name from url as it will be taken from PIT
	query, err := buildQuery(filters)
	if err != nil {
		log.Errorf("can not build query from filters: %s", err.Error())
		return
	}
	params := &opensearchapi.SearchParams{
		TrackTotalHits: true,
		Source:         true,
		SourceIncludes: includeFields,
		// DocvalueFields: includeFields,
	}

	if size <= 0 || size > MAX_OS_PAGESIZE {
		size = MAX_OS_PAGESIZE
	}
	baseQuery := osquery.Search().Size(uint64(size)).Query(query).Sort(osquery.FieldSort("sampleID").Order(osquery.OrderAsc))
	q, _ := baseQuery.MarshalJSON()
	fmt.Println(string(q))
	searchResponse, err := runQuery(os.client.Client, *baseQuery, INDEX_NAME, params)
	if err != nil {
		log.Errorf("can not run query: %s", err.Error())
		return
	}

	resultCount := len(searchResponse.Hits.Hits)
	page, err := parseIndexPage(searchResponse.Hits, searchResponse.Aggregations)
	if err != nil {
		log.Errorf("can not parse results: %s", err.Error())
		return
	}
	resultChan <- page

	// if more hits than size - follow up requests with search_after: [LAST_SORT_VALUE]
	for searchResponse.Hits.Total.Value > resultCount {
		// cancel retrieving hits if request context is canceled
		select {
		case <-ctx.Done():
			log.Info("Request context is done - stopping search")
			return
		case <-time.After(1 * time.Millisecond):
			// continue search
		}
		numReturned := len(searchResponse.Hits.Hits)
		if numReturned == 0 {
			log.Warnf("search returned no results. Status: %d", searchResponse.Inspect().Response.StatusCode)
			return
		}
		lastSortVal := searchResponse.Hits.Hits[numReturned-1].Sort
		pageQuery := osquery.Search().Size(uint64(size)).Query(query).SearchAfter(lastSortVal...).Sort(osquery.FieldSort("sampleID").Order(osquery.OrderAsc))
		searchResponse, err = runQuery(os.client.Client, *pageQuery, INDEX_NAME, params)
		if err != nil {
			log.Errorf("can not run subsequent query: %s", err.Error())
			return
		}

		resultCount += len(searchResponse.Hits.Hits)
		page, err := parseIndexPage(searchResponse.Hits, searchResponse.Aggregations)
		if err != nil {
			log.Errorf("can not parse subsequent results: %s", err.Error())
			return
		}
		resultChan <- page
	}
}

// buildQuery constructs a osquery.BoolQuery from given filters
func buildQuery(filters map[string]string) (*osquery.BoolQuery, error) {
	osFilters := []osquery.Mappable{}
	for k, v := range filters {
		nestedPath := getNested(k)
		k = translateKey(k)
		var f []osquery.Mappable
		if len(nestedPath) > 0 {
			// construct fieldName from original key and nestedPath
			fieldName := fmt.Sprintf("%s.%s", nestedPath, k)
			// parse custom chemistry query DSL
			if k == FILTER_CHEMISTRY {
				chemFilters, err := model.ParseChemQuery(v)
				if err != nil {
					return nil, err
				}
				// add a filter conjunction for each analyte
				for _, chemFilter := range chemFilters.Expressions {
					minValue, err := strconv.ParseFloat(chemFilter.MinValue, 64)
					if err != nil {
						return nil, err
					}
					maxValue, err := strconv.ParseFloat(chemFilter.MaxValue, 64)
					if err != nil {
						return nil, err
					}
					rng := osquery.Range(fmt.Sprintf("%s.%s", nestedPath, FIELD_VALUE))
					if !math.IsNaN(minValue) {
						rng = rng.Gte(minValue)
					}
					if !math.IsNaN(maxValue) {
						rng = rng.Lte(maxValue)
					}
					filters := []osquery.Mappable{
						rng,
					}
					if chemFilter.Type != "" {
						filters = append(filters, osquery.Term(fmt.Sprintf("%s.%s", nestedPath, FIELD_ITEMGROUP), chemFilter.Type))
					}
					if chemFilter.Element != "" {
						filters = append(filters, osquery.Term(fmt.Sprintf("%s.%s", nestedPath, FIELD_ITEMNAME), chemFilter.Element))
					}
					q := osquery.Bool().Filter(filters...)
					f = append(f, osquery.Nested(nestedPath, q))
				}
			} else {
				f = append(f, osquery.Nested(nestedPath, dslToFilterQuery(fieldName, v)))
			}
		} else {
			switch k {
			case FILTER_POLYGON:
				polygon, err := geometry.ParsePointArray(v)
				if err != nil {
					return nil, err
				}
				should := []osquery.Mappable{}
				// imitate a wrap around +/-180 meridian by using the original polygon and a copy moved by +/-180 depending on the crossed boundary
				partial1, partial2, err := geometry.WrapPolygonLon(polygon)
				if err != nil {
					return nil, err
				}
				points1 := []model.GeoPoint{}
				for _, coords := range partial1 {
					points1 = append(points1, model.GeoPoint{Lat: coords.Y, Lon: coords.X})
				}
				// do a geopolygon filter
				should = append(should, osquery.CustomQuery(map[string]any{"geo_polygon": map[string]any{FIELD_GEOPOINT: map[string]any{"points": points1}}}))
				points2 := []model.GeoPoint{}
				for _, coords := range partial2 {
					points2 = append(points1, model.GeoPoint{Lat: coords.Y, Lon: coords.X})
				}
				// do a geopolygon filter
				should = append(should, osquery.CustomQuery(map[string]any{"geo_polygon": map[string]any{FIELD_GEOPOINT: map[string]any{"points": points2}}}))
				f = append(f, osquery.Bool().MinimumShouldMatch(1).Should(should...))
			case FILTER_BBOX:
				bbox, err := geometry.ParsePointArray(v)
				if err != nil {
					return nil, err
				}
				should := []osquery.Mappable{}
				// imitate a wrap around +/-180 meridian by using the original polygon and a copy moved by +/-180 depending on the crossed boundary
				partial1, partial2, err := geometry.WrapPolygonLon(bbox)
				if err != nil {
					return nil, err
				}
				should = append(should, osquery.CustomQuery(map[string]any{"geo_bounding_box": map[string]any{FIELD_GEOPOINT: map[string]any{"top_right": map[string]any{"lon": partial1[2].X, "lat": partial1[2].Y}, "bottom_left": map[string]any{"lon": partial1[0].X, "lat": partial1[0].Y}}}}))
				should = append(should, osquery.CustomQuery(map[string]any{"geo_bounding_box": map[string]any{FIELD_GEOPOINT: map[string]any{"top_right": map[string]any{"lon": partial2[2].X, "lat": partial2[2].Y}, "bottom_left": map[string]any{"lon": partial2[0].X, "lat": partial2[0].Y}}}}))
				f = append(f, osquery.Bool().MinimumShouldMatch(1).Should(should...))
			default:
				// do a normal term filter
				f = append(f, dslToFilterQuery(k, v))
			}
		}
		osFilters = append(osFilters, f...)
	}
	query := osquery.Bool().Filter(osFilters...)
	return query, nil
}

// dslToFilterQuery takes a field name and a value string and parses the custom query dsl to return search index queries as osquery.Mappable objects
func dslToFilterQuery(field string, v string) osquery.Mappable {
	var fq osquery.Mappable
	operator, value, found := strings.Cut(v, ":")
	if !found {
		// if no operator is specified, "EQ" is assumed as default
		operator = PREFIX_EQ
		// retain original value
		value = v
	}
	operator = strings.ToUpper(operator)
	switch operator {
	case PREFIX_IN:
		generified := []any{}
		for s := range strings.SplitSeq(value, DELIM) {
			generified = append(generified, s)
		}
		fq = osquery.Terms(field, generified...)
	case PREFIX_EQ:
		fq = osquery.Term(field, value)
	default:
		fq = osquery.Term(field, value)
	}
	return fq
}

// translateKey returns the document key corresponding to a search term
func translateKey(k string) string {
	switch k {
	case "doi":
		return "externalIdentifier"
	case "setting":
		return "tectonicSetting"
	case "geoage":
		return "geologicalAge"
	case "geoageprefix":
		return "geologicalAgePrefix"
	case "mineral", "inclusionmaterial", "hostmaterial", "rocktype", "rockclass":
		return "value"
	case "standards":
		return "standardName"
	}
	return k
}

// getNested returns for a given key in the filters map, the nested field in the documents it belongs to; or empty string if the key is top level
func getNested(key string) string {
	switch key {
	case "batchName", "batchID", "crystal", "inclusionTypes", "material", "rimOrCoreInclusion", "rimOrCoreMineral", "specimenMedium":
		return "batchData"
	case "mineral":
		return "batchData.minerals"
	case "inclusionmaterial":
		return "batchData.inclusionMinerals"
	case "hostmaterial":
		return "batchData.hostMinerals"
	case "itemName", "itemGroup", "medium", "method", "unit", "value", "valueCount", FILTER_CHEMISTRY:
		return "batchData.results"
	case "standards":
		return "batchData.results.standards"
	case "rocktype":
		return "rockTypes"
	case "rockclass":
		return "rockClasses"
	case "citationID", "title", "publisher", "publicationYear", "citationLink", "journal", "volume", "issue", "firstPage", "lastPage", "bookTitle", "editors", "externalIdentifier":
		return "references"
	case "personID", "firstName", "lastName", "order":
		return "references.authors"
	}
	return ""
}

func parseIndexPage(hits opensearchapi.SearchHits, aggregations json.RawMessage) (model.SearchIndexPage, error) {
	page := model.SearchIndexPage{
		TotalHits: hits.Total.Value,
	}
	if len(aggregations) > 0 {
		aggs := map[string]any{}
		err := json.Unmarshal(aggregations, &aggs)
		if err != nil {
			return page, err
		}
		page.Aggregations = aggs
	}
	results := make([]map[string]any, 0, len(hits.Hits))
	for _, hit := range hits.Hits {
		doc := map[string]any{}
		err := json.Unmarshal(hit.Source, &doc)
		if err != nil {
			return page, err
		}
		results = append(results, doc)
	}
	page.Documents = results
	return page, nil
}

func parseClusterResponse(resp *opensearchapi.SearchResp) (model.ClusterResponse, error) {
	clusterResp := model.ClusterResponse{}
	if len(resp.Aggregations) == 0 {
		return clusterResp, nil
	}
	aggs := model.ClusterAggregations{}
	err := json.Unmarshal(resp.Aggregations, &aggs)
	if err != nil {
		return clusterResp, err
	}
	for i, b := range aggs.Clustering.Buckets {
		clusterResp.Clusters = append(clusterResp.Clusters, model.GeoJSONCluster{
			ClusterID: i,
			Centroid: model.GeoJSONFeature{
				Type: model.GEOJSONTYPE_FEATURE,
				ID:   b.Key,
				Geometry: model.Geometry{
					Type:        model.GEOJSON_GEOMETRY_POINT,
					Coordinates: []any{b.Centroid.Location.Lon, b.Centroid.Location.Lat},
				},
				Properties: map[string]any{
					"clusterID":   i,
					"clusterSize": b.DocCount,
				},
			},
			// add dummy geometry
			ConvexHull: model.GeoJSONFeature{
				Type: model.GEOJSONTYPE_FEATURE,
				ID:   b.Key,
				Geometry: model.Geometry{
					Type:        model.GEOJSON_GEOMETRY_POINT,
					Coordinates: []any{b.Centroid.Location.Lon, b.Centroid.Location.Lat},
				},
			},
		})
	}
	// set to non-null
	clusterResp.Points = []model.GeoJSONFeature{}
	return clusterResp, nil
}

// collectResults marshals result data and aggregates it into the result array
func collectResults(hits []opensearchapi.SearchHit, results []model.FullData) ([]model.FullData, error) {
	start := time.Now()
	for _, hit := range hits {
		sample := model.FullData{}
		err := json.Unmarshal(hit.Source, &sample)
		if err != nil {
			return nil, fmt.Errorf("can not unmarshal hit result: %w", err)
		}
		results = append(results, sample)
	}
	took := time.Since(start).String()
	fmt.Printf("Parsed %d samples in %s\n", len(hits), took)
	return results, nil
}

// runQuery executes the given search query with the opensearch client and logs the time
func runQuery(client *opensearch.Client, searchQuery osquery.SearchRequest, indexName string, params *opensearchapi.SearchParams) (*opensearchapi.SearchResp, error) {
	start := time.Now()
	searchResponse, err := searchQuery.Run(
		context.Background(),
		client,
		&osquery.Options{
			Indices: []string{indexName},
			Params:  params,
		},
	)
	took := time.Since(start).String()
	if err != nil {
		return nil, fmt.Errorf("can not query for documents: %w", err)
	}
	fmt.Printf("Matched %d of %d results in %s (%dms search).\n", len(searchResponse.Hits.Hits), searchResponse.Hits.Total.Value, took, searchResponse.Took)
	return searchResponse, nil
}
