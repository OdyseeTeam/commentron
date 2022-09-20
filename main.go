package main

import "github.com/OdyseeTeam/commentron/cmd"

//go:generate go-bindata -o migration/bindata.go -nometadata -pkg migration -ignore bindata.go migration/
//go:generate go fmt ./migration/bindata.go
//go:generate goimports -l ./migration/bindata.go

func main() {
	cmd.Execute()
}
