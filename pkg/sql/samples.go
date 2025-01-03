// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

const GetSampleByIDQuery = `
select s.samplingfeatureid,
s.samplingfeatureuuid,
s.samplingfeaturename,
s.samplingfeaturedescription,
s.samplingfeaturegeotypecv,
s.featuregeometrywkt,
s.elevation_m as elevationM,
s.elevationdatumcv,
s.elevationprecision,
s.elevationprecisioncomment
from odm2.samplingfeatures s
where s.samplingfeatureid = $1
`

// Modular query to get specimenids filtered by various features
// BaseQuery is extended with JOIN-modules depending on the selected filter options
// Filter query-modules can be configured with feature comparisons that are concatenated either with "and" or "or"
const GetSamplingfeatureIdsByFilterBaseQuery = `
-- modular query for specimenids and coordinates with all filter options
select distinct spec.sampleid,
coalesce(coords.latitude, 0) as latitude,
coalesce(coords.longitude, 0) as longitude,
spec.sampleName,
spec.batches,
spec.sites,
spec.publicationYear,
spec.doi,
spec.authors,
spec.minerals,
spec.hostMinerals,
spec.inclusionMinerals,
spec.rockTypes,
spec.rockClasses,
spec.inclusionTypes,
spec.geologicalSettings,
spec.geologicalAges,
spec.geologicalAgesMin,
spec.geologicalAgesMax,
spec.selectedMeasurements
from odm2.samplelistinformationextended spec
`

// Same as GetSamplingfeatureIdsByFilterBaseQuery but with translated geometries for points outside -180 to 180
// Depends on QueryModule Geometry being added
const GetSamplingFeatureIdsByFilteBaseQueryTranslated = `
-- modular query for specimenids and translated geometries with all filter options
select distinct spec.sampleid,
spec.sampleName,
spec.batches,
spec.sites,
spec.publicationYear,
spec.doi,
spec.authors,
spec.minerals,
spec.hostMinerals,
spec.inclusionMinerals,
spec.rockTypes,
spec.rockClasses,
spec.inclusionTypes,
spec.geologicalSettings,
spec.geologicalAges,
spec.geologicalAgesMin,
spec.geologicalAgesMax,
spec.selectedMeasurements,
case when geom.isInOriginalBBOX then geom.geometry else st_translate(geom.geometry, 360 * translationFactor, 0) end as translatedGeom
from odm2.samplelistinformationextended spec
`

// Filter query-module Locations
// Filter options are:
//
//	SettingName
//	Locationname lvl1
//	Locationname lvl2
//	Locationname lvl3
//	Latitude
//	Longitude
const GetSamplingfeatureIdsByFilterLocationsStart = `
join (
	-- location data
	select r_sample.samplingfeatureid as sample
	from odm2.sites s
	left join odm2.sitegeologicalsettings sgs on sgs.samplingfeatureid = s.samplingfeatureid
	left join odm2.geologicalsettings gs on gs.settingid = sgs.settingid
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
`
const GetSamplingfeatureIdsByFilterLocationsEnd = `
) loc on loc.sample = spec.sampleid
`

// Filter query-module TaxonomicClassifiers
// Filter options are:
//
//	RockType
//	RockClass
//	Mineral
//	HostMineral
//	InclusionMineral
const GetSamplingfeatureIdsByFilterTaxonomicClassifiersStart = `
join (
	-- taxonomic classifiers
	select st.samplingfeatureid
	from odm2.sampletaxonomicclassifiers st
`
const GetSamplingfeatureIdsByFilterTaxonomicClassifiersEnd = `
) tax on tax.samplingfeatureid = spec.sampleid
`

// Filter query-module Annotations
// Filter options are:
//
//	Material
//	InclusionType
//	SamplingTechnique
const GetSamplingfeatureIdsByFilterAnnotationsStart = `
join (
	-- annotations
	select sr.sampleid
	from odm2.samplerelations sr
`

const GetSamplingfeatureIdsByFilterAnnotationsMaterial = `
left join odm2.annotations ann_mat on ann_mat.annotationid = sr.annotationid and ann_mat.annotationcode = 'g_batches.material'
`

const GetSamplingfeatureIdsByFilterAnnotationsIncType = `
left join odm2.annotations ann_inc_type on ann_inc_type.annotationid = sr.annotationid and ann_inc_type.annotationcode = 'g_inclusions.inclusion_type'
`

const GetSamplingfeatureIdsByFilterAnnotationsSampTech = `
left join odm2.annotations ann_stech on ann_stech.annotationid = sr.annotationid and ann_stech.annotationcode = 'g_samples.samp_technique'
`

const GetSamplingfeatureIdsByFilterAnnotationsRimOrCore = `
left join odm2.annotations ann_roc on ann_roc.annotationid = sr.annotationid and ann_roc.annotationcode = 'g_inclusions.rim_or_core_inc'
`

