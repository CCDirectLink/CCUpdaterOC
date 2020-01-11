package middle

import (
	"github.com/CCDirectLink/CCUpdaterUI/design"
	"github.com/CCDirectLink/ccmu/game"
	"github.com/CCDirectLink/ccmu/pkg"
)

// FakeError should be enabled to prevent internet access by CCUpdaterUI.
const FakeError bool = false

// InternetConnectionWarning is true if the last GetRemotePackages() call actually resulted in error.
var InternetConnectionWarning bool = false

// GetRemotePackages retrieves remote packages from the server. (The CCUpdaterCLI-level cache semantics still apply.)
func GetRemotePackages() []pkg.Package {
	if FakeError {
		InternetConnectionWarning = true
		return []pkg.Package{}
	}

	pkgs, err := game.Default.Available()
	if err != nil {
		InternetConnectionWarning = true
		return []pkg.Package{}
	}

	return pkgs
}

// PackageIcon returns the relevant icon ID for a package.
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
