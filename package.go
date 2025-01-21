package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"

	"golang.org/x/mod/modfile"
)

type Package struct {
	Path    string
	Version string
	Module  string
}

func (p Package) String() string {
	if p.Version != "" {
		return p.Path + "@" + p.Version
	}
	return p.Path
}

func (p Package) CommandName() string {
	return path.Base(p.Path)
}

func (p Package) BinName() string {
	return p.CommandName() + "@" + p.Version
}

func ListPackages(toolsSrc, modSrc string) ([]Package, error) {
	packages, err := ListTools(toolsSrc)
	if err != nil {
		return nil, err
	}

	packages, err = PopulatePackageVersions(packages, modSrc)
	if err != nil {
		return nil, err
	}

	return packages, nil
}

func ListTools(toolsSrc string) ([]Package, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "tools.go", toolsSrc, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	packages := make([]Package, 0)
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ImportSpec:
			path := n.Path.Value
			// Remove quotes
			path = path[1 : len(path)-1]
			packages = append(packages, Package{Path: path})
		}
		return true
	})

	return packages, nil
}

func PopulatePackageVersions(packages []Package, modSrc string) ([]Package, error) {
	f, err := modfile.ParseLax("go.mod", []byte(modSrc), nil)
	if err != nil {
		return nil, err
	}
	modules := NewModuleTrie()
	for _, r := range f.Require {
		if !r.Indirect {
			modules.Add(r.Mod.Path, r.Mod.Version)
		}
	}

	for i := range packages {
		p := &packages[i]

		module, version := modules.DeepestMatch(p.Path)

		if module == "" {
			return nil, fmt.Errorf(
				"could not find module for package %s in go.mod, did you forget to `go mod tidy`?",
				p.Path,
			)
		}

		p.Module = module
		p.Version = version
	}

	return packages, nil
}

type ModuleTrie struct {
	Children map[string]*ModuleTrie
	Module   string
	Version  string
}

func NewModuleTrie() *ModuleTrie {
	return &ModuleTrie{Children: make(map[string]*ModuleTrie)}
}

func (m *ModuleTrie) Add(module, version string) {
	components := strings.Split(module, "/")
	current := m

	for _, c := range components {
		if current.Children[c] == nil {
			current.Children[c] = NewModuleTrie()
		}
		current = current.Children[c]
	}

	current.Module = module
	current.Version = version
}

func (m *ModuleTrie) DeepestMatch(path string) (string, string) {
	components := strings.Split(path, "/")
	current := m

	for _, c := range components {
		if current.Children[c] == nil {
			break
		}
		current = current.Children[c]
	}

	return current.Module, current.Version
}
