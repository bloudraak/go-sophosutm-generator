package main

type PackageImport struct {
    importedPackage string
    importedClass   *Class
}

func (p *PackageImport) ImportedPackage() string {
    return p.importedPackage
}

func (p *PackageImport) SetImportedPackage(importedPackage string) {
    p.importedPackage = importedPackage
}

func (p *PackageImport) WithImportedPackage(importedPackage string) *PackageImport {
    p.importedPackage = importedPackage
    return p
}

func (p *PackageImport) ImportedClass() *Class {
    return p.importedClass
}

func NewPackageImport(class *Class, pkg string) *PackageImport {
    return &PackageImport{importedPackage: pkg, importedClass: class}
}
