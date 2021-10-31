package estimator

//type (
//	Configuration struct {
//		Fallback string
//		fallback FallbackStrategy
//		Keywords KeywordConfiguration
//
//		PrereleasePrefix string
//	}
//
//	KeywordConfiguration struct {
//		Major []string
//		Minor []string
//		Patch []string
//	}
//
//	FallbackStrategy int
//)
//
//const (
//	FallbackStrategyNone FallbackStrategy = iota
//	FallbackStrategyPatch
//
//	fallbackStrategyNone  = "none"
//	fallbackStrategyPatch = "patch"
//)
//
//var (
//	defaultKeywordsMajor = []string{"^BREAKING CHANGE:"}
//	defaultKeywordsMinor = []string{"^feat:"}
//	defaultKeywordsPatch = []string{"^fix:", "^chore:"}
//)
//
//func (c *Configuration) Valid() error {
//	c.applyDefaults()
//
//	return c.validate()
//}
//
//func (c *Configuration) FallbackStrategy() FallbackStrategy {
//	return c.fallback
//}
//
//func (c *Configuration) applyDefaults() {
//	if len(c.Keywords.Major) == 0 {
//		c.Keywords.Major = defaultKeywordsMajor
//	}
//
//	if len(c.Keywords.Minor) == 0 {
//		c.Keywords.Minor = defaultKeywordsMinor
//	}
//
//	if len(c.Keywords.Patch) == 0 {
//		c.Keywords.Patch = defaultKeywordsPatch
//	}
//
//	if c.Fallback == "" {
//		c.Fallback = fallbackStrategyNone
//	}
//
//	if c.PrereleasePrefix == "" {
//		c.PrereleasePrefix = "rc"
//	}
//}
//
//func (c *Configuration) validate() error {
//	switch c.Fallback {
//	case fallbackStrategyPatch:
//		c.fallback = FallbackStrategyPatch
//
//	case fallbackStrategyNone:
//		c.fallback = FallbackStrategyNone
//
//	default:
//		return fmt.Errorf("invalid fallback value: %v", c.Fallback)
//	}
//
//	return nil
//}
