package validation

import "regexp"

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var PhoneRegex = regexp.MustCompile(`^\+?[1-9]\d{9,15}$`)
