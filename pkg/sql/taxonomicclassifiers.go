package sql

const RockClassQueryStart = `
select distinct (array_agg(t.taxonomicclassifiername))[1] as value,
(array_agg(t.taxonomicclassifiercommonname))[1] as label,
count(distinct s.samplingfeatureid)
from odm2.taxonomicclassifiers t
left join odm2.specimentaxonomicclassifiers s on s.taxonomicclassifierid = t.taxonomicclassifierid
join (
	select s2.samplingfeatureid
	from odm2.taxonomicclassifiers t2
	left join odm2.specimentaxonomicclassifiers s2 on s2.taxonomicclassifierid = t2.taxonomicclassifierid
	where t2.taxonomicclassifiertypecv = 'Rock'
`

const RockClassQueryMid = `
) rocktype on rocktype.samplingfeatureid = s.samplingfeatureid 
where t.taxonomicclassifiertypecv = 'Lithology'
`

const RockClassQueryEnd = `
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
