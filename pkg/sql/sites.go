package sql

const SitesQuery = `
SELECT * FROM odm2.sites
`

const SiteByIDQuery = `
SELECT * FROM odm2.sites s
WHERE s.samplingfeatureID = $1
`

const GeoSettingsQuery = `
SELECT distinct gs.settingname as setting
FROM odm2.geologicalsettings gs
where gs.settingname is not null
`
