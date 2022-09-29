package sql

const FullDataByIdQuery = `
select
samples.SamplingFeatureID as sample_num,
(array_agg(samples.SamplingFeatureUUID))[1] as unique_id,
array_agg(samples.batches) as batches,
(array_agg(jsonb_build_object('reference', refs.reference, 'authors', refs.authors)))[1] as references,
array_agg(distinct samples.samplingfeaturename) as sampleids,
(array_agg(distinct coalesce (loc.loc_names, array['Unknown']))) as location_names,
(array_agg(distinct coalesce (loc.loc_types, array['Unknown']))) as location_types,
(array_agg(distinct loc.loc_data))[1][1] as loc_data,
(array_agg(distinct loc.elevation))[1] as elevation_min,
(array_agg(distinct loc.elevation))[1] as elevation_max,
-- no samp_technique in odm2
array_agg(distinct loc.land_or_sea) as land_or_sea,
(array_agg(distinct samples.rock_type))[1] as rock_type,
(array_agg(distinct samples.rock_class))[1] as rock_class,
(array_agg(distinct samples.specimentexture))[1] as rock_texture,
(array_agg(distinct samples.specimenagemin))[1] as age_min,
(array_agg(distinct samples.specimenagemax))[1] as age_max,
array_agg(samples.material) as material,
array_agg(samples.mineral) as mineral,
(array_agg(samples.inclusion_type)) as inclusion_type,
(array_agg(distinct loc.loc_data))[1][1].samplingfeatureid as location_num,
(array_agg(distinct loc.loc_data))[1][1].latitude as latitude,
(array_agg(distinct loc.loc_data))[1][1].longitude as longitude,
(regexp_match((array_agg(distinct loc.loc_data))[1][1].locationcomment, 'LATITUDE_MIN=(-?\d+\.?\d*);'))[1] as latitude_min,
(regexp_match((array_agg(distinct loc.loc_data))[1][1].locationcomment, 'LATITUDE_MAX=(-?\d+\.?\d*);'))[1] as latitude_max,
(regexp_match((array_agg(distinct loc.loc_data))[1][1].locationcomment, 'LONGITUDE_MIN=(-?\d+\.?\d*);'))[1] as longitude_min,
(regexp_match((array_agg(distinct loc.loc_data))[1][1].locationcomment, 'LONGITUDE_MAX=(-?\d+\.?\d*)'))[1] as longitude_max,
array_agg(distinct loc.setting) as tectonic_setting,
array_agg(methods.method_acronyms) as method,
array_agg(methods.method_comments) as comment,
array_agg(methods.institution) as institution, -- actionBy seems to be sparsely filled
array_agg(distinct results.items_measured) as item_name,
array_agg(distinct results.item_types) as item_group,
--rockmode_num not in odm2, 
(json_agg(results.standard_names)) as standard_names,
(json_agg(results.standard_values)) as standard_values, -- somehow different values than in the dbo-query
array_agg(results.values_meas) as values, -- somehow different values than in the dbo-query
array_agg(results.units) as units
from 
(
	select samples.samplingfeatureid, samples.samplingfeatureuuid, r.samplingfeatureid as batches,samples.samplingfeaturename,
	tax_type.taxonomicclassifiername as rock_type, tax_class.taxonomicclassifiername as rock_class, tax_min.taxonomicclassifiercommonname as mineral,
	spec.specimentexture, sage.specimenagemin, sage.specimenagemax,
	ann_mat.annotationtext as material, ann_inc_type.annotationtext as inclusion_type
	from odm2.samplingfeatures samples
	left join odm2.relatedfeatures r on r.relatedfeatureid  = samples.SamplingFeatureID
	left join odm2.specimens spec on spec.samplingfeatureid = samples.samplingfeatureid
	left join odm2.specimentaxonomicclassifiers stc on stc.samplingfeatureid = spec.samplingfeatureid 
	left join odm2.taxonomicclassifiers tax_type on tax_type.taxonomicclassifierid = stc.taxonomicclassifierid and tax_type.taxonomicclassifiertypecv = 'Rock'
	left join odm2.taxonomicclassifiers tax_class on tax_class.taxonomicclassifierid = stc.taxonomicclassifierid and tax_class.taxonomicclassifiertypecv = 'Lithology'
	left join odm2.taxonomicclassifiers tax_min on tax_min.taxonomicclassifierid = stc.taxonomicclassifierid and tax_min.taxonomicclassifierdescription  = 'Mineral Classification from GEOROC'
	left join odm2.specimenages sage on sage.samplingfeatureid = samples.samplingfeatureid
	left join odm2.samplingfeatureannotations sann on sann.samplingfeatureid = samples.samplingfeatureid
	left join odm2.annotations ann_mat on ann_mat.annotationid = sann.annotationid and ann_mat.annotationcode = 'g_batches.material'
	left join odm2.annotations ann_inc_type on ann_inc_type.annotationid = sann.annotationid and ann_inc_type.annotationcode = 'g_inclusions.inclusion_type'
	where samples.samplingfeaturedescription = 'Sample' and samples.samplingfeatureid = $1
) samples
left join 
(
	-- query references and authors
	select stc_ref.samplingfeatureid,
	json_build_object('ref_num', c_ref.citationid , 'title', c_ref.title , 'journal', c_ref.journal , 'pages', c_ref.firstpage || '-' || c_ref.lastpage , 'year', c_ref.publicationyear , 'doi', cei_ref.citationexternalidentifier) as reference,
	array_agg(distinct p_ref.*) as authors 
	from odm2.specimentaxonomicclassifiers stc_ref
	left join odm2.citations c_ref on c_ref.citationid = stc_ref.citationid
	left join odm2.citationexternalidentifiers cei_ref on cei_ref.citationid = c_ref.citationid and cei_ref.externalidentifiersystemname = 'DOI'
	left join odm2.authorlists a_ref on a_ref.citationid = c_ref.citationid
	left join odm2.people p_ref on p_ref.personid = a_ref.personid 
	left join odm2.affiliations af_ref on af_ref.personid = p_ref.personid
	where stc_ref.samplingfeatureid = $1
	group by stc_ref.samplingfeatureid, c_ref.citationid, c_ref.title, c_ref.journal, c_ref.firstpage, c_ref.lastpage, c_ref.publicationyear, cei_ref.citationexternalidentifier 
) refs on refs.samplingfeatureid = samples.samplingfeatureid
left join
(
	-- query locations
	select rel_loc.samplingfeatureid, 
	(array_agg(distinct si_loc.*)) as loc_data, 
	array_remove(array_agg(sg_loc.locationname), null) as loc_names, 
	array_remove(array_agg(g_loc.geolocationtype), null) as loc_types ,
	array_remove(array_agg(loc.elevationcomment), null) as elevation,
	array_remove(array_agg(distinct si_loc.sitedescription), null) as land_or_sea,
	array_remove(array_agg(distinct si_loc.setting), null) as setting
	from odm2.relatedfeatures rel_loc
	left join odm2.samplingfeatures loc on loc.samplingfeatureid = rel_loc.relatedfeatureid 
	left join odm2.sites si_loc on si_loc.samplingfeatureid = rel_loc.relatedfeatureid 
	left join odm2.sitegeolocations sg_loc on sg_loc.samplingfeatureid  = si_loc.samplingfeatureid
	left join odm2.geolocations g_loc on g_loc.geolocationid = sg_loc.geolocationid 
	where rel_loc.samplingfeatureid = $1
	group by rel_loc.samplingfeatureid 
) loc on loc.samplingfeatureid = samples.samplingfeatureid
left join (
	-- query methods
	select rel_meth.relatedfeatureid as id,
	(array_agg(distinct meth.methodcode))[1] as method_acronyms,
	(array_agg(distinct a_meth.actiondescription))[1] as method_comments,
	array_agg(distinct org.organizationname) as institution
	from odm2.relatedfeatures rel_meth
	left join odm2.featureactions f_meth on f_meth.samplingfeatureid = rel_meth.samplingfeatureid
	left join odm2.actions a_meth on a_meth.actionid = f_meth.actionid
	left join odm2.methods meth on meth.methodid = a_meth.methodid
	left join odm2.actionby ab_meth on ab_meth.actionid = a_meth.actionid 
	left join odm2.organizations org on org.organizationid = ab_meth.organizationid::int4
	where rel_meth.relatedfeatureid = $1
	group by rel_meth.relatedfeatureid
) methods on methods.id = samples.samplingfeatureid
left join (
	-- query results
	select rel_res.relatedfeatureid as id,
	array_agg(vars.variablecode) as items_measured,
	array_agg(vars.variabletypecode) as item_types,
	coalesce(array_agg(distinct rr.relationdescription), array['Unknown']) as standard_names,
	coalesce(array_agg(distinct std_vals.datavalue), array[-999]) as standard_values,
	array_agg(chem_vals.datavalue) as values_meas,
	array_agg(u.unitgeoroc) as units 
	from odm2.relatedfeatures rel_res
	left join odm2.featureactions f_res on f_res.samplingfeatureid = rel_res.samplingfeatureid -- add sample-featureactions? or f_res.samplingfeatureid = rel_res.relatedfeatureid
	left join odm2.results res on res.featureactionid = f_res.featureactionid
	left join odm2.variables vars on vars.variableid = res.variableid 
	left join odm2.relatedresults rr on rr.resultid = res.resultid 
	left join odm2.measurementresultvalues std_vals on std_vals.valueid = rr.relatedresultid 
	left join odm2.measurementresultvalues chem_vals on chem_vals.valueid = res.resultid 
	left join odm2.units u on u.unitsid = res.unitsid 
	where rel_res.relatedfeatureid = $1
	group by rel_res.relatedfeatureid
) results on results.id = samples.samplingfeatureid
group by samples.SamplingFeatureID
;
`
