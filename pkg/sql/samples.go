package sql

const GetSampleByIDQuery = `
select * from odm2.samplingfeatures s
where s.samplingfeatureid = $1
`

// Modular query to get specimenids filtered by various features
// BaseQuery is extended with JOIN-modules depending on the selected filter options
// Filter query-modules can be configured with feature comparisons that are concatenated either with "and" or "or"
const GetSamplingfeatureIdsByFilterBaseQuery = `
-- modular query for specimenids with all filter options
select distinct (case when spec.samplingfeaturedescription = 'Sample' then spec.samplingfeatureid else r.relatedfeatureid end) as sampleid
from odm2.samplingfeatures spec
left join odm2.relatedfeatures r on r.samplingfeatureid = spec.samplingfeatureid
`

// Same as GetSamplingfeatureIdsByFilterBaseQuery but with select on latitude and longitude
// Depends on QueryModule Coordinates being added
const GetSamplingfeatureIdsByFilterBaseQueryWithCoords = `
-- modular query for specimenids and coordinates with all filter options
select distinct (case when spec.samplingfeaturedescription = 'Sample' then spec.samplingfeatureid else r.relatedfeatureid end) as sampleid,
coords.latitude,
coords.longitude
from odm2.samplingfeatures spec
left join odm2.relatedfeatures r on r.samplingfeatureid = spec.samplingfeatureid
`

// Alternative query beginning to GetSamplingfeatureIdsByFilterBaseQuery but with translated geometries for points outside -180 to 180
// and refactored for two-step clustering
// Depends on QueryModule Geometry being added
const GetSamplingFeatureIdsByFilterBaseQueryForClusters = `
-- modular query for specimenids and translated geometries with all filter options
select string_agg(tmp.SampleString, ',') as valuesString,
count(tmp.SampleID) as numSamples
from (
	select distinct (case when spec.samplingfeaturedescription = 'Sample' then spec.samplingfeatureid::varchar else r.relatedfeatureid::varchar end) as sampleid,
	case when geom.isInOriginalBBOX then st_astext(geom.geometry) else st_astext(st_translate(geom.geometry, 360 * translationFactor, 0)) end as translatedGeom,
	'('||(case when spec.samplingfeaturedescription = 'Sample' then spec.samplingfeatureid::varchar else r.relatedfeatureid::varchar end)||','''||
	case when geom.isInOriginalBBOX then st_astext(geom.geometry) else st_astext(st_translate(geom.geometry, 360 * 1, 0)) end || ''')' as SampleString
	from odm2.samplingfeatures spec
	left join odm2.relatedfeatures r on r.samplingfeatureid = spec.samplingfeatureid
`

const GetSamplingFeatureIdsByFilterBaseQueryTranslatedEnd = `
) tmp
`

// Filter query-module Locations
// Filter options are:
// 		Setting
//		Locationname lvl1
//		Locationname lvl2
//		Locationname lvl3
//		Latitude
//		Longitude
const GetSamplingfeatureIdsByFilterLocationsStart = `
join (
	-- location data
	select r_sample.samplingfeatureid as sample,
	r_batch.samplingfeatureid as batch
	from odm2.sites s
	left join 
	(
		select sg.samplingfeatureid, sg.locationname
		from odm2.sitegeolocations sg
		left join odm2.geolocations g on g.geolocationid = sg.geolocationid
		where right(g.locationhierarchy::varchar, 3) = '100' --toplevel
		group by sg.samplingfeatureid, sg.locationname -- multiple entries for same locationname
	) toplevelloc on toplevelloc.samplingfeatureid = s.samplingfeatureid 
	left join
	(
		select sg.samplingfeatureid, sg.locationname
		from odm2.sitegeolocations sg
		left join odm2.geolocations g on g.geolocationid = sg.geolocationid
		where right(g.locationhierarchy::varchar, 3) = '200' --second level
		group by sg.samplingfeatureid, sg.locationname -- multiple entries for same locationname
	) secondlevelloc on secondlevelloc.samplingfeatureid = s.samplingfeatureid 
	left join
	(
		select sg.samplingfeatureid, sg.locationname
		from odm2.sitegeolocations sg
		left join odm2.geolocations g on g.geolocationid = sg.geolocationid
		where right(g.locationhierarchy::varchar, 3) = '300' --third level
		group by sg.samplingfeatureid, sg.locationname -- multiple entries for same locationname
	) thirdlevelloc on thirdlevelloc.samplingfeatureid = s.samplingfeatureid
	left join odm2.relatedfeatures r_sample on r_sample.relatedfeatureid = s.samplingfeatureid -- samples for each location
	left join odm2.relatedfeatures r_batch on r_batch.relatedfeatureid = r_sample.samplingfeatureid -- batches for each sample
`
const GetSamplingfeatureIdsByFilterLocationsEnd = `
) loc on loc.sample = spec.samplingfeatureid or loc.batch = spec.samplingfeatureid
`

