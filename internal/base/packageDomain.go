package base

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	regextVersion = regexp.MustCompile(`^v[\d]+`)
)

type PackageDomain map[string]string

func NewPackageDomainList() (pd PackageDomain) {
	pd = make(PackageDomain, 6)
	return
}

func (i PackageDomain) Add(pkg string) PackageDomain {
	pkgName, alias := i.GetPackageName(pkg)
	if alias != "" {
		i[pkgName] = alias + " " + filepath.Dir(pkg)
		return i
	}
	i[pkgName] = filepath.Dir(pkg)
	return i
}

func (i PackageDomain) GetPackageName(pkg string) (pkgName string, alias string) {
	pkgName = filepath.Base(filepath.Dir(pkg))
	if regextVersion.MatchString(pkgName) {
		pkgs := strings.Split(pkg, "/")
		l := len(pkgs)
		if l >= 3 {
			pkgName = pkgs[l-3] + pkgs[l-2]
			alias = pkgName
		}
	}
	return
}

func (i PackageDomain) Fix(pkg string) (packageName, paramType string) {
	if strings.Index(pkg, ".") == -1 {
		packageName = pkg
		paramType = pkg
		return
	}
	pkg = strings.ReplaceAll(pkg, ".", "/")
	pkgName, _ := i.GetPackageName(pkg)
	var ok bool
	if packageName, ok = i[pkgName]; ok {
		packageName = packageName
	}
	paramType = pkgName + "." + filepath.Base(pkg)
	return
}

func (i PackageDomain) ParsePackageInParam(param string) (paramType, packageName string) {
	if strings.Contains(param, ".") == false {
		paramType = param
		return
	}
	packageName, paramType = i.Fix(param)
	return
}
