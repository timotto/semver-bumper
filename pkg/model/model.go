package model

type BumpLevel int

const (
	BumpLevelNone BumpLevel = iota
	BumpLevelPatch
	BumpLevelMinor
	BumpLevelMajor
)
