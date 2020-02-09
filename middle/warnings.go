package middle

import (
	"github.com/CCDirectLink/ccmu/game"
	"github.com/Masterminds/semver"
)

// WarningID represents a kind of warning action.
type WarningID int

const (
	// NullActionWarningID cannot be automatically fixed.
	NullActionWarningID WarningID = iota
	// InstallOrUpdatePackageWarningID warnings can be solved by installing/updating the package Parameter.
	InstallOrUpdatePackageWarningID
)

// Warning represents a warning to show the user on the primary view.
type Warning struct {
	Text      string
	Action    WarningID
	Parameter string
}

// FindWarnings detects issues with the installation to show on the primary view.
func FindWarnings(game game.Game) []Warning {
	warnings := []Warning{}
	crosscode, err := game.Get("crosscode")
	if err != nil {
		warnings = append(warnings, Warning{
			Text: "CrossCode is not installed in this installation. (Ok, come on, how'd you manage this? - CCDirectLink)",
		})
	} else {
		info, err := crosscode.Info()

		if err != nil {
			warnings = append(warnings, Warning{
				Text: "The CrossCode version could not be read.",
			})
		} else {
			cvers, err := semver.NewVersion(info.CurrentVersion)
			if err != nil {
				warnings = append(warnings, Warning{
					Text: "The CrossCode version could not be parsed.",
				})
			}
			if cvers.LessThan(semver.MustParse("1.1.0")) {
				warnings = append(warnings, Warning{
					Text: "The CrossCode version is " + info.CurrentVersion + "; mods usually expect 1.1.0 or higher.",
				})
			}
		}
	}
	return warnings
}
