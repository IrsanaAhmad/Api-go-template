package utils

import (
	"database/sql"
	"encoding/json"
	"time"
)

// NullString Nullable String that overrides sql. NullString
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		// Return null as JSON, omitempty will handle this
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// IsZero returns true if the value is not valid (for omitempty support)
func (ns NullString) IsZero() bool {
	return !ns.Valid
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

// NullInt64 Nullable Int64 that overrides util. NullInt64
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		ni.Valid = true
		ni.Int64 = *i
	} else {
		ni.Valid = false
	}
	return nil
}

// NullFloat64 Nullable Float64 that overrides util. NullFloat64
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullFloat64
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON for NullFloat64
func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

// NullTime Nullable Time that overrides util. NullTime
type NullTime struct {
	sql.NullTime
}

// MarshalJSON for NullTime
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON for NullTime
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		nt.Valid = true
		nt.Time = *t
	} else {
		nt.Valid = false
	}
	return nil
}

// NullDateTime Nullable DateTime that overrides util. NullDateTime
type NullDateTime struct {
	sql.NullTime
}

// MarshalJSON for NullDateTime
func (ndt NullDateTime) MarshalJSON() ([]byte, error) {
	if ndt.Valid {
		return json.Marshal(ndt.Time.Format("2006-01-02 15:04:05"))
	}
	return json.Marshal(nil)
}

// UnmarshalJSON for NullDateTime
func (ndt *NullDateTime) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		ndt.Valid = true
		ndt.Time = *t
	} else {
		ndt.Valid = false
	}
	return nil
}

// NullDate Nullable Date that overrides util. NullDate
type NullDate struct {
	sql.NullTime
}

// MarshalJSON for NullDate
func (nd NullDate) MarshalJSON() ([]byte, error) {
	if nd.Valid {
		return json.Marshal(nd.Time.Format("2006-01-02"))
	}
	return json.Marshal(nil)
}

// UnmarshalJSON for NullDate
func (nd *NullDate) UnmarshalJSON(data []byte) error {
	var t *time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		nd.Valid = true
		nd.Time = *t
	} else {
		nd.Valid = false
	}
	return nil
}

// NullBoolean Nullable Boolean that overrides sql.NullBool
type NullBoolean struct {
	sql.NullBool
}

// MarshalJSON for NullBoolean
func (nb NullBoolean) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON for NullBoolean
func (nb *NullBoolean) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}

// NullBool Nullable Bool that overrides util. NullBool
type NullBool struct {
	sql.NullBool
}

// MarshalJSON for NullBool
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON for NullBool
func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var b *bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		nb.Valid = true
		nb.Bool = *b
	} else {
		nb.Valid = false
	}
	return nil
}
