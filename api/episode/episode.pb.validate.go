// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: proto/episode.proto

package episode

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

// Validate checks the field values on EpisodeResp with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *EpisodeResp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on EpisodeResp with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in EpisodeRespMultiError, or
// nil if none found.
func (m *EpisodeResp) ValidateAll() error {
	return m.validate(true)
}

func (m *EpisodeResp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Episode

	// no validation rules for Url

	// no validation rules for Download

	// no validation rules for Ext

	// no validation rules for FileSize

	for idx, item := range m.GetSubtitles() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, EpisodeRespValidationError{
						field:  fmt.Sprintf("Subtitles[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, EpisodeRespValidationError{
						field:  fmt.Sprintf("Subtitles[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return EpisodeRespValidationError{
					field:  fmt.Sprintf("Subtitles[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for LastPlayedPosition

	if all {
		switch v := interface{}(m.GetLastPlayedTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, EpisodeRespValidationError{
					field:  "LastPlayedTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, EpisodeRespValidationError{
					field:  "LastPlayedTime",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetLastPlayedTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return EpisodeRespValidationError{
				field:  "LastPlayedTime",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for SkipIntro

	// no validation rules for SkipEnding

	// no validation rules for Title

	// no validation rules for Poster

	// no validation rules for Logo

	if all {
		switch v := interface{}(m.GetAirDate()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, EpisodeRespValidationError{
					field:  "AirDate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, EpisodeRespValidationError{
					field:  "AirDate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAirDate()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return EpisodeRespValidationError{
				field:  "AirDate",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Overview

	// no validation rules for Favorite

	// no validation rules for SeasonId

	// no validation rules for Season

	// no validation rules for SeasonTitle

	// no validation rules for SeriesTitle

	for idx, item := range m.GetActors() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, EpisodeRespValidationError{
						field:  fmt.Sprintf("Actors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, EpisodeRespValidationError{
						field:  fmt.Sprintf("Actors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return EpisodeRespValidationError{
					field:  fmt.Sprintf("Actors[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Filename

	if len(errors) > 0 {
		return EpisodeRespMultiError(errors)
	}

	return nil
}

// EpisodeRespMultiError is an error wrapping multiple validation errors
// returned by EpisodeResp.ValidateAll() if the designated constraints aren't met.
type EpisodeRespMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EpisodeRespMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EpisodeRespMultiError) AllErrors() []error { return m }

// EpisodeRespValidationError is the validation error returned by
// EpisodeResp.Validate if the designated constraints aren't met.
type EpisodeRespValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EpisodeRespValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EpisodeRespValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EpisodeRespValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EpisodeRespValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EpisodeRespValidationError) ErrorName() string { return "EpisodeRespValidationError" }

// Error satisfies the builtin error interface
func (e EpisodeRespValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEpisodeResp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EpisodeRespValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EpisodeRespValidationError{}

// Validate checks the field values on GetEpisodeRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetEpisodeRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetEpisodeRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetEpisodeRequestMultiError, or nil if none found.
func (m *GetEpisodeRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetEpisodeRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return GetEpisodeRequestMultiError(errors)
	}

	return nil
}

// GetEpisodeRequestMultiError is an error wrapping multiple validation errors
// returned by GetEpisodeRequest.ValidateAll() if the designated constraints
// aren't met.
type GetEpisodeRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetEpisodeRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetEpisodeRequestMultiError) AllErrors() []error { return m }

// GetEpisodeRequestValidationError is the validation error returned by
// GetEpisodeRequest.Validate if the designated constraints aren't met.
type GetEpisodeRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetEpisodeRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetEpisodeRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetEpisodeRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetEpisodeRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetEpisodeRequestValidationError) ErrorName() string {
	return "GetEpisodeRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetEpisodeRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetEpisodeRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetEpisodeRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetEpisodeRequestValidationError{}
