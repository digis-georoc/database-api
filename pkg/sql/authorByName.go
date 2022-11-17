package sql

const AuthorByNameQuery = `
SELECT * FROM odm2.people WHERE lower(personlastname) = lower($1)
`
