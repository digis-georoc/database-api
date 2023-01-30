package sql

const SitesQuery = `
SELECT * FROM odm2.sites
`

const SiteByIDQuery = `
SELECT * FROM odm2.sites s
WHERE s.samplingfeatureID = $1
`

const GeoSettingsQuery = `
SELECT distinct s.setting FROM odm2.sites s
`
const LandOrSeaQuery = `
select distinct upper(s.sitedescription) as name
from odm2.sites s 
`