const GetSamplingfeatureIdsByFilterAnnotationsEnd = `
) ann on ann.sampleid = spec.sampleid
`

// Filter query-module Results
// Filter options are defined by conjunctive filter tuples:
// (TYPE, ELEMENT, MIN, MAX)
const GetSamplingfeatureIdsByFilterResultsStartPre = `
join (
	-- results
	select distinct coalesce(
`

const GetSamplingfeatureIdsByFilterResultsStartPost = `
) as sampleid
`

const GetSamplingfeatureIdsByFilterResultsExpression = `
select distinct n.batchid as samplingfeatureid, n.sampleid
from odm2.normalizedchemistry n
`

const GetSamplingfeatureIdsByFilterResultsEnd = `
) results on results.sampleid = spec.sampleid
`

// Filter query-module Citations
// Filter options are:
//
//	ExternalIdentifier
//	Title
//	PublicationYear
const GetSamplingfeatureIdsByFilterCitationsStart = `
join (
	select distinct scd.samplingfeatureid
	from odm2.samplecitationdata scd
`

const GetSamplingfeatureIdsByFilterCitationsEnd = `
) citations on citations.samplingfeatureid = spec.sampleid
`

// Filter query-module Ages
// Filter options are:
//
//	AgeMin
//	AgeMax
//	GeologicalAge
//	GeologicalAgePrefix
const GetSamplingfeatureIdsByFilterAgesStart = `
join (
	select sa.samplingfeatureid
	from odm2.specimenages sa
`

const GetSamplingfeatureIdsByFilterAgesEnd = `
) ages on ages.samplingfeatureid = spec.sampleid
`

// Filter query-module Organizations
// Filter options are:
//
//	OrganizationName
const GestSamplingfeatureIdsByFilterOrganizationsStart = `
join (
	select 
	s.samplingfeatureid as sid,
	s.samplingfeaturedescription
	from odm2.organizations o 
	left join odm2.actionby a on a.organizationid = o.organizationid and a.roledescription != 'chief scientist'
	left join odm2.featureactions f on f.actionid = a.actionid
	left join odm2.relatedfeatures r on r.samplingfeatureid = f.samplingfeatureid and r.relationshiptypecv != 'Is identical to'
	left join odm2.samplingfeatures s on s.samplingfeatureid = r.relatedfeatureid
`

const GestSamplingfeatureIdsByFilterOrganizationsEnd = `
) organizations on spec.sampleid = organizations.sid
`

// Filter query-module Geometry
// Filter options are:
//
//	ST_WITHIN(geometry, st_wrapx(given-polygon))
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
) geom on geom.sampleid = spec.sampleid
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
) coords on coords.sampleid = spec.sampleid
`

// Wrapper for clustering
const GetSamplesClusteredWrapperPrefix = `
-- filter query with clustering
select
clusters.clusterid,
st_asText(st_convexhull(st_collect(clusters.translatedGeom))) as convexHullString,
st_asText(ST_Centroid(ST_Union(clusters.translatedGeom))) as centroidString,
array_agg(st_astext(clusters.translatedGeom)) as pointStrings,
array_agg((
	clusters.sampleID,
	clusters.sampleName,
	clusters.batches,
	clusters.sites,
	clusters.publicationYear,
	clusters.doi,
	clusters.authors,
	clusters.minerals,
	clusters.hostMinerals,
	clusters.inclusionMinerals,
	clusters.rockTypes,
	clusters.rockClasses,
	clusters.inclusionTypes,
	clusters.geologicalSettings,
	clusters.geologicalAges,
	clusters.geologicalAgesMin,
	clusters.geologicalAgesMax,
	clusters.selectedMeasurements
)) as samples
from (
	select samples.*,
	st_clusterkmeans(samples.translatedGeom, numClusters, maxDistance) over () as clusterid
	from (
`

const GetSamplesClusteredWrapperNoClusteringPrefix = `
-- filter query with clustering
select
clusters.clusterid,
'NoClustering' as convexHullString,
'NoClustering' as centroidString,
array_agg(st_astext(clusters.translatedGeom)) as pointStrings,
array_agg((
	clusters.sampleID,
	clusters.sampleName,
	clusters.batches,
	clusters.sites,
	clusters.publicationYear,
	clusters.doi,
	clusters.authors,
	clusters.minerals,
	clusters.hostMinerals,
	clusters.inclusionMinerals,
	clusters.rockTypes,
	clusters.rockClasses,
	clusters.inclusionTypes,
	clusters.geologicalSettings,
	clusters.geologicalAges,
	clusters.geologicalAgesMin,
	clusters.geologicalAgesMax,
	clusters.selectedMeasurements
)) as samples
from (
	select samples.*,
	fakeID::int4 as clusterid
	from (
`

const GetSamplesClusteredWrapperPostfix = `
	) samples
) clusters
group by clusters.clusterid
`
