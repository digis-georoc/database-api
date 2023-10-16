package sql

const RockClassQuery = `
select distinct (array_agg(t.taxonomicclassifiername))[1] as value,
(array_agg(t.taxonomicclassifiercommonname))[1] as label,
count(distinct s.samplingfeatureid) as count
from odm2.taxonomicclassifiers t
left join odm2.specimentaxonomicclassifiers s on s.taxonomicclassifierid = t.taxonomicclassifierid 
where t.taxonomicclassifiertypecv = 'Lithology'
group by t.taxonomicclassifierid 
order by count desc
`

const RockTypeQuery = `
select distinct (array_agg(t.taxonomicclassifiername))[1] as value,
(array_agg(t.taxonomicclassifiercommonname))[1] as label,
count(distinct s.samplingfeatureid)
from odm2.taxonomicclassifiers t
left join odm2.specimentaxonomicclassifiers s on s.taxonomicclassifierid = t.taxonomicclassifierid 
where t.taxonomicclassifiertypecv = 'Rock'
group by t.taxonomicclassifierid 
order by count desc
`

const MineralQuery = `
select distinct (array_agg(t.taxonomicclassifiername))[1] as value,
(array_agg(t.taxonomicclassifiercommonname))[1] as label,
count(distinct s.samplingfeatureid)
from odm2.taxonomicclassifiers t
left join odm2.specimentaxonomicclassifiers s on s.taxonomicclassifierid = t.taxonomicclassifierid 
where t.taxonomicclassifierdescription = 'Mineral Classification from GEOROC'
group by t.taxonomicclassifierid 
order by count desc
`
