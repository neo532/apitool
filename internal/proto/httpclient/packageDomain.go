package httpclient

import (
	"path/filepath"
	"strings"
)

type PackageDomain map[string]string

func NewPackageDomainList() (pd PackageDomain) {
	pd = make(PackageDomain, 6)
	return
}

func (i PackageDomain) Add(pkg string) PackageDomain {
	i[filepath.Base(filepath.Dir(pkg))] = filepath.Dir(pkg)
	return i
}

func (i PackageDomain) Fix(pkg string) (packageName, paramType string) {
	if strings.Index(pkg, ".") == -1 {
		packageName = pkg
		paramType = pkg
		return
	}
	pkg = strings.ReplaceAll(pkg, ".", "/")
	pkgName := filepath.Base(filepath.Dir(pkg))
	packageName, _ = i[pkgName]
	paramType = pkgName + "." + filepath.Base(pkg)
	return
}