// Filter query-module TaxonomicClassifiers
// Filter options are:
// 		RockType
//		RockClass
//		Mineral
const GetSamplingfeatureIdsByFilterTaxonomicClassifiersStart = `
join (
	-- taxonomic classifiers
	select stc.samplingfeatureid
	from odm2.specimentaxonomicclassifiers stc
	left join odm2.taxonomicclassifiers tax_type on tax_type.taxonomicclassifierid = stc.taxonomicclassifierid and tax_type.taxonomicclassifiertypecv = 'Rock'
	left join odm2.taxonomicclassifiers tax_class on tax_class.taxonomicclassifierid = stc.taxonomicclassifierid and tax_class.taxonomicclassifiertypecv = 'Lithology'
	left join odm2.taxonomicclassifiers tax_min on tax_min.taxonomicclassifierid = stc.taxonomicclassifierid and tax_min.taxonomicclassifierdescription  = 'Mineral Classification from GEOROC'
`
const GetSamplingfeatureIdsByFilterTaxonomicClassifiersEnd = `
) tax on tax.samplingfeatureid = spec.samplingfeatureid
`

// Filter query-module Annotations
// Filter options are:
// 		Material
//		InclusionType
//		SamplingTechnique
const GetSamplingfeatureIdsByFilterAnnotationsStart = `
join (
	-- annotations
	select r.relatedfeatureid as sampleid
	from odm2.samplingfeatureannotations sann
	left join odm2.annotations ann_mat on ann_mat.annotationid = sann.annotationid and ann_mat.annotationcode = 'g_batches.material'
	left join odm2.annotations ann_inc_type on ann_inc_type.annotationid = sann.annotationid and ann_inc_type.annotationcode = 'g_inclusions.inclusion_type'
	left join odm2.annotations ann_samp_tech on ann_samp_tech.annotationid = sann.annotationid and ann_samp_tech.annotationcode = 'g_samples.samp_technique'
	left join odm2.annotations ann_rim_or_core on ann_rim_or_core.annotationid = sann.annotationid and ann_rim_or_core.annotationcode = 'g_inclusions.rim_or_core_inc'
	left join odm2.relatedfeatures r on r.samplingfeatureid = sann.samplingfeatureid and r.relationshiptypecv != 'Is identical to'
`
const GetSamplingfeatureIdsByFilterAnnotationsEnd = `
) ann on ann.sampleid = spec.samplingfeatureid
`

// Filter query-module Results
// Filter options are:
// 		MeasuredItem
//		ItemType
//		Unit
//		Value
const GetSamplingfeatureIdsByFilterResultsStart = `
join (
	-- results
	select mv.samplingfeatureid 
	from odm2.measuredvalues mv
`
const GetSamplingfeatureIdsByFilterResultsEnd = `
) results on results.samplingfeatureid = spec.samplingfeatureid
`

// Filter query-module Citations
// Filter options are:
// 		DOI
// 		Title
//		PublicationYear
const GetSamplingfeatureIdsByFilterCitationsStart = `
join (
	select distinct cs.samplingfeatureid
	from odm2.citations c
	left join odm2.authorlists al on al.citationid = c.citationid
	left join odm2.people p on p.personid = al.personid
	left join odm2.citationexternalidentifiers cid on cid.citationid = c.citationid
	left join odm2.externalidentifiersystems e on e.externalidentifiersystemid = cid.externalidentifiersystemid
	left join odm2.citationsamplingfeatures cs on cs.citationid = c.citationid
`

