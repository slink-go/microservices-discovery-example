go build -ldflags "-w -s" -o build/bulk   bulk/main.go
go build -ldflags "-w -s" -o build/single single/main.go
go build -ldflags "-w -s" -o build/foo    foo/main.go