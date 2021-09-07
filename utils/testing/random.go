package testing

import "math/rand"

// RandomString generates a random string (alphanumeric).
func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

// RandomInt63 generates a random positive number of type int64, avoiding the given exceptions.
func RandomInt63(maxValue int64, except ...int64) int64 {
	var i int64

	for {
		i = rand.Int63n(maxValue)

		var match bool
		for _, n := range except {
			if n == i {
				match = true
			}
		}

		if !match {
			break
		}
	}

	return i
}
