package sql

const GeoJSONQuery = `
SELECT s.latitude,
s.longitude,
(array_agg(s.samplingfeatureid))[1] AS locationID,
count(DISTINCT r.samplingfeatureid) AS numSamplingFeatureIDs,
array_agg(DISTINCT r.samplingfeatureid) AS samplingFeatureIDs,
(array_agg(DISTINCT gs.settingname))[1] AS setting,
(array_agg(DISTINCT toplevelloc.locationname))[1] AS loc1,
(array_agg(DISTINCT secondlevelloc.locationname))[1] AS loc2,
(array_agg(DISTINCT thirdlevelloc.locationname))[1] AS loc3,
(array_agg(DISTINCT s.sitedescription))[1] AS landOrSea
FROM odm2.sites s
LEFT JOIN odm2.sitegeologicalsettings sgs on sgs.samplingfeatureid = s.samplingfeatureid
LEFT JOIN odm2.geologicalsettings gs on gs.settingid = sgs.settingid
LEFT JOIN ( 
	SELECT sg.samplingfeatureid,
	sg.locationname
	FROM odm2.sitegeolocations sg
	LEFT JOIN odm2.geolocations g ON g.geolocationid = sg.geolocationid
WHERE "right"(g.locationhierarchy::character varying::text, 3) = '100'::text
) toplevelloc ON toplevelloc.samplingfeatureid = s.samplingfeatureid
LEFT JOIN ( 
	SELECT sg.samplingfeatureid,
	sg.locationname
	FROM odm2.sitegeolocations sg
	LEFT JOIN odm2.geolocations g ON g.geolocationid = sg.geolocationid
WHERE "right"(g.locationhierarchy::character varying::text, 3) = '200'::text
) secondlevelloc ON secondlevelloc.samplingfeatureid = s.samplingfeatureid
LEFT JOIN ( 
	SELECT sg.samplingfeatureid,
	sg.locationname
	FROM odm2.sitegeolocations sg
	LEFT JOIN odm2.geolocations g ON g.geolocationid = sg.geolocationid
	WHERE "right"(g.locationhierarchy::character varying::text, 3) = '300'::text
) thirdlevelloc ON thirdlevelloc.samplingfeatureid = s.samplingfeatureid
LEFT JOIN odm2.relatedfeatures r ON r.relatedfeatureid = s.samplingfeatureid
GROUP BY s.latitude, s.longitude
`
