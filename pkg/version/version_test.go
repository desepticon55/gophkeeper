package version

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestPrintBuildInfo(t *testing.T) {
	for _, tc := range []struct {
		name     string
		version  string
		date     string
		commit   string
		expected string
	}{
		{
			name:    "Simple test",
			version: "0.0.1", date: "2022/07/19", commit: "01234567",
			expected: "Build version: 0.0.1\nBuild date: 2022/07/19\nCommit: 01234567",
		},
		{
			name:     "Empty fields",
			expected: "Build version: N/A\nBuild date: N/A\nCommit: N/A",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			BuildVersion = tc.version
			BuildDate = tc.date
			BuildCommit = tc.commit

			info := MakeBuildInfo(zap.NewNop())
			assert.Equal(t, tc.expected, info)
		})
	}
}
