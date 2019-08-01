package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperationsLoneAnonymousOperation(t *testing.T) {
	assert.Empty(t, validateSource(t, `{scalar}`))
	assert.NotEmpty(t, validateSource(t, `{scalar} {scalar}`))
}
