// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: proto/movie.proto

package movie

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
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
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on MovieReply with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *MovieReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MovieReply with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in MovieReplyMultiError, or
// nil if none found.
func (m *MovieReply) ValidateAll() error {
	return m.validate(true)
}

func (m *MovieReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Name

	if len(errors) > 0 {
		return MovieReplyMultiError(errors)
	}

	return nil
}

// MovieReplyMultiError is an error wrapping multiple validation errors
// returned by MovieReply.ValidateAll() if the designated constraints aren't met.
type MovieReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MovieReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MovieReplyMultiError) AllErrors() []error { return m }

// MovieReplyValidationError is the validation error returned by
// MovieReply.Validate if the designated constraints aren't met.
type MovieReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MovieReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MovieReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MovieReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MovieReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MovieReplyValidationError) ErrorName() string { return "MovieReplyValidationError" }

// Error satisfies the builtin error interface
func (e MovieReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMovieReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MovieReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MovieReplyValidationError{}

// Validate checks the field values on CreateMovieRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateMovieRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateMovieRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateMovieRequestMultiError, or nil if none found.
func (m *CreateMovieRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateMovieRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	if len(errors) > 0 {
		return CreateMovieRequestMultiError(errors)
	}

	return nil
}

// CreateMovieRequestMultiError is an error wrapping multiple validation errors
// returned by CreateMovieRequest.ValidateAll() if the designated constraints
// aren't met.
type CreateMovieRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateMovieRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateMovieRequestMultiError) AllErrors() []error { return m }

// CreateMovieRequestValidationError is the validation error returned by
// CreateMovieRequest.Validate if the designated constraints aren't met.
type CreateMovieRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateMovieRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateMovieRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateMovieRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateMovieRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateMovieRequestValidationError) ErrorName() string {
	return "CreateMovieRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateMovieRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateMovieRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateMovieRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateMovieRequestValidationError{}

// Validate checks the field values on GetMovieRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetMovieRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetMovieRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetMovieRequestMultiError, or nil if none found.
func (m *GetMovieRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetMovieRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return GetMovieRequestMultiError(errors)
	}

	return nil
}

// GetMovieRequestMultiError is an error wrapping multiple validation errors
// returned by GetMovieRequest.ValidateAll() if the designated constraints
// aren't met.
type GetMovieRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetMovieRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetMovieRequestMultiError) AllErrors() []error { return m }

// GetMovieRequestValidationError is the validation error returned by
// GetMovieRequest.Validate if the designated constraints aren't met.
type GetMovieRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetMovieRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetMovieRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetMovieRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetMovieRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetMovieRequestValidationError) ErrorName() string { return "GetMovieRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetMovieRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetMovieRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetMovieRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetMovieRequestValidationError{}

// Validate checks the field values on GetMovieReply with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetMovieReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetMovieReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetMovieReplyMultiError, or
// nil if none found.
func (m *GetMovieReply) ValidateAll() error {
	return m.validate(true)
}

func (m *GetMovieReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Name

	if len(errors) > 0 {
		return GetMovieReplyMultiError(errors)
	}

	return nil
}

// GetMovieReplyMultiError is an error wrapping multiple validation errors
// returned by GetMovieReply.ValidateAll() if the designated constraints
// aren't met.
type GetMovieReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetMovieReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetMovieReplyMultiError) AllErrors() []error { return m }

// GetMovieReplyValidationError is the validation error returned by
// GetMovieReply.Validate if the designated constraints aren't met.
type GetMovieReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetMovieReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetMovieReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetMovieReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetMovieReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetMovieReplyValidationError) ErrorName() string { return "GetMovieReplyValidationError" }

// Error satisfies the builtin error interface
func (e GetMovieReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetMovieReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetMovieReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetMovieReplyValidationError{}

// Validate checks the field values on FindMovieRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *FindMovieRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FindMovieRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// FindMovieRequestMultiError, or nil if none found.
func (m *FindMovieRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *FindMovieRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetPage()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FindMovieRequestValidationError{
					field:  "Page",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FindMovieRequestValidationError{
					field:  "Page",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPage()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FindMovieRequestValidationError{
				field:  "Page",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Name != nil {
		// no validation rules for Name
	}

	if len(errors) > 0 {
		return FindMovieRequestMultiError(errors)
	}

	return nil
}

// FindMovieRequestMultiError is an error wrapping multiple validation errors
// returned by FindMovieRequest.ValidateAll() if the designated constraints
// aren't met.
type FindMovieRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FindMovieRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FindMovieRequestMultiError) AllErrors() []error { return m }

// FindMovieRequestValidationError is the validation error returned by
// FindMovieRequest.Validate if the designated constraints aren't met.
type FindMovieRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FindMovieRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FindMovieRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FindMovieRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FindMovieRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FindMovieRequestValidationError) ErrorName() string { return "FindMovieRequestValidationError" }

// Error satisfies the builtin error interface
func (e FindMovieRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFindMovieRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FindMovieRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FindMovieRequestValidationError{}

// Validate checks the field values on FindMovieReply with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *FindMovieReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FindMovieReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in FindMovieReplyMultiError,
// or nil if none found.
func (m *FindMovieReply) ValidateAll() error {
	return m.validate(true)
}

func (m *FindMovieReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetPage()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FindMovieReplyValidationError{
					field:  "Page",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FindMovieReplyValidationError{
					field:  "Page",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPage()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FindMovieReplyValidationError{
				field:  "Page",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetList() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, FindMovieReplyValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, FindMovieReplyValidationError{
						field:  fmt.Sprintf("List[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return FindMovieReplyValidationError{
					field:  fmt.Sprintf("List[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return FindMovieReplyMultiError(errors)
	}

	return nil
}

// FindMovieReplyMultiError is an error wrapping multiple validation errors
// returned by FindMovieReply.ValidateAll() if the designated constraints
// aren't met.
type FindMovieReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FindMovieReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FindMovieReplyMultiError) AllErrors() []error { return m }

// FindMovieReplyValidationError is the validation error returned by
// FindMovieReply.Validate if the designated constraints aren't met.
type FindMovieReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FindMovieReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FindMovieReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FindMovieReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FindMovieReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FindMovieReplyValidationError) ErrorName() string { return "FindMovieReplyValidationError" }

// Error satisfies the builtin error interface
func (e FindMovieReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFindMovieReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FindMovieReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FindMovieReplyValidationError{}

// Validate checks the field values on UpdateMovieRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateMovieRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateMovieRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateMovieRequestMultiError, or nil if none found.
func (m *UpdateMovieRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateMovieRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if m.Name != nil {
		// no validation rules for Name
	}

	if len(errors) > 0 {
		return UpdateMovieRequestMultiError(errors)
	}

	return nil
}

// UpdateMovieRequestMultiError is an error wrapping multiple validation errors
// returned by UpdateMovieRequest.ValidateAll() if the designated constraints
// aren't met.
type UpdateMovieRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateMovieRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateMovieRequestMultiError) AllErrors() []error { return m }

// UpdateMovieRequestValidationError is the validation error returned by
// UpdateMovieRequest.Validate if the designated constraints aren't met.
type UpdateMovieRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateMovieRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateMovieRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateMovieRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateMovieRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateMovieRequestValidationError) ErrorName() string {
	return "UpdateMovieRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateMovieRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateMovieRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateMovieRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateMovieRequestValidationError{}