// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

// "IN $1" wont work with array but "= ANY ($1)" does
const FullDataByMultiIdQuery = `
select
samples.SamplingFeatureID as sampleID,
samples.uuid as uniqueID,
samples.references,
samples.name as samplename,
coalesce (loc.loc_names, array['Unknown']) as locationNames,
coalesce (loc.loc_types, array['Unknown']) as locationTypes,
(regexp_match(loc.elevation[1], 'ELEVATION_MIN=(-?\d+\.?\d*);'))[1] as elevationMin,
(regexp_match(loc.elevation[1], 'ELEVATION_MAX=(-?\d+\.?\d*)'))[1] as elevationMax,
(loc.land_or_sea)[1] as landOrSea,
samples.rock_types as rockTypes,
samples.rock_classes as rockClasses,
samples.specimentextures as rockTextures,
(samples.specimenagemin)[1] as ageMin,
(samples.specimenagemax)[1] as ageMax,
(samples.eruptiondate)[1] as eruptiondate,
(samples.geologicalage)[1] as geologicalage,
loc.loc_data[1].samplingfeatureid as locationNum,
loc.loc_data[1].latitude as latitude,
loc.loc_data[1].longitude as longitude,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LATITUDE_MIN=(-?\d+\.?\d*);'))[1] as latitudeMin,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LATITUDE_MAX=(-?\d+\.?\d*);'))[1] as latitudeMax,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LONGITUDE_MIN=(-?\d+\.?\d*);'))[1] as longitudeMin,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LONGITUDE_MAX=(-?\d+\.?\d*)'))[1] as longitudeMax,
(loc.setting)[1] as tectonicSetting,
coalesce(loc.location_comments, array[]::varchar[])  as locationComments,
coalesce(methods.method_acronyms, array[]::varchar[]) as methods,
coalesce(methods.method_comments, array[]::varchar[])  as comments,
coalesce(methods.institutions, array[]::varchar[])  as institutions, -- actionBy seems to be sparsely filled
coalesce(results.results, '[]'::jsonb) as results ,
samples.alteration as alteration,
samples.samp_technique as samplingTechnique,
samples.dd_min as drillDepthMin,
samples.dd_max as drillDepthMax,
samples.batchdata
from 
(
	select samples.samplingfeatureid, 
	(array_agg(samples.samplingfeatureuuid))[1] as uuid,
	(array_agg(samples.samplingfeaturename))[1] as name,
	array(select unnest(array_agg(distinct stc.rockTypeObj))) as rock_types,
	array(select unnest(array_agg(distinct stc.rockClassObj))) as rock_classes,
	(array_remove(array_agg(distinct spec.specimentexture), null)) as specimentextures,
	(array_agg(sage.specimenagemin)) as specimenagemin,
	(array_agg(sage.specimenagemax)) as specimenagemax,
	(array_agg(case when sage.specimengeolageprefix is not null then sage.specimengeolageprefix || '-' || sage.specimengeolage else sage.specimengeolage end)) as geologicalage,
	(array_agg(sage.eruptionday || '.' || sage.eruptionmonth || '.' || sage.eruptionyear)) as eruptiondate,
	(array_remove(array_agg(ann_alt.annotationtext), null))[1] as alteration,
	(array_remove(array_agg(ann_samptech.annotationtext), null))[1] as samp_technique,
	(array_remove(array_agg(ann_ddmin.annotationtext), null))[1] as dd_min,
	(array_remove(array_agg(ann_ddmax.annotationtext), null))[1] as dd_max,
	jsonb_agg(distinct scd.*) as references,
	(jsonb_agg(distinct batchdata)) as batchdata
	from 
	(
		select * 
		from odm2.samplingfeatures s
		where s.samplingfeatureid = any ($1)
	) samples
	left join odm2.specimens spec on spec.samplingfeatureid = samples.samplingfeatureid
	left join odm2.sampletaxonomicclassifiers stc on stc.samplingfeatureid = samples.samplingfeatureid
	left join odm2.specimenages sage on sage.samplingfeatureid = samples.samplingfeatureid
	left join odm2.samplerelations sann on sann.sampleid = samples.samplingfeatureid
	left join odm2.annotations ann_alt on ann_alt.annotationid = sann.annotationid and ann_alt.annotationcode = 'g_samples.alteration'
	left join odm2.annotations ann_samptech on ann_samptech.annotationid = sann.annotationid and ann_samptech.annotationcode = 'g_samples.samp_technique'
	left join odm2.annotations ann_ddmin on ann_ddmin.annotationid = sann.annotationid and ann_ddmin.annotationcode = 'g_samples.drill_depth_min'
	left join odm2.annotations ann_ddmax on ann_ddmax.annotationid = sann.annotationid and ann_ddmax.annotationcode = 'g_samples.drill_depth_max'
	left join 
	(
		select batches.batchid,
		(array_agg(batches.samplingfeaturename))[1] as batchname,
		(array_agg(batches.sampleid))[1] as sampleid,
		(array_agg(spec_batch.specimentexture))[1] as crystal,
		(array_agg(spec_batch.specimenmediumcv))[1] as specimenmedium,
		array_remove(array_agg(distinct (case when tax_min.taxonomicclassifierid is not null then jsonb_build_object('id', tax_min.taxonomicclassifierid, 'value', tax_min.taxonomicclassifiername, 'label', tax_min.taxonomicclassifiercommonname) end)), null) as minerals,
		array_remove(array_agg(distinct (case when tax_hmin.taxonomicclassifierid is not null then jsonb_build_object('id', tax_hmin.taxonomicclassifierid, 'value', tax_hmin.taxonomicclassifiername, 'label', tax_hmin.taxonomicclassifiercommonname) end)), null) as hostMinerals,
		array_remove(array_agg(distinct (case when tax_imin.taxonomicclassifierid is not null then jsonb_build_object('id', tax_imin.taxonomicclassifierid, 'value', tax_imin.taxonomicclassifiername, 'label', tax_imin.taxonomicclassifiercommonname) end)), null) as inclusionMinerals,
		(array_remove(array_agg(distinct ann_mat.annotationtext), null))[1] as material,
		array_remove(array_agg(distinct ann_inc_type.annotationtext), null) as inclusionTypes,
		(array_remove(array_agg(distinct ann_rocinc.annotationtext), null))[1] as rimOrCoreInclusion,
		(array_remove(array_agg(distinct ann_rocmin.annotationtext), null))[1] as rimOrCoreMineral,
		jsonb_agg(distinct batch_res.*) as results
		from
		(
			select sr.sampleid,
			sr.batch as batchid,
			s2.samplingfeaturename
			from odm2.samplerelations sr
			left join odm2.samplingfeatures s2 on s2.samplingfeatureid = sr.batch
			where sr.sampleid = any ($1)
		) batches 
		left join odm2.specimens spec_batch on spec_batch.samplingfeatureid = batches.batchid
		left join
		(
			select stc.samplingfeatureid, tax_min.taxonomicclassifierid, tax_min.taxonomicclassifiername, tax_min.taxonomicclassifiercommonname 
			from odm2.specimentaxonomicclassifiers stc
			left join odm2.taxonomicclassifiers tax_min on tax_min.taxonomicclassifierid = stc.taxonomicclassifierid and tax_min.taxonomicclassifiertypecv = 'Mineral'
			where tax_min.taxonomicclassifierid is not null
		) tax_min on tax_min.samplingfeatureid = spec_batch.samplingfeatureid
		left join
		(
			select stc.samplingfeatureid, tax_hmin.taxonomicclassifierid,  tax_hmin.taxonomicclassifiername, tax_hmin.taxonomicclassifiercommonname 
			from odm2.specimentaxonomicclassifiers stc
			left join odm2.taxonomicclassifiers tax_hmin on tax_hmin.taxonomicclassifierid = stc.taxonomicclassifierid
			where tax_hmin.taxonomicclassifierid is not null
			and stc.specimentaxonomicclassifiertype = 'host mineral'
		) tax_hmin on tax_hmin.samplingfeatureid = spec_batch.samplingfeatureid
		left join
		(
			select stc.samplingfeatureid, tax_imin.taxonomicclassifierid, tax_imin.taxonomicclassifiername, tax_imin.taxonomicclassifiercommonname  
			from odm2.specimentaxonomicclassifiers stc
			left join odm2.taxonomicclassifiers tax_imin on tax_imin.taxonomicclassifierid = stc.taxonomicclassifierid
			where tax_imin.taxonomicclassifierid is not null
			and stc.specimentaxonomicclassifiertype = 'mineral inclusion'
		) tax_imin on tax_imin.samplingfeatureid = spec_batch.samplingfeatureid
		left join odm2.specimenages sage on sage.samplingfeatureid = spec_batch.samplingfeatureid 
		left join odm2.samplingfeatureannotations sann_batch on sann_batch.samplingfeatureid = spec_batch.samplingfeatureid
		left join odm2.annotations ann_mat on ann_mat.annotationid = sann_batch.annotationid and ann_mat.annotationcode = 'g_batches.material'
		left join odm2.annotations ann_inc_type on ann_inc_type.annotationid = sann_batch.annotationid and ann_inc_type.annotationcode = 'g_inclusions.inclusion_type'
		left join odm2.annotations ann_rocmin on ann_rocmin.annotationid = sann_batch.annotationid and ann_rocmin.annotationcode = 'g_minerals.rim_or_core_min'
		left join odm2.annotations ann_rocinc on ann_rocinc.annotationid = sann_batch.annotationid and ann_rocinc.annotationcode = 'g_inclusions.rim_or_core_inc'
		left join 
		(
			-- batch results
			select mv.samplingfeatureid,
			mv.sampledmediumcv as medium,
			mv.valuecount,
			mv.variabletypecode as itemgroup,
			mv.variablecode as itemname,
			mv.datavalue as value,
			mv.unitgeoroc as unit,
			mv.methodcode as method,
			std.standardname,
			std.standardvalue,
			std.standardvariable
			from odm2.samplerelations sr
			left join odm2.measuredvalues mv on mv.samplingfeatureid = sr.batch
			left join odm2.featureactions fa on fa.samplingfeatureid = mv.samplingfeatureid
			left join odm2.standards std on std.actionid = fa.actionid
			where sr.sampleid = any($1)
		) batch_res on batch_res.samplingfeatureid = batches.batchid
		group by batches.batchid
	) batchdata on batchdata.sampleid = samples.samplingfeatureid
	left join odm2.samplecitationdata scd on scd.samplingfeatureid = samples.samplingfeatureid
	where samples.samplingfeaturedescription = 'Sample' 
	and samples.samplingfeatureid = any ($1)
	group by samples.samplingfeatureid
) samples
left join
(
	-- query locations - Until geolocations is refactored, we get multiple outputs here
	select rel_loc.samplingfeatureid, 
	(array_agg(distinct si_loc.*)) as loc_data, 
	array_remove(array_agg(distinct sg_loc.locationname), null) as loc_names, 
	array_remove(array_agg(distinct g_loc.geolocationtype), null) as loc_types ,
	array_remove(array_agg(loc.elevationprecisioncomment), null) as elevation,
	array_remove(array_agg(distinct si_loc.sitedescription), null) as land_or_sea,
	array_remove(array_agg(distinct gs.settingname), null) as setting,
	array_remove(array_agg(distinct a_loc.annotationtext), null) as location_comments
	from odm2.relatedfeatures rel_loc
	left join odm2.samplingfeatures loc on loc.samplingfeatureid = rel_loc.relatedfeatureid 
	left join odm2.sites si_loc on si_loc.samplingfeatureid = rel_loc.relatedfeatureid 
	left join odm2.sitegeologicalsettings sgs on sgs.samplingfeatureid = si_loc.samplingfeatureid
	left join odm2.geologicalsettings gs on gs.settingid = sgs.settingid
	left join odm2.sitegeolocations sg_loc on sg_loc.samplingfeatureid  = si_loc.samplingfeatureid
	left join odm2.geolocations g_loc on g_loc.geolocationid = sg_loc.geolocationid 
	left join odm2.samplingfeatureannotations sa_loc on sa_loc.samplingfeatureid = si_loc.samplingfeatureid
	left join odm2.annotations a_loc on a_loc.annotationid = sa_loc.annotationid and a_loc.annotationcode = 'g_locations.location_comment'
	where rel_loc.samplingfeatureid = any ($1)
	group by rel_loc.samplingfeatureid 
) loc on loc.samplingfeatureid = samples.samplingfeatureid
left join (
	-- query sample methods
	select f_meth.samplingfeatureid as id,
	array_remove(array_agg(distinct meth.methodcode), null) as method_acronyms,
	(array_remove((array_agg(distinct a_meth.actiondescription)),null)) as method_comments,
	array_agg(distinct org.organizationname) as institutions
	from odm2.featureactions f_meth
	left join odm2.actions a_meth on a_meth.actionid = f_meth.actionid
	left join odm2.methods meth on meth.methodid = a_meth.methodid
	left join odm2.actionby ab_meth on ab_meth.actionid = a_meth.actionid 
	left join odm2.organizations org on org.organizationid = ab_meth.organizationid
	where f_meth.samplingfeatureid = any ($1)
	group by f_meth.samplingfeatureid
) methods on methods.id = samples.samplingfeatureid
left join (
	-- sample results
	select res.samplingfeatureid,
	jsonb_agg(distinct res) as results
	from (
		select mv.samplingfeatureid,
		mv.sampledmediumcv as medium,
		mv.valuecount,
		mv.variabletypecode as itemgroup,
		mv.variablecode as itemname,
		mv. datavalue as value,
		mv.unitgeoroc as unit,
		mv.methodcode as method,
		std.standardname,
		std.standardvalue,
		std.standardvariable
		from odm2.samplerelations sr
		left join odm2.measuredvalues mv on mv.samplingfeatureid = sr.batch
		left join odm2.featureactions fa on fa.samplingfeatureid = mv.samplingfeatureid
		left join odm2.standards std on std.actionid = fa.actionid
		where sr.sampleid = any($1)
	)res
	group by res.samplingfeatureid
) results on results.samplingfeatureid = samples.samplingfeatureid
`
