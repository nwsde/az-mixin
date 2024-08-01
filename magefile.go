//go:build mage

package main

import (
	"get.porter.sh/magefiles/git"
	"get.porter.sh/magefiles/mixins"
	"get.porter.sh/magefiles/tools"
	"github.com/carolynvs/magex/shx"
	"github.com/magefile/mage/mg"
)

const (
	mixinName    = "az"
	mixinPackage = "get.porter.sh/mixin/az"
	mixinBin     = "bin/mixins/" + mixinName
)

var magefile = mixins.NewMagefile(mixinPackage, mixinName, mixinBin)
var must = shx.CommandBuilder{StopOnError: true}

func ConfigureAgent() {
	magefile.ConfigureAgent()
}

// Build the mixin
func Build() {
	magefile.Build()
}

// Cross-compile the mixin before a release
func XBuildAll() {
	magefile.XBuildAll()
}

// Run unit tests
func TestUnit() {
	magefile.TestUnit()
}

func Test() {
	magefile.Test()
}

// Publish the mixin to github
func Publish() {
	magefile.Publish()
}

// Run Go Vet on the project
func Vet() {
	must.RunV("go", "vet", "./...")
}

// Run golangci-lint on the project
func Lint() {
	mg.Deps(tools.EnsureGolangCILint)
	must.RunV("golangci-lint", "run", "--max-issues-per-linter", "0", "--max-same-issues", "0", "./...")
}

// TestPublish tries out publish locally, with your github forks
// Assumes that you forked and kept the repository name unchanged.
func TestPublish(username string) {
	magefile.TestPublish(username)
}

// Install the mixin
func Install() {
	magefile.Install()
}

// Remove generated build files
func Clean() {
	magefile.Clean()
}

// SetupDCO configures your git repository to automatically sign your commits
// to comply with our DCO
func SetupDCO() error {
	return git.SetupDCO()
}
