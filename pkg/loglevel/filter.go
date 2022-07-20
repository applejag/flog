package loglevel

type Filter struct {
	MinLevel      Level
	MaxLevel      Level
	BlacklistMask Level
	WhitelistMask Level
}
