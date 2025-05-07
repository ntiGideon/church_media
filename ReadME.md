
## Performing A Detailed Migration
go run -mod=mod entgo.io/ent/cmd/ent generate ./schema

## Creating a new Entity
go run -mod=mod entgo.io/ent/cmd/ent new LogAudit