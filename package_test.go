package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestListTools(t *testing.T) {
	const good = `
	//go:build tools
	// +build tools
	
	package main
	
	import (
		_ "github.com/a-h/templ/cmd/templ"
		_ "github.com/pressly/goose/v3/cmd/goose"
		_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	)
	`

	type args struct {
		src string
	}
	tests := []struct {
		name    string
		args    args
		want    []Package
		wantErr bool
	}{
		{
			name: "lists packages",
			args: args{src: good},
			want: []Package{
				{
					Path: "github.com/a-h/templ/cmd/templ",
				},
				{
					Path: "github.com/pressly/goose/v3/cmd/goose",
				},
				{
					Path: "github.com/sqlc-dev/sqlc/cmd/sqlc",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListTools(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTools() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPopulatePackageVersions(t *testing.T) {
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

	type args struct {
		packages []Package
		modSrc   string
	}
	tests := []struct {
		name    string
		args    args
		want    []Package
		wantErr bool
	}{
		{
			name: "populates versions",
			args: args{
				packages: []Package{
					{
						Path: "github.com/a-h/templ/cmd/templ",
					},
					{
						Path: "github.com/pressly/goose/v3/cmd/goose",
					},
					{
						Path: "github.com/sqlc-dev/sqlc/cmd/sqlc",
					},
				},
				modSrc: mod,
			},
			want: []Package{
				{
					Path:    "github.com/a-h/templ/cmd/templ",
					Version: "v0.3.819",
					Module:  "github.com/a-h/templ",
				},
				{
					Path:    "github.com/pressly/goose/v3/cmd/goose",
					Version: "v3.24.1",
					Module:  "github.com/pressly/goose/v3",
				},
				{
					Path:    "github.com/sqlc-dev/sqlc/cmd/sqlc",
					Version: "v1.27.0",
					Module:  "github.com/sqlc-dev/sqlc",
				},
			},

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PopulatePackageVersions(tt.args.packages, tt.args.modSrc)
			if (err != nil) != tt.wantErr {
				t.Errorf("PopulatePackageVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PopulatePackageVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleTrie_AddModule(t *testing.T) {
	type fields struct {
		Children map[string]*ModuleTrie
		Module   string
		Version  string
	}
	type args struct {
		module  string
		version string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ModuleTrie
	}{
		{
			name: "adds module",
			fields: fields{
				Children: make(map[string]*ModuleTrie),
				Module:   "",
				Version:  "",
			},
			args: args{
				module:  "github.com/a-h/templ",
				version: "v0.3.819",
			},
			want: &ModuleTrie{
				Children: map[string]*ModuleTrie{
					"github.com": {
						Children: map[string]*ModuleTrie{
							"a-h": {
								Children: map[string]*ModuleTrie{
									"templ": {
										Children: make(map[string]*ModuleTrie),
										Module:   "github.com/a-h/templ",
										Version:  "v0.3.819",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModuleTrie{
				Children: tt.fields.Children,
				Module:   tt.fields.Module,
				Version:  tt.fields.Version,
			}
			m.Add(tt.args.module, tt.args.version)

			if !reflect.DeepEqual(m, tt.want) {
				t.Errorf(
					"AddModule() = %v, want %v",
					stringify(m),
					stringify(tt.want),
				)
			}
		})
	}
}

func stringify(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
