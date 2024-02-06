package sql

// "IN $1" wont work with array but "= ANY ($1)" does
const FullDataByMultiIdQuery = `
select
samples.SamplingFeatureID as sampleID,
samples.uuid as uniqueID,
refs.references,
samples.name as samplename,
coalesce (loc.loc_names, array['Unknown']) as locationNames,
coalesce (loc.loc_types, array['Unknown']) as locationTypes,
(regexp_match(loc.elevation[1], 'ELEVATION_MIN=(-?\d+\.?\d*);'))[1] as elevationMin,
(regexp_match(loc.elevation[1], 'ELEVATION_MAX=(-?\d+\.?\d*)'))[1] as elevationMax,
-- no samp_technique in odm2
(loc.land_or_sea)[1] as landOrSea,
samples.rock_type as rockType,
samples.rock_class as rockClass,
samples.specimentexture as rockTextures,
(samples.specimenagemin)[1] as ageMin,
(samples.specimenagemax)[1] as ageMax,
(samples.eruptiondate)[1] as eruptiondate,
(samples.geologicalage)[1] as geologicalage,
samples.minerals as minerals,
samples.host_mineral as hostMinerals,
loc.loc_data[1].samplingfeatureid as locationNum,
loc.loc_data[1].latitude as latitude,
loc.loc_data[1].longitude as longitude,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LATITUDE_MIN=(-?\d+\.?\d*);'))[1] as latitudeMin,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LATITUDE_MAX=(-?\d+\.?\d*);'))[1] as latitudeMax,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LONGITUDE_MIN=(-?\d+\.?\d*);'))[1] as longitudeMin,
(regexp_match(loc.loc_data[1].locationprecisioncomment, 'LONGITUDE_MAX=(-?\d+\.?\d*)'))[1] as longitudeMax,
(loc.setting)[1] as tectonicSetting,
methods.method_acronyms as method,
methods.method_comments as comment,
(methods.institution) as institutions, -- actionBy seems to be sparsely filled
results.results,
samples.alterations as alterations,
samples.samp_techniques as samplingTechniques,
samples.dd_min as drillDepthMin,
samples.dd_max as drillDepthMax,
samples.batchdata
from 
(
	select samples.samplingfeatureid, 
	(array_agg(samples.samplingfeatureuuid))[1] as uuid,
	(array_agg(samples.samplingfeaturename))[1] as name,
	(array_agg(distinct tax_type.taxonomicclassifiername))[1] as rock_type,
	(array_agg(distinct tax_class.taxonomicclassifiername))[1] as rock_class,
	(array_agg(distinct tax_min.taxonomicclassifiercommonname))[1] as minerals,
	(array_agg(distinct tax_hostmin.taxonomicclassifiercommonname))[1] as host_mineral,
	(array_remove(array_agg(distinct spec.specimentexture), null)) as specimentexture,
	(array_agg(sage.specimenagemin)) as specimenagemin,
	(array_agg(sage.specimenagemax)) as specimenagemax,
	(array_agg(case when sage.specimengeolageprefix is not null then sage.specimengeolageprefix || '-' || sage.specimengeolage else sage.specimengeolage end)) as geologicalage,
	(array_agg(sage.eruptionday || '.' || sage.eruptionmonth || '.' || sage.eruptionyear)) as eruptiondate,
	(array_remove(array_agg(distinct ann_alt.annotationtext), null)) as alterations,
	(array_remove(array_agg(distinct ann_samptech.annotationtext), null)) as samp_techniques,
	(array_remove(array_agg(ann_ddmin.annotationtext), null)) as dd_min,
	(array_remove(array_agg(ann_ddmax.annotationtext), null)) as dd_max,
	(jsonb_agg(batchdata)) as batchdata
	from 
	(
		select * 
		from odm2.samplingfeatures s
		where s.samplingfeatureid = any ($1)
	) samples
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
		left join odm2.taxonomicclassifiers tax_min on tax_min.taxonomicclassifierid = stc.taxonomicclassifierid and tax_min.taxonomicclassifiertypecv = 'Mineral'
		where tax_min.taxonomicclassifierid is not null
	) tax_min on tax_min.samplingfeatureid = spec.samplingfeatureid
	left join
	(
		select stc.samplingfeatureid, tax_hostmin.taxonomicclassifiercommonname 
		from odm2.specimentaxonomicclassifiers stc
		left join odm2.taxonomicclassifiers tax_hostmin on tax_hostmin.taxonomicclassifierid = stc.taxonomicclassifierid
		where tax_hostmin.taxonomicclassifierid is not null
		and stc.specimentaxonomicclassifiertype = 'host mineral'
	) tax_hostmin on tax_hostmin.samplingfeatureid = spec.samplingfeatureid
	left join odm2.specimenages sage on sage.samplingfeatureid = samples.samplingfeatureid
	left join odm2.samplingfeatureannotations sann on sann.samplingfeatureid = samples.samplingfeatureid
	left join odm2.annotations ann_alt on ann_alt.annotationid = sann.annotationid and ann_alt.annotationcode = 'g_samples.alteration'
	left join odm2.annotations ann_samptech on ann_samptech.annotationid = sann.annotationid and ann_samptech.annotationcode = 'g_samples.samp_technique'
	left join odm2.annotations ann_ddmin on ann_ddmin.annotationid = sann.annotationid and ann_ddmin.annotationcode = 'g_samples.drill_depth_min'
	left join odm2.annotations ann_ddmax on ann_ddmax.annotationid = sann.annotationid and ann_ddmax.annotationcode = 'g_samples.drill_depth_max'
	left join 
	(
		select batches.batchid,
		(array_agg(batches.samplingfeaturename))[1] as batchname,
		(array_agg(batches.samplingfeatureid))[1] as sampleid,
		(array_agg(spec_batch.specimentexture))[1] as crystal,
		(array_agg(spec_batch.specimenmediumcv))[1] as specimenmedium, 
		array_remove(array_agg(distinct tax_type.taxonomicclassifiername), null) as rocktypes,
		array_remove(array_agg(distinct tax_class.taxonomicclassifiername), null) as rockclasses, 
		array_remove(array_agg(distinct tax_min.taxonomicclassifiercommonname), null) as minerals ,
		array_remove(array_agg(distinct ann_mat.annotationtext), null) as materials,
		array_remove(array_agg(distinct ann_inc_type.annotationtext), null) as inclusionTypes,
		(array_remove(array_agg(distinct ann_rocinc.annotationtext), null))[1] as rimOrCoreInclusion,
		(array_remove(array_agg(distinct ann_rocmin.annotationtext), null))[1] as rimOrCoreMineral
		from
		(
			select s.samplingfeatureid,
			r.samplingfeatureid as batchid,
			s2.samplingfeaturename
			from odm2.samplingfeatures s 
			left join odm2.relatedfeatures r on r.relatedfeatureid = s.SamplingFeatureID and r.relationshiptypecv = 'Is child of'
			left join odm2.samplingfeatures s2 on s2.samplingfeatureid = r.samplingfeatureid 
			where s.samplingfeatureid = any ($1)
		) batches 
		left join odm2.specimens spec_batch on spec_batch.samplingfeatureid = batches.batchid
		left join
		(
			select stc.samplingfeatureid, tax_type.taxonomicclassifiername 
			from odm2.specimentaxonomicclassifiers stc
			left join odm2.taxonomicclassifiers tax_type on tax_type.taxonomicclassifierid = stc.taxonomicclassifierid and tax_type.taxonomicclassifiertypecv = 'Rock'
			where tax_type.taxonomicclassifierid is not null
		) tax_type on tax_type.samplingfeatureid = spec_batch.samplingfeatureid
		left join
		(
			select stc.samplingfeatureid, tax_class.taxonomicclassifiername 
			from odm2.specimentaxonomicclassifiers stc
			left join odm2.taxonomicclassifiers tax_class on tax_class.taxonomicclassifierid = stc.taxonomicclassifierid and tax_class.taxonomicclassifiertypecv = 'Lithology'
			where tax_class.taxonomicclassifierid is not null
		) tax_class on tax_class.samplingfeatureid = spec_batch.samplingfeatureid
		left join
		(
			select stc.samplingfeatureid, tax_min.taxonomicclassifiercommonname 
			from odm2.specimentaxonomicclassifiers stc
			left join odm2.taxonomicclassifiers tax_min on tax_min.taxonomicclassifierid = stc.taxonomicclassifierid and tax_min.taxonomicclassifiertypecv = 'Mineral'
			where tax_min.taxonomicclassifierid is not null
		) tax_min on tax_min.samplingfeatureid = spec_batch.samplingfeatureid
		left join odm2.specimenages sage on sage.samplingfeatureid = spec_batch.samplingfeatureid 
		left join odm2.samplingfeatureannotations sann_batch on sann_batch.samplingfeatureid = spec_batch.samplingfeatureid
		left join odm2.annotations ann_mat on ann_mat.annotationid = sann_batch.annotationid and ann_mat.annotationcode = 'g_batches.material'
		left join odm2.annotations ann_inc_type on ann_inc_type.annotationid = sann_batch.annotationid and ann_inc_type.annotationcode = 'g_inclusions.inclusion_type'
		left join odm2.annotations ann_rocmin on ann_rocmin.annotationid = sann_batch.annotationid and ann_rocmin.annotationcode = 'g_minerals.rim_or_core_min'
		left join odm2.annotations ann_rocinc on ann_rocinc.annotationid = sann_batch.annotationid and ann_rocinc.annotationcode = 'g_inclusions.rim_or_core_inc'
		group by batches.batchid
	) batchdata on batchdata.sampleid = samples.samplingfeatureid
	where samples.samplingfeaturedescription = 'Sample' 
	and samples.samplingfeatureid = any ($1)
	group by samples.samplingfeatureid 
) samples
left join 
(	
	select q.samplingfeatureid, array_agg(q.reference) as references from (
		-- query references and authors
		select stc_ref.samplingfeatureid,
		json_build_object('citationid', c_ref.citationid , 'title', c_ref.title , 'publisher', c_ref.publisher, 'publicationyear', c_ref.publicationyear , 'link', c_ref.citationlink, 'journal', c_ref.journal , 'volume', c_ref.volume, 'issue', c_ref.issue, 'firstpage', c_ref.firstpage, 'lastpage', c_ref.lastpage , 'booktitle', c_ref.booktitle, 'editors', c_ref.editors, 'authors', array_agg(distinct authors), 'doi', cei_ref.citationexternalidentifier) as reference 
		from odm2.specimentaxonomicclassifiers stc_ref
		left join odm2.citations c_ref on c_ref.citationid = stc_ref.citationid
		left join odm2.citationexternalidentifiers cei_ref on cei_ref.citationid = c_ref.citationid and cei_ref.externalidentifiersystemid = 1 -- id of externalidentifiersystem "DOI"
		left join 
		(
			select distinct p_ref.personid,
			a_ref.citationid,
			a_ref.authororder,
			p_ref.personfirstname,
			p_ref.personlastname 
			from odm2.authorlists a_ref
			left join odm2.people p_ref on p_ref.personid = a_ref.personid
		) authors on authors.citationid = c_ref.citationid
		where stc_ref.samplingfeatureid = any ($1)
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
	array_remove(array_agg(distinct gs.settingname), null) as setting
	from odm2.relatedfeatures rel_loc
	left join odm2.samplingfeatures loc on loc.samplingfeatureid = rel_loc.relatedfeatureid 
	left join odm2.sites si_loc on si_loc.samplingfeatureid = rel_loc.relatedfeatureid 
	left join odm2.sitegeologicalsettings sgs on sgs.samplingfeatureid = si_loc.samplingfeatureid
	left join odm2.geologicalsettings gs on gs.settingid = sgs.settingid
	left join odm2.sitegeolocations sg_loc on sg_loc.samplingfeatureid  = si_loc.samplingfeatureid
	left join odm2.geolocations g_loc on g_loc.geolocationid = sg_loc.geolocationid 
	where rel_loc.samplingfeatureid = any ($1)
	group by rel_loc.samplingfeatureid 
) loc on loc.samplingfeatureid = samples.samplingfeatureid
left join (
	-- query methods
	select rel_meth.relatedfeatureid as id,
	(array_agg(distinct meth.methodcode)) as method_acronyms,
	(array_remove((array_agg(distinct a_meth.actiondescription)),null)) as method_comments,
	array_agg(distinct org.organizationname) as institution
	from odm2.relatedfeatures rel_meth
	left join odm2.featureactions f_meth on f_meth.samplingfeatureid = rel_meth.samplingfeatureid
	left join odm2.actions a_meth on a_meth.actionid = f_meth.actionid
	left join odm2.methods meth on meth.methodid = a_meth.methodid
	left join odm2.actionby ab_meth on ab_meth.actionid = a_meth.actionid 
	left join odm2.organizations org on org.organizationid = ab_meth.organizationid
	where rel_meth.relatedfeatureid = any ($1)
	group by rel_meth.relatedfeatureid
) methods on methods.id = samples.samplingfeatureid
left join (
	-- query results
	select res.sampleID,
	json_agg(res) as results
	from (
		select rel_res.relatedfeatureid as sampleID,
		mv.variablecode as itemName,
		mv.variabletypecode as itemGroup,
		std.std_names as standardName,
		std.std_values as standardValue,
		mv.datavalue as value,
		mv.unitgeoroc as unit
		from odm2.relatedfeatures rel_res
		join odm2.measuredvalues mv on mv.samplingfeatureid = rel_res.samplingfeatureid
		left join 
		(
			select relf.samplingfeatureid,
			std.standardname as std_names,
			std.standardvalue as std_values,
			std.standardvariable as std_var
			from odm2.relatedfeatures relf
			join
			(	
				select fa.samplingfeatureid,
				standards.standardname,
				standards.standardvalue,
				standards.standardvariable
				from odm2.featureactions fa 
				join odm2.standards standards on standards.actionid = fa.actionid
			) std on std.samplingfeatureid = relf.samplingfeatureid
			where relf.relatedfeatureid = any ($1)
		)std on std.samplingfeatureid = rel_res.samplingfeatureid and std.std_var = mv.variablecode
		where rel_res.relatedfeatureid = any ($1)
		and rel_res.relationshiptypecv = 'Is child of'
	)res
	group by res.sampleID
) results on results.sampleID = samples.samplingfeatureid
`
