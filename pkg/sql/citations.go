package sql

const CitationsQuery = `
SELECT * FROM odm2.citations
`

const CitationByIDQuery = `
SELECT * FROM odm2.citations c
WHERE c.citationID = $1
`
