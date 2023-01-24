package sql

const FullDataByIdQuery = `
select
samples.SamplingFeatureID as sample_num,
samples.SamplingFeatureUUID as unique_id,
array(select unnest(samples.batchids)) as batches,
refs.references,
samples.samplingfeaturename as sampleids,
coalesce (loc.loc_names, array['Unknown']) as location_names,
coalesce (loc.loc_types, array['Unknown']) as location_types,
loc.loc_data[1] as loc_data,
(regexp_match(loc.elevation[1], 'ELEVATION_MIN=(-?\d+\.?\d*);'))[1] as elevation_min,
(regexp_match(loc.elevation[1], 'ELEVATION_MAX=(-?\d+\.?\d*)'))[1] as elevation_max,
-- no samp_technique in odm2
(loc.land_or_sea)[1] as land_or_sea,
samples.rock_type as rock_types,
samples.rock_class as rock_classes,
samples.specimentexture as rock_textures,
(samples.specimenagemin)[1] as age_min,
(samples.specimenagemax)[1] as age_max,
samples.material as materials,
samples.mineral as minerals,
samples.inclusion_type as inclusion_types,
loc.loc_data[1].samplingfeatureid as location_num,
loc.loc_data[1].latitude as latitude,
loc.loc_data[1].longitude as longitude,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LATITUDE_MIN=(-?\d+\.?\d*);'))[1] as latitude_min,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LATITUDE_MAX=(-?\d+\.?\d*);'))[1] as latitude_max,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LONGITUDE_MIN=(-?\d+\.?\d*);'))[1] as longitude_min,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LONGITUDE_MAX=(-?\d+\.?\d*)'))[1] as longitude_max,
(loc.setting)[1] as tectonic_setting,
methods.method_acronyms as method,
methods.method_comments as comment,
(methods.institution)[1] as institution, -- actionBy seems to be sparsely filled
results.items_measured as item_name,
results.item_types as item_group,
--rockmode_num not in odm2, 
results.standard_names as standard_names,
results.standard_values as standard_values, -- somehow different values than in the dbo-query
results.values_meas as values, -- somehow different values than in the dbo-query
results.units as units
from 
(
	select samples.samplingfeatureid, 
	samples.samplingfeatureuuid, 
	(array_agg(batches.batches)) as batchids,
	samples.samplingfeaturename,
	(array_agg(tax_type.taxonomicclassifiername)) as rock_type,
	(array_agg(tax_class.taxonomicclassifiername)) as rock_class,
	(array_agg(tax_min.taxonomicclassifiercommonname)) as mineral,
	(array_agg(spec.specimentexture)) as specimentexture,
	(array_agg(sage.specimenagemin)) as specimenagemin,
	(array_agg(sage.specimenagemax)) as specimenagemax,
	(array_agg(ann_mat.annotationtext)) as material,
	(array_agg(ann_inc_type.annotationtext)) as inclusion_type
	from odm2.samplingfeatures samples
	left join 
	(
		select s.samplingfeatureid,
		array_agg(r.samplingfeatureid) as batches 
		from odm2.samplingfeatures s 
		left join odm2.relatedfeatures r on r.relatedfeatureid = s.SamplingFeatureID and r.relationshiptypecv = 'Is child of'
		where s.samplingfeatureid = $1
		group by s.samplingfeatureid
	) batches on batches.samplingfeatureid = samples.samplingfeatureid 
	left join odm2.specimens spec on spec.samplingfeatureid = samples.samplingfeatureid
	left join
	(
		select stc.samplingfeatureid, tax_type.taxonomicclassifiername 
		from odm2.specimentaxonomicclassifiers stc
		left join odm2.taxonomicclassifiers tax_type on tax_type.taxonomicclassifierid = stc.taxonomicclassifierid and tax_type.taxonomicclassifiertypecv = 'Rock'
		where tax_type.taxonomicclassifierid is not null
	) tax_type on tax_type.samplingfeatureid = spec.samplingfeatureid
	left join
	(
		select stc.samplingfeatureid, tax_class.taxonomicclassifiername 
		from odm2.specimentaxonomicclassifiers stc
		left join odm2.taxonomicclassifiers tax_class on tax_class.taxonomicclassifierid = stc.taxonomicclassifierid and tax_class.taxonomicclassifiertypecv = 'Lithology'
		where tax_class.taxonomicclassifierid is not null
	) tax_class on tax_class.samplingfeatureid = spec.samplingfeatureid
	left join
	(
		select stc.samplingfeatureid, tax_min.taxonomicclassifiercommonname 
		from odm2.specimentaxonomicclassifiers stc
		left join odm2.taxonomicclassifiers tax_min on tax_min.taxonomicclassifierid = stc.taxonomicclassifierid and tax_min.taxonomicclassifierdescription  = 'Mineral Classification from GEOROC'
		where tax_min.taxonomicclassifierid is not null
	) tax_min on tax_min.samplingfeatureid = spec.samplingfeatureid
	left join odm2.specimenages sage on sage.samplingfeatureid = samples.samplingfeatureid
	left join odm2.samplingfeatureannotations sann on sann.samplingfeatureid = samples.samplingfeatureid -- x2
	left join odm2.annotations ann_mat on ann_mat.annotationid = sann.annotationid and ann_mat.annotationcode = 'g_batches.material'
	left join odm2.annotations ann_inc_type on ann_inc_type.annotationid = sann.annotationid and ann_inc_type.annotationcode = 'g_inclusions.inclusion_type'
	where samples.samplingfeaturedescription = 'Sample' 
	and samples.samplingfeatureid = $1
	group by samples.samplingfeatureid 
) samples
left join 
(	
	select q.samplingfeatureid, array_agg(q.*) as references from (
		-- query references and authors
		select stc_ref.samplingfeatureid,
		json_build_object('ref_num', c_ref.citationid , 'title', c_ref.title , 'journal', c_ref.journal , 'pages', c_ref.firstpage || '-' || c_ref.lastpage , 'year', c_ref.publicationyear , 'doi', cei_ref.citationexternalidentifier) as reference,
		array_agg(distinct p_ref.*) as authors 
		from odm2.specimentaxonomicclassifiers stc_ref
		left join odm2.citations c_ref on c_ref.citationid = stc_ref.citationid
		left join odm2.citationexternalidentifiers cei_ref on cei_ref.citationid = c_ref.citationid and cei_ref.externalidentifiersystemid = 1 -- id of externalidentifiersystem "DOI"
		left join odm2.authorlists a_ref on a_ref.citationid = c_ref.citationid
		left join odm2.people p_ref on p_ref.personid = a_ref.personid 
		left join odm2.affiliations af_ref on af_ref.personid = p_ref.personid
		where stc_ref.samplingfeatureid = $1
		group by stc_ref.samplingfeatureid, c_ref.citationid, c_ref.title, c_ref.journal, c_ref.firstpage, c_ref.lastpage, c_ref.publicationyear, cei_ref.citationexternalidentifier 
	) q
	group by q.samplingfeatureid
) refs on refs.samplingfeatureid = samples.samplingfeatureid
left join
(
	-- query locations - Until geolocations is refactored, we get multiple outputs here
	select rel_loc.samplingfeatureid, 
	(array_agg(distinct si_loc.*)) as loc_data, 
	array_remove(array_agg(distinct sg_loc.locationname), null) as loc_names, 
	array_remove(array_agg(distinct g_loc.geolocationtype), null) as loc_types ,
	array_remove(array_agg(loc.elevationprecisioncomment), null) as elevation,
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
	left join odm2.organizations org on org.organizationid = ab_meth.organizationid
	where rel_meth.relatedfeatureid = $1
	group by rel_meth.relatedfeatureid
) methods on methods.id = samples.samplingfeatureid
left join (
	-- query results
	select rel_res.relatedfeatureid as id,
	array_agg(vars.variablecode) as items_measured,
	array_agg(vars.variabletypecode) as item_types,
	coalesce(array_agg(std.std_names), array['Unknown']) as standard_names,
	coalesce(array_agg(std.std_values), array[-999]) as standard_values,
	array_agg(chem_vals.datavalue) as values_meas,
	array_agg(u.unitgeoroc) as units 
	from odm2.relatedfeatures rel_res
	join odm2.featureactions f_res on f_res.samplingfeatureid = rel_res.samplingfeatureid
	left join odm2.results res on res.featureactionid = f_res.featureactionid
	left join odm2.variables vars on vars.variableid = res.variableid 
	left join 
	(
		select relf.samplingfeatureid,
		array_agg(standards.standardname) as std_names,
		array_agg(standards.standardvalue) as std_values 
		from odm2.relatedfeatures relf
		join odm2.featureactions fa on fa.samplingfeatureid = relf.samplingfeatureid
		left join odm2.standards standards on standards.actionid = fa.actionid
		where relf.relatedfeatureid = $1
		group by relf.samplingfeatureid
	)std on std.samplingfeatureid = rel_res.samplingfeatureid 
	left join odm2.measurementresultvalues chem_vals on chem_vals.valueid = res.resultid 
	left join odm2.units u on u.unitsid = res.unitsid 
	where rel_res.relatedfeatureid = $1
	and rel_res.relationshiptypecv = 'Is child of'
	group by rel_res.relatedfeatureid
) results on results.id = samples.samplingfeatureid
`
