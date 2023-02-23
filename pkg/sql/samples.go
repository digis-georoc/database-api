package sql

// Modular query to get specimenids filtered by various features
// BaseQuery is extended with JOIN-modules depending on the selected filter options
// Filter query-modules can be configured with feature comparisons that are concatenated either with "and" or "or"
const GetSamplingfeatureIdsByFilterBaseQuery = `
-- modular query for specimenids with all filter options
select distinct spec.samplingfeatureid
from odm2.specimens spec
`

// Filter query-module Locations
// Filter options are:
// 		Setting
//		Locationname lvl1
//		Locationname lvl2
//		Locationname lvl3
const GetSamplingfeatureIdsByFilterLocationsStart = `
join (
	-- location data
	select r_sample.samplingfeatureid as samples,
	r_batch.samplingfeatureid as batches
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
) loc on loc.samples = spec.samplingfeatureid or loc.batches = spec.samplingfeatureid
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
	select sann.samplingfeatureid
	from odm2.samplingfeatureannotations sann
	left join odm2.annotations ann_mat on ann_mat.annotationid = sann.annotationid and ann_mat.annotationcode = 'g_batches.material'
	left join odm2.annotations ann_inc_type on ann_inc_type.annotationid = sann.annotationid and ann_inc_type.annotationcode = 'g_inclusions.inclusion_type'
	left join odm2.annotations ann_samp_tech on ann_samp_tech.annotationid = sann.annotationid and ann_samp_tech.annotationcode = 'g_samples.samp_technique'
	left join odm2.annotations ann_rim_or_core on ann_rim_or_core.annotationid = sann.annotationid and ann_rim_or_core.annotationcode = 'g_inclusions.rim_or_core_inc'
`
const GetSamplingfeatureIdsByFilterAnnotationsEnd = `
) ann on ann.samplingfeatureid = spec.samplingfeatureid
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
