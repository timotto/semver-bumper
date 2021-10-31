package estimator

import (
	"fmt"
	. "github.com/timotto/semver-bumper/pkg/config"
	. "github.com/timotto/semver-bumper/pkg/model"
	"regexp"
	"strconv"
	"strings"
)

type estimator struct {
	config *Options
}

func NewEstimator(config *Options) *estimator {
	return &estimator{
		config: config,
	}
}

func (e estimator) BumpLevelFrom(commitMessages []string) BumpLevel {
	lvl := BumpLevelNone
	switch e.config.NoMatchBumpValue() {
	case FallbackStrategyPatch:
		lvl = BumpLevelPatch
	}

	for _, message := range commitMessages {
		if containsAny(message, e.config.KeywordsMajor) {
			lvl = BumpLevelMajor
			break
		}

		if containsAny(message, e.config.KeywordsMinor) {
			if lvl < BumpLevelMinor {
				lvl = BumpLevelMinor
			}
			continue
		}

		if containsAny(message, e.config.KeywordsPatch) {
			if lvl < BumpLevelMinor {
				lvl = BumpLevelPatch
			}
			continue
		}
	}

	return lvl
}

func (e estimator) NextPrerelease(pre string) (string, error) {
	prefix := e.config.Prerelease + "."

	if pre == "" {
		return prefix + "1", nil
	}

	if !strings.HasPrefix(pre, prefix) {
		return "", fmt.Errorf("expected prerelease to start with %v but found %v", e.config.Prerelease, pre)
	}

	strVal := strings.TrimPrefix(pre, prefix)
	val, err := strconv.ParseInt(strVal, 10, 32)
	if err != nil {
		return "", fmt.Errorf("expected prerelease to be an integer but found %v", strVal)
	}

	return fmt.Sprintf("%s%d", prefix, val+1), nil
}

func containsAny(msg string, regexps []string) bool {
	for _, re := range regexps {
		if regexp.MustCompile(re).MatchString(msg) {
			return true
		}
	}

	return false
}
