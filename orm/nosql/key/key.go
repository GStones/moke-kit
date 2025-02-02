package key

import (
	"errors"
	"regexp"
	"strings"
)

const (
	KeySeparator        = "/"
	PrefixSeparator     = "$"
	validatePathPattern = `^([/]([0-9a-f]+\$)*[a-zA-Z0-9_.-]+){2,}$`
	// Keys must contain 2 or more sections that begin with a /.
	// Sections must be alphanumeric with the addition of underscores, periods, or dashes.
	// Additionally, a section may contain 0 or more hexadecimal prefixes, separated from the main identifier by a $,
	// but intended use is that only the final (Base) section should contain any prefixes.
)

var (
	validatePathExp *regexp.Regexp

	ErrNotEnoughParts   = errors.New("ErrNotEnoughParts")
	ErrInvalidKeyFormat = errors.New("ErrInvalidKeyFormat")
)

func init() {
	validatePathExp = regexp.MustCompile(validatePathPattern)
}

// A Key is a structured abstraction of nosql identifier.
type Key struct {
	value string
}

// NewKey constructs a new keys from the given fully qualified string.  It will panic if the string isn't valid.  Use
// NewKeyFromString if you don't want to panic.
func NewKey(value string) Key {
	return MustSucceed(NewKeyFromString(value))
}

func NewKeys(values ...string) []Key {
	keys := make([]Key, len(values))
	for k, v := range values {
		keys[k] = NewKey(v)
	}
	return keys
}

// NewKeyFromParts creates a new keys within the current core.Namespace.
func NewKeyFromParts(parts ...string) (key Key, err error) {
	if len(parts) == 0 {
		return key, ErrNotEnoughParts
	}

	var builder strings.Builder
	// Preallocate capacity for all parts + separators
	builder.Grow(len(parts) + len(KeySeparator)*len(parts))

	for _, part := range parts {
		builder.WriteString(KeySeparator)
		builder.WriteString(part)
	}

	return NewKeyFromString(builder.String())
}

// NewKeyFromString creates a new keys within the current core.Namespace from a string
// representation, validating each part.
func NewKeyFromString(value string) (key Key, err error) {
	if namespaceKeyPrefix != "" && !strings.HasPrefix(value, namespaceKeyPrefix) {
		if strings.HasPrefix(value, KeySeparator) {
			value = value[1:]
		}
		value = namespaceKeyPrefix + value
	}

	if !validatePathExp.MatchString(value) {
		err = ErrInvalidKeyFormat
	} else {
		key.value = value
	}
	return
}

// NewKeyFromBytesUnchecked creates a new keys from the given byte array without validation of any kind.
func NewKeyFromBytesUnchecked(value []byte) Key {
	return Key{string(value)}
}

// NewKeyFromStringUnchecked creates a new keys from the given byte array without validation of any kind.
func NewKeyFromStringUnchecked(value string) Key {
	return Key{value}
}

// MustSucceed is a handy wrapper for panicking when creating a new Key.  Example: MustSucceed(NewKeyFromString("xyz")).
func MustSucceed(key Key, err error) Key {
	if err != nil {
		panic(err)
	} else {
		return key
	}
}

// Clear any value in this keys.
func (k *Key) Clear() {
	k.value = ""
}

// Parts returns a slice of the component parts of the keys.
func (k *Key) Parts() []string {
	if k.IsEmpty() {
		return nil
	}
	// Skip first empty string from split since key starts with separator
	return strings.Split(k.String(), KeySeparator)[1:]
}

// Base returns the final KeySeparator-separated element of the keys.
func (k *Key) Base() string {
	if k.IsEmpty() {
		return ""
	}
	str := k.String()
	if idx := strings.LastIndex(str, KeySeparator); idx >= 0 {
		return str[idx+1:]
	}
	return str
}

func (k *Key) Prefix() string {
	return k.String()[:strings.LastIndex(k.String(), KeySeparator)]
}

// Prefixes returns a slice of PrefixSeparator-separated elements of the Base.
func (k *Key) Prefixes() []string {
	p := strings.Split(k.Base(), PrefixSeparator)
	if len(p) > 0 {
		p = p[:len(p)-1]
	}
	return p
}

func (k *Key) String() string {
	return k.value
}

func (k *Key) Bytes() []byte {
	return []byte(k.value)
}

func (k *Key) IsEqual(s string) bool {
	return k.String() == s
}

func (k *Key) IsEmpty() bool {
	return k.value == ""
}

func (k *Key) HasValue() bool {
	return k.value != ""
}
