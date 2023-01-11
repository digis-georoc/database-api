package sql

const GetSamplesByGeoSettingQuery = `
select 
s.samplingfeatureid,
spec.samplingfeatureid as specimen,
array_agg(distinct s.latitude) as lat ,
array_agg(distinct s.longitude) as long,
array_agg(distinct s.setting) as setting,
array_agg(distinct toplevelloc.locationname) as loc1,
array_agg(distinct secondlevelloc.locationname) as loc2,
array_agg(distinct thirdlevelloc.locationname) as loc3,
array_agg(spec.specimentexture) as texture,
array_agg(tax_type.taxonomicclassifiername) as rock_type,
array_agg(tax_class.taxonomicclassifiername) as rock_class,
array_agg(tax_min.taxonomicclassifiercommonname) as mineral, -- missing values
array_agg(ann_mat.annotationtext) as material,
array_agg(ann_inc_type.annotationtext) as inclusion_type, -- missing values
array_agg(ann_samp_tech.annotationtext) as samp_technique,
array_agg(distinct sf.samplingfeaturename) as sample_names,
array_agg(distinct s.sitedescription) as land_or_sea,
array_agg(ann_rim_or_core.annotationtext) as rim_or_core -- missing values
from odm2.sites s
left join 
(
	select sg.samplingfeatureid, sg.locationname
	from odm2.sitegeolocations sg
	left join odm2.geolocations g on g.geolocationid = sg.geolocationid
	where right(g.locationhierarchy::varchar, 3) = '100' --toplevel
) toplevelloc on toplevelloc.samplingfeatureid = s.samplingfeatureid 
left join
(
	select sg.samplingfeatureid, sg.locationname
	from odm2.sitegeolocations sg
	left join odm2.geolocations g on g.geolocationid = sg.geolocationid
	where right(g.locationhierarchy::varchar, 3) = '200' --second level
) secondlevelloc on secondlevelloc.samplingfeatureid = s.samplingfeatureid 
left join
(
	select sg.samplingfeatureid, sg.locationname
	from odm2.sitegeolocations sg
	left join odm2.geolocations g on g.geolocationid = sg.geolocationid
	where right(g.locationhierarchy::varchar, 3) = '300' --third level
) thirdlevelloc on thirdlevelloc.samplingfeatureid = s.samplingfeatureid
left join odm2.relatedfeatures rsamp on rsamp.relatedfeatureid = s.samplingfeatureid
left join odm2.relatedfeatures rbatches on rbatches.relatedfeatureid = rsamp.samplingfeatureid 
left join odm2.specimens spec on spec.samplingfeatureid = rsamp.samplingfeatureid or spec.samplingfeatureid = rbatches.samplingfeatureid
left join odm2.specimentaxonomicclassifiers stc on stc.samplingfeatureid = spec.samplingfeatureid
left join odm2.taxonomicclassifiers tax_type on tax_type.taxonomicclassifierid = stc.taxonomicclassifierid and tax_type.taxonomicclassifiertypecv = 'Rock'
left join odm2.taxonomicclassifiers tax_class on tax_class.taxonomicclassifierid = stc.taxonomicclassifierid and tax_class.taxonomicclassifiertypecv = 'Lithology'
left join odm2.taxonomicclassifiers tax_min on tax_min.taxonomicclassifierid = stc.taxonomicclassifierid and tax_min.taxonomicclassifierdescription  = 'Mineral Classification from GEOROC'
left join odm2.samplingfeatureannotations sann on sann.samplingfeatureid = spec.samplingfeatureid
left join odm2.annotations ann_mat on ann_mat.annotationid = sann.annotationid and ann_mat.annotationcode = 'g_batches.material'
left join odm2.annotations ann_inc_type on ann_inc_type.annotationid = sann.annotationid and ann_inc_type.annotationcode = 'g_inclusions.inclusion_type'
left join odm2.annotations ann_samp_tech on ann_samp_tech.annotationid = sann.annotationid and ann_samp_tech.annotationcode = 'g_samples.samp_technique'
left join odm2.annotations ann_rim_or_core on ann_rim_or_core.annotationid = sann.annotationid and ann_rim_or_core.annotationcode = 'g_inclusions.rim_or_core_inc'
left join odm2.samplingfeatures sf on sf.samplingfeatureid = rsamp.samplingfeatureid
group by s.samplingfeatureid, spec.samplingfeatureid
`
