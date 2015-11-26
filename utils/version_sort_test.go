package utils_test

import (
	"sort"
	"testing"

	"github.com/3zcurdia/reportcopter/utils"
	"github.com/stretchr/testify/assert"
)

func TestVersionArray(t *testing.T) {
	arr := utils.VersionArray("1.2.3.4")
	assert.Equal(t, []int{1, 2, 3, 4}, arr)
}

func TestSortByVersion(t *testing.T) {
	versions := []string{"0.999", "1.0.0", "1.0.123", "1.0.123.1"}
	sort.Sort(utils.ByVersion(versions))
	assert.Equal(t, []string{"1.0.123.1", "1.0.123", "1.0.0", "0.999"}, versions)
}

func TestSortByVersionWithAflfanumeric(t *testing.T) {
	versions := []string{"release-0.999", "release-1.0.0", "release-1.0.123", "release-1.0.123.1"}
	sort.Sort(utils.ByVersion(versions))
	assert.Equal(t, []string{"release-1.0.123.1", "release-1.0.123", "release-1.0.0", "release-0.999"}, versions)
}

func TestOnlyStable(t *testing.T) {
	versions := []string{"1.2.rc2", "1.3.beta3", "0.0.1alfa", "0.1.1"}
	versions = utils.OnlyStable(versions)
	assert.Equal(t, []string{"0.1.1"}, versions)
}
