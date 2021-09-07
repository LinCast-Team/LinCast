package safe

import (
	"math"
	"strconv"
)

const DefaultAllocate = 0

// SafeParseUint does the parsing of a string into an uint securely, avoiding possible overflows.
func SafeParseUint(value string) uint {
	parsed, err := strconv.ParseUint(value, 10, strconv.IntSize)
	if err != nil {
		return DefaultAllocate
	}

	var max uint64

	if strconv.IntSize == 64 { // uint of 64 bits
		max = math.MaxUint64
	} else { // uint of 32 bits
		max = math.MaxUint32
	}

	// GOOD: check for lower and upper bounds
	if parsed > 0 && parsed <= max {
		return uint(parsed)
	}

	return DefaultAllocate
}

// SafeParseInt does the parsing of a string into an int securely, avoiding possible overflows.
func SafeParseInt(value string) int {
	parsed, err := strconv.ParseInt(value, 10, strconv.IntSize)
	if err != nil {
		return DefaultAllocate
	}

	var max int64

	if strconv.IntSize == 64 { // int of 64 bits
		max = math.MaxInt64
	} else { // int of 32 bits
		max = math.MaxInt32
	}

	// GOOD: check for lower and upper bounds
	if parsed > 0 && parsed <= max {
		return int(parsed)
	}

	return DefaultAllocate
}
