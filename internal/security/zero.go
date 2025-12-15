package security

// ZeroBytes overwrites a byte slice with zeros.
// Safe to call with nil or empty slices.
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
