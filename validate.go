package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func containsStr(list []string, str string) bool {
	for _, lstr := range list {
		if lstr == str {
			return true
		}
	}
	return false
}

// IsValid checks a single string to see if it meets certain requirements based on the other
// parameters as follows:
// - allow (ignored if empty string) is a comma-separated list of strings that are the only
//   valid values for the givne string.
// - goodchars (ignored if empty string) is a string of characters that are the only valid
//   characters that may be in the string.
// - badchars (ignored if empty string) is a string of characters that are explicitly disallowed
//   from being in the string.
// - match (ignored if empty string) is a regular expression that the given string must match
// Returned values will be true if valid and false with a descriptive string if invalid.
func IsValid(str, allow, goodchars, badchars, match string) (bool, string) {
	if len(allow) > 0 && !containsStr(strings.Split(allow, ","), str) {
		return false, "Value " + str + "not in " + allow
	}

	if len(goodchars) > 0 {
		for c := 0; c < len(str); c++ {
			var good bool
			for g := 0; g < len(goodchars); g++ {
				if str[c] == goodchars[g] {
					good = true
					break
				}
			}
			if !good {
				return false, "Value " + str + " contains invalid char " + str[c:c+1]
			}
		}
	}

	if len(badchars) > 0 {
		for c := 0; c < len(str); c++ {
			for g := 0; g < len(badchars); g++ {
				if str[c] == badchars[g] {
					return false, "Value " + str + " contains invalid char " + str[c:c+1]
				}
			}
		}
	}

	if len(match) > 0 {
		if matched, err := regexp.MatchString(match, str); !matched || err != nil {
			return false, "Value " + str + " does not match " + match
		}
	}

	return true, ""
}

// Validate checks all string fields and string slice fields in the given struct to see if
// all meet certain requirements given in the struct field tags.
// - allowVals (ignored if empty string) is a comma-separated list of strings that are the only
//   valid values for the givne string.
// - goodChars (ignored if empty string) is a string of characters that are the only valid
//   characters that may be in the string.
// - badChars (ignored if empty string) is a string of characters that are explicitly disallowed
//   from being in the string.
// - match (ignored if empty string) is a regular expression that the given string must match
// - require (ignored if empty string) if this contains any string at all, the length of the field
//   must be greater than zero
// This function returns an error containing a description of which field is invalid and why or nil
// if the struct is valid.
func Validate(inobj interface{}) error {
	v := reflect.ValueOf(inobj)
	t := reflect.TypeOf(inobj)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		tag := t.Field(i).Tag

		req := tag.Get("require")
		allow := tag.Get("allowVals")
		goodchars := tag.Get("goodChars")
		badchars := tag.Get("badChars")
		match := tag.Get("match")

		doIt := len(req) > 0 || len(allow) > 0 || len(goodchars) > 0 || len(badchars) > 0 || len(match) > 0

		if field.CanInterface() && doIt {
			switch field.Interface().(type) {
			case string:
				if field.Len() > 0 {
					valid, str := IsValid(field.String(), allow, goodchars, badchars, match)
					if !valid {
						return fmt.Errorf("Field %v: %v", t.Field(i).Name, str)
					}
				} else if len(req) > 0 {
					return fmt.Errorf("Field %v: Missing required field", t.Field(i).Name)
				}
			case []string:
				if field.Len() > 0 {
					list, _ := field.Interface().([]string)
					for _, str := range list {
						valid, str := IsValid(str, allow, goodchars, badchars, match)
						if !valid {
							return fmt.Errorf("Field %v: %v", t.Field(i).Name, str)
						}
					}
				} else if len(req) > 0 {
					return fmt.Errorf("Missing required field %v", t.Field(i).Name)
				}
			default:
				continue
			}
		}
	}

	return nil
}

// ValidateUseName checks all string fields and string slice fields in the given struct to see if
// all meet certain requirements given in the struct field tags and returns errors strings
// using field names from the given tagName (for example "json")
// - allowVals (ignored if empty string) is a comma-separated list of strings that are the only
//   valid values for the givne string.
// - goodChars (ignored if empty string) is a string of characters that are the only valid
//   characters that may be in the string.
// - badChars (ignored if empty string) is a string of characters that are explicitly disallowed
//   from being in the string.
// - match (ignored if empty string) is a regular expression that the given string must match
// - require (ignored if empty string) if this contains any string at all, the length of the field
//   must be greater than zero
// This function returns an error containing a description of which field is invalid and why or nil
// if the struct is valid.
func ValidateUseName(inobj interface{}, tagName string) error {
	v := reflect.ValueOf(inobj)
	t := reflect.TypeOf(inobj)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		tag := t.Field(i).Tag

		req := tag.Get("require")
		allow := tag.Get("allowVals")
		goodchars := tag.Get("goodChars")
		badchars := tag.Get("badChars")
		match := tag.Get("match")

		doIt := len(req) > 0 || len(allow) > 0 || len(goodchars) > 0 || len(badchars) > 0 || len(match) > 0

		if field.CanInterface() && doIt {
			switch field.Interface().(type) {
			case string:
				if field.Len() > 0 {
					valid, str := IsValid(field.String(), allow, goodchars, badchars, match)
					if !valid {
						return fmt.Errorf("Field %v: %v", tag.Get(tagName), str)
					}
				} else if len(req) > 0 {
					return fmt.Errorf("Field %v: Missing required field", tag.Get(tagName))
				}
			case []string:
				if field.Len() > 0 {
					list, _ := field.Interface().([]string)
					for _, str := range list {
						valid, str := IsValid(str, allow, goodchars, badchars, match)
						if !valid {
							return fmt.Errorf("Field %v: %v", tag.Get(tagName), str)
						}
					}
				} else if len(req) > 0 {
					return fmt.Errorf("Missing required field %v", tag.Get(tagName))
				}
			default:
				continue
			}
		}
	}

	return nil
}
