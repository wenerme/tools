package apk

import "strings"

type DependencyType string

const (
	DependencyPackage DependencyType = ""
	DependencyPkgConf DependencyType = "pc"
	DependencyCommand DependencyType = "cmd"
	DependencySo      DependencyType = "so"
)

type Dependency struct {
	Type    string // Type may set to pc,cmd,so if not package
	Name    string
	Version string
}

func (d Dependency) String() string {
	s := d.Name
	if d.Type != "" {
		s = d.Type + ":" + s
	}
	if d.Version != "" {
		s += "=" + d.Version
	}
	return s
}
func ParseDependency(dep string) Dependency {
	d := Dependency{}
	s := dep
	{
		sp := strings.SplitN(s, ":", 2)
		if len(sp) > 1 {
			d.Type = sp[0]
			s = sp[1]
		}
	}
	{
		sp := strings.SplitN(s, "=", 2)
		d.Name = sp[0]
		if len(sp) > 1 {
			d.Version = sp[1]
		}
	}
	return d
}
