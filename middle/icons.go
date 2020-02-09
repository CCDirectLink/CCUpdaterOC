package middle

import (
	"github.com/CCDirectLink/CCUpdaterOC/design"
	"github.com/CCDirectLink/ccmu/pkg"
)

func PackageIcon(p pkg.Package) design.IconID {
	info, err := p.Info()
	if err != nil {
		return design.ModIconID
	}

	switch info.Name {
	case "crosscode",
		"ccloader":
		return design.ToolIconID
	default:
		return design.ModIconID
	}
}