// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: system/v1alpha1/zone.proto

package v1alpha1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _zone_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on Zone with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *Zone) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetIngress()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ZoneValidationError{
				field:  "Ingress",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// ZoneValidationError is the validation error returned by Zone.Validate if the
// designated constraints aren't met.
type ZoneValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ZoneValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ZoneValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ZoneValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ZoneValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ZoneValidationError) ErrorName() string { return "ZoneValidationError" }

// Error satisfies the builtin error interface
func (e ZoneValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sZone.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ZoneValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ZoneValidationError{}

// Validate checks the field values on Zone_Ingress with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *Zone_Ingress) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Address

	return nil
}

// Zone_IngressValidationError is the validation error returned by
// Zone_Ingress.Validate if the designated constraints aren't met.
type Zone_IngressValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e Zone_IngressValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e Zone_IngressValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e Zone_IngressValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e Zone_IngressValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e Zone_IngressValidationError) ErrorName() string { return "Zone_IngressValidationError" }

// Error satisfies the builtin error interface
func (e Zone_IngressValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sZone_Ingress.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = Zone_IngressValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = Zone_IngressValidationError{}
