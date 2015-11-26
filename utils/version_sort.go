package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// ByVersion sorting type
type ByVersion []string

func (s ByVersion) Len() int {
	return len(s)
}
func (s ByVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByVersion) Less(i, j int) bool {
	if s[i] == s[j] {
		return false
	}
	vA := VersionArray(s[i])
	vB := VersionArray(s[j])
	lenA, lenB := len(vA), len(vB)
	var subVersionA int
	var subVersionB int
	for i := 0; i < maxInt(lenA, lenB); i++ {
		if i < lenA {
			subVersionA = vA[i]
		} else {
			subVersionA = 0
		}
		if i < lenB {
			subVersionB = vB[i]
		} else {
			subVersionB = 0
		}
		if subVersionA == subVersionB {
			continue
		}
		return subVersionB < subVersionA
	}
	return false
}

// VersionArray split version string into array
func VersionArray(version string) []int {
	digits := regexp.MustCompile(`\D?(\d+)`)
	vsplited := strings.Split(version, ".")
	var versionArray []int
	for _, v := range vsplited {
		match := digits.FindAllStringSubmatch(v, -1)
		vint, _ := strconv.Atoi(match[0][1])
		versionArray = append(versionArray, vint)
	}
	return versionArray
}

func maxInt(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// OnlyStable remove all beta versions
func OnlyStable(s []string) []string {
	var copy []string
	for _, v := range s {
		if hasBeta(v) {
			continue
		}
		copy = append(copy, v)
	}
	return copy
}

func hasBeta(version string) bool {
	betaVersion := regexp.MustCompile(`.*(alfa|beta|rc).*`)
	match := betaVersion.FindAllStringSubmatch(version, -1)
	return len(match) > 0
}
