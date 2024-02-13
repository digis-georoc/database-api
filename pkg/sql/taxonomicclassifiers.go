// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

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

const HostMatQuery = `
-- get all host materials
select t.taxonomicclassifiername as value, t.taxonomicclassifiercommonname as label, count(distinct s.samplingfeatureid)
from odm2.taxonomicclassifiers t 
left join odm2.specimentaxonomicclassifiers s on s.taxonomicclassifierid = t.taxonomicclassifierid 
where s.specimentaxonomicclassifiertype = 'host mineral'
group by t.taxonomicclassifiername, t.taxonomicclassifiercommonname 
order by count desc
`

const IncMatQuery = `
-- get all inclusion materials
select t.taxonomicclassifiername as value, t.taxonomicclassifiercommonname as label, count(distinct s.samplingfeatureid)
from odm2.taxonomicclassifiers t 
left join odm2.specimentaxonomicclassifiers s on s.taxonomicclassifierid = t.taxonomicclassifierid 
where s.specimentaxonomicclassifiertype = 'mineral inclusion'
group by t.taxonomicclassifiername, t.taxonomicclassifiercommonname 
order by count desc
`
