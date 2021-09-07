package safe

import (
	"fmt"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestSafeParseUint(t *testing.T) {
	assert := assert2.New(t)

	value := uint(8888)
	parsedValue := SafeParseUint(fmt.Sprintf("%d", value))

	assert.Equal(value, parsedValue, "the value should be parsed without problems")

	value2 := "18446744073709551616" // equals to math.MaxUint64 + 1
	parsedValue = SafeParseUint(value2)

	assert.Equal(uint(DefaultAllocate), parsedValue, "the value should not be parsed (overflow), returning the default")
}

func TestSafeParseInt(t *testing.T) {
	assert := assert2.New(t)

	value := int(8888)
	parsedValue := SafeParseInt(fmt.Sprintf("%d", value))

	assert.Equal(value, parsedValue, "the value should be parsed without problems")

	value2 := "9223372036854775808" // equals to math.MaxInt64 + 1
	parsedValue = SafeParseInt(value2)

	assert.Equal(int(DefaultAllocate), parsedValue, "the value should not be parsed (overflow), returning the default")
}
