package sqlcheck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/harilsatra/sqllint/sqlcheck" // Corrected import path
)

func TestSqllint(t *testing.T) {
	// The analysistest.Run function expects the testdata directory to be
	// 'testdata/src/<pkgname>/...'. Adjust <pkgname> if your package name
	// in test files is different from 'a'.
	// For this setup, we'll use 'a' as the package name in testdata.
	testdata := analysistest.TestData()
	results := analysistest.Run(t, testdata, sqlcheck.Analyzer, "a") // "a" is the package pattern
	assert.NotNil(t, results) // Basic check to ensure results are processed
	// analysistest.Run will automatically fail the test if diagnostics don't match 'want' comments.
}
