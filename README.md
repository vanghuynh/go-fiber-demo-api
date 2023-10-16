swag init 
air   

# Generate migration script for data versoning
atlas migrate diff --env gorm 

# Apply migration script generated
atlas schema apply --env gorm -u "postgres://admin:password@localhost:5432/fiber?sslmode=disable"

curl -sSf https://atlasgo.sh | sh
go get ariga.io/atlas-go-sdk/atlasexec
https://atlasgo.io/integrations/go-sdk

https://pkg.go.dev/ariga.io/atlas@v0.14.1#section-readme

