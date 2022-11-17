# database-api

API to access the GEOROC database in the new ODM2 schema

## Documentation

The api documentation is generated with [swagger](https://github.com/swaggo/swag).
For installation guides see [the documentation](https://github.com/swaggo/swag#getting-started). You may need to add the `GOPATH` to your `PATH` variable to be able to exectue `swag` from your commandline: `export PATH=$(go env GOPATH)/bin:$PATH`.

To generate the documentation files under **cmd/docs/**, move to the **cmd/** directory where the main.go is located and execute

`swag init -g ../pkg/api/api.go`