const GetSamplingfeatureIdsByFilterCitationsEnd = `
) citations on citations.samplingfeatureid = spec.samplingfeatureid
`

// Filter query-module Ages
// Filter options are:
//		AgeMin
//		AgeMax
//		GeologicalAge
//		GeologicalAgePrefix
const GetSamplingfeatureIdsByFilterAgesStart = `
join (
	select sa.samplingfeatureid
	from odm2.specimenages sa
`

const GetSamplingfeatureIdsByFilterAgesEnd = `
) ages on ages.samplingfeatureid = spec.samplingfeatureid
`

// Filter query-module Organizations
// Filter options are:
// 		OrganizationName
const GestSamplingfeatureIdsByFilterOrganizationsStart = `
join (
	select 
	f.samplingfeatureid as fid,
	s.samplingfeatureid as sid,
	s.samplingfeaturedescription
	from odm2.organizations o 
	left join odm2.actionby a on a.organizationid = o.organizationid and a.roledescription != 'chief scientist'
	left join odm2.featureactions f on f.actionid = a.actionid
	left join odm2.relatedfeatures r on r.samplingfeatureid = f.samplingfeatureid and r.relationshiptypecv != 'Is identical to'
	left join odm2.samplingfeatures s on s.samplingfeatureid = r.relatedfeatureid
`

const GestSamplingfeatureIdsByFilterOrganizationsEnd = `
) organizations on spec.samplingfeatureid = case when organizations.samplingfeaturedescription = 'Sample' then organizations.sid else fid end
`

// Filter query-module Geometry
// Filter options are:
// 		ST_WITHIN(geometry, st_wrapx(given-polygon))
const GestSamplingfeatureIdsByFilterGeometryStart = `
join (
select r.samplingfeatureid as sampleid,
sg.geometry
from odm2.sitegeometries sg 
join odm2.relatedfeatures r on r.relatedfeatureid = sg.samplingfeatureid
`

// Same as GestSamplingfeatureIdsByFilterGeometryStart but with added check if geometries are in original bbox
const GestSamplingfeatureIdsByFilterGeometryBBOXStart = `
join (
select r.samplingfeatureid as sampleid,
sg.geometry,
case when st_within(sg.geometry, ST_GEOMETRYFROMTEXT(bboxPolygon, 4326)) then true else false end as isInOriginalBBOX
from odm2.sitegeometries sg 
join odm2.relatedfeatures r on r.relatedfeatureid = sg.samplingfeatureid
`

const GestSamplingfeatureIdsByFilterGeometryEnd = `
) geom on geom.sampleid = spec.samplingfeatureid
`

// Filter query-module coordinates
// No filter option but adds coordinates to each sampleID
const GetGestSamplingfeatureIdsByFilterCoordinates = `
left join 
(
	select r.samplingfeatureid as sampleid,
	s.latitude,
	s.longitude
	from odm2.sites s
	join odm2.relatedfeatures r on r.relatedfeatureid = s.samplingfeatureid
) coords on coords.sampleid = case when spec.samplingfeaturedescription = 'Sample' then spec.samplingfeatureid else r.relatedfeatureid end
`

// Wrapper for clustering
const GetSamplesClusteredWrapperPrefix = `
-- filter query with clustering
select
clusters.clusterid,
st_convexhull(st_collect(clusters.translatedGeom)) as convexHull,
ST_Centroid(ST_Union(clusters.translatedGeom)) as centroid,
array_agg(clusters.sampleid) as samples,
array_agg(clusters.sampleid || ',' || clusters.translatedGeom) as pointsWithIds
from (
	select samples.sampleid,
	samples.translatedGeom,
	st_clusterkmeans(samples.translatedGeom, numClusters, maxDistance) over () as clusterid
	from (
`

const GetSamplesClusteredWrapperPostfix = `
	) as samples (sampleid, translatedGeom)
) clusters
group by clusters.clusterid
`
