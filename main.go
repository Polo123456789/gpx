package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	_ = ctx

	const src = `
	//go:build tools
	// +build tools

	package main

	import (
		_ "github.com/a-h/templ/cmd/templ"
		_ "github.com/pressly/goose/v3/cmd/goose"
		_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	)
	`

	const mod = `
	module some-project-of-mine.com/you-dont-care
	
	go 1.23.2
	
	require (
		github.com/Polo123456789/assert v0.1.4
		github.com/a-h/templ v0.3.819
		github.com/charmbracelet/log v0.4.0
		github.com/google/uuid v1.6.0
		github.com/gorilla/sessions v1.4.0
		github.com/pressly/goose/v3 v3.24.1
		github.com/sqlc-dev/sqlc v1.27.0
		golang.org/x/crypto v0.32.0
		modernc.org/sqlite v1.34.4
	)

	require (
		cel.dev/expr v0.19.1 // indirect
		filippo.io/edwards25519 v1.1.0 // indirect
		github.com/ClickHouse/ch-go v0.63.1 // indirect
		github.com/ClickHouse/clickhouse-go/v2 v2.30.0 // indirect
		github.com/PuerkitoBio/goquery v1.10.1 // indirect
		github.com/a-h/parse v0.0.0-20240121214402-3caf7543159a // indirect
		github.com/a-h/protocol v0.0.0-20240821172110-e94e5c43897f // indirect
		github.com/andybalholm/brotli v1.1.1 // indirect
		github.com/andybalholm/cascadia v1.3.3 // indirect
	)
	`

	binPath, err := GetCachePath()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Binaries path:", binPath)

	err = CheckCachePath(binPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	packages, err := ListPackages(src, mod)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(packages)

	for _, p := range packages {
		if err := InstallPackage(ctx, p, binPath); err != nil {
			fmt.Println(err)
			return
		}
		if err := AddVersionToExecutable(p, binPath); err != nil {
			fmt.Println(err)
			return
		}
	}
}
