package strings

import (
	"strings"
)

// EscapeQualifiedName converts a plugin name, which might contain a / into a
// string that is safe to use on-disk.  This assumes that the input has already
// been validates as a qualified name.  we use "~" rather than ":" here in case
// we ever use a filesystem that doesn't allow ":".
func EscapeQualifiedName(in string) string {
	return strings.Replace(in, "/", "~", -1)
}

// UnescapeQualifiedName converts an escaped plugin name (as per EscapeQualifiedName)
// back to its normal form.  This assumes that the input has already been
// validates as a qualified name.
func UnescapeQualifiedName(in string) string {
	return strings.Replace(in, "~", "/", -1)
}
