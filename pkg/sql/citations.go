// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

const CitationsQuery = `
SELECT c_ref.citationid,
c_ref.title,
c_ref.publisher,
c_ref.publicationyear,
c_ref.citationlink,
c_ref.journal,
c_ref.volume,
c_ref.issue,
c_ref.firstpage,
c_ref.lastpage,
c_ref.booktitle,
c_ref.editors,
array_agg(json_build_object('personid', p_ref.personid, 'personfirstname', p_ref.personfirstname, 'personlastname', p_ref.personlastname)) as authors,
(array_agg(distinct cei_ref.citationexternalidentifier))[1] as doi
FROM odm2.citations c_ref
left join odm2.citationexternalidentifiers cei_ref on cei_ref.citationid = c_ref.citationid and cei_ref.externalidentifiersystemid = 1 -- id of externalidentifiersystem "DOI"
left join odm2.authorlists a_ref on a_ref.citationid = c_ref.citationid
left join odm2.people p_ref on p_ref.personid = a_ref.personid 
group by c_ref.citationid
`

const CitationByIDQuery = `
SELECT c_ref.citationid,
c_ref.title,
c_ref.publisher,
c_ref.publicationyear,
c_ref.citationlink,
c_ref.journal,
c_ref.volume,
c_ref.issue,
c_ref.firstpage,
c_ref.lastpage,
c_ref.booktitle,
c_ref.editors,
array_agg(json_build_object('personid', p_ref.personid, 'personfirstname', p_ref.personfirstname, 'personlastname', p_ref.personlastname)) as authors,
(array_agg(distinct cei_ref.citationexternalidentifier))[1] as doi
FROM odm2.citations c_ref
left join odm2.citationexternalidentifiers cei_ref on cei_ref.citationid = c_ref.citationid and cei_ref.externalidentifiersystemid = 1 -- id of externalidentifiersystem "DOI"
left join odm2.authorlists a_ref on a_ref.citationid = c_ref.citationid
left join odm2.people p_ref on p_ref.personid = a_ref.personid 
WHERE c_ref.citationID = $1
group by c_ref.citationid
`
