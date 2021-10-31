package config

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"reflect"
)

const (
	FallbackStrategyNone FallbackStrategy = iota
	FallbackStrategyPatch

	fallbackStrategyNone  = "none"
	fallbackStrategyPatch = "patch"
)

type FallbackStrategy int

type Options struct {
	ConfigFile  string `json:"-" yaml:"-" short:"C" long:"config-file" description:"load parameters from a JSON or YAML file"`
	Prerelease  string `json:"pre,omitempty" yaml:"prerelease,omitempty" short:"p" long:"pre" description:"bump prerelease with given keyword, eg \"rc\" for \"1.2.3-rc.4\"'"`
	TagPrefix   string `json:"tag_prefix,omitempty" yaml:"tag_prefix,omitempty" short:"t" long:"tag-prefix" description:"only detect tags matching the expression, eg \"v\" for \"v1.2.3\""`
	NoMatchBump string `json:"no_match_bump,omitempty" yaml:"no_match_bump,omitempty" short:"n" long:"no-match-bump" choice:"none" choice:"patch" description:"bump patch or nothing when no commits match"`

	Output  string `json:"output,omitempty" yaml:"output,omitempty" short:"o" long:"output" description:"write result into file, defaults to stdout"`
	Commits string `json:"commits,omitempty" yaml:"commits,omitempty" short:"c" long:"commits" description:"write commit messages into file"`

	PathInclude []string `json:"path_include,omitempty" yaml:"path_include,omitempty" short:"i" long:"path-include" description:"only detect commits at the given path, can be supplied multiple times"`
	PathExclude []string `json:"path_exclude,omitempty" yaml:"path_exclude,omitempty" short:"x" long:"path-exclude" description:"ignore commits at the given path, can be supplied multiple times"`

	InitialVersion string   `json:"initial_version,omitempty" yaml:"initial_version,omitempty" short:"0" long:"initial-version" description:"release version if there are no tags yet, defaults to \"1.0.0\""`
	KeywordsMajor  []string `json:"keywords_major,omitempty" yaml:"keywords_major,omitempty" short:"1" long:"major" description:"commit message keywords justifying a major version bump, can be supplied multiple times"`
	KeywordsMinor  []string `json:"keywords_minor,omitempty" yaml:"keywords_minor,omitempty" short:"2" long:"minor" description:"commit message keywords justifying a minor version bump, can be supplied multiple times"`
	KeywordsPatch  []string `json:"keywords_patch,omitempty" yaml:"keywords_patch,omitempty" short:"3" long:"patch" description:"commit message keywords justifying a patch version bump, can be supplied multiple times"`

	PrintKeywords bool   `json:"-" yaml:"-" short:"k" long:"print-keywords" description:"print the configured version bump keywords and exit"`
	WriteConfig   string `json:"-" yaml:"-" short:"W" long:"write-config" description:"write the given parameters into a JSON or YAML config file and exit"`

	initialVersion *semver.Version
	noMatchBump    FallbackStrategy
}

func (o *Options) InitialVersionValue() *semver.Version {
	return o.initialVersion
}

func (o *Options) ReadConfigFile() bool {
	return o.ConfigFile != ""
}

func (o *Options) WriteConfigFile() bool {
	return o.WriteConfig != ""
}

func (o *Options) BumpPrerelease() bool {
	return o.Prerelease != ""
}

func (o *Options) NoMatchBumpValue() FallbackStrategy {
	return o.noMatchBump
}

func (o *Options) OutputCommits() bool {
	return o.Commits != ""
}

func (o *Options) SetMissingFrom(other *Options) {
	higher := reflect.ValueOf(o).Elem()
	lower := reflect.ValueOf(other).Elem()

	n := higher.NumField()
	for i := 0; i < n; i++ {
		fieldHigher := higher.Field(i)
		if !fieldHigher.IsZero() {
			continue
		}

		fieldLower := lower.Field(i)
		if fieldLower.IsZero() {
			continue
		}

		fieldHigher.Set(fieldLower)
	}
}

func (o *Options) Valid() error {
	if o.InitialVersion == "" {
		o.InitialVersion = "1.0.0"
	}
	if o.NoMatchBump == "" {
		o.NoMatchBump = fallbackStrategyNone
	}
	if len(o.KeywordsMajor) == 0 {
		o.KeywordsMajor = []string{"^BREAKING CHANGE:"}
	}
	if len(o.KeywordsMinor) == 0 {
		o.KeywordsMinor = []string{"^feat:"}
	}
	if len(o.KeywordsPatch) == 0 {
		o.KeywordsPatch = []string{"^fix:", "^chore:"}
	}

	var err error
	if o.initialVersion, err = semver.StrictNewVersion(o.InitialVersion); err != nil {
		return fmt.Errorf("invalid initial version %v: %w", o.InitialVersion, err)
	}

	switch o.NoMatchBump {
	case fallbackStrategyNone:
		o.noMatchBump = FallbackStrategyNone
	case fallbackStrategyPatch:
		o.noMatchBump = FallbackStrategyPatch
	default:
		return fmt.Errorf("invalid no match bump value: %v", o.NoMatchBump)
	}

	return nil
}
