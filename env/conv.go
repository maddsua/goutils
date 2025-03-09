package env

import (
	"os"
	"strconv"
	"strings"
)

type Value string

func Env(key string) Value {
	return Value(strings.TrimSpace(os.Getenv(strings.ToUpper(key))))
}

func (this Value) IsTrue() bool {
	return strings.ToLower(string(this)) == "true"
}

func (this Value) IsFalse() bool {
	return strings.ToLower(string(this)) == "false"
}

func (this Value) IsEmpty() bool {
	return this == ""
}

func (this Value) ToLower() string {
	return strings.ToLower(string(this))
}

func (this Value) ToUpper() string {
	return strings.ToUpper(string(this))
}

func (this Value) CommaSeparated() []string {

	if this == "" {
		return nil
	}

	var entries []string
	for _, item := range strings.Split(string(this), ",") {

		item = strings.TrimSpace(item)
		if len(item) == 0 {
			continue
		}

		entries = append(entries, item)
	}

	return entries
}

func (this Value) AsInt() IntValue {
	val, err := strconv.Atoi(string(this))
	return IntValue{
		Val:   val,
		Valid: err == nil,
	}
}

type IntValue struct {
	Val   int
	Valid bool
}

func (this IntValue) IntOr(val int) int {
	if !this.Valid {
		return val
	}
	return this.Val
}

func (this IntValue) ToRange(min int, max int) IntValue {

	if !this.Valid || min > max {
		return this
	}

	clamped := this.Val
	if clamped < min {
		clamped = min
	} else if clamped > max {
		clamped = max
	}

	return IntValue{Val: clamped, Valid: true}
}

func EnvAnyOf(keys []string) *EnvMatch {

	for _, key := range keys {
		if val := Env(key); !val.IsEmpty() {
			return &EnvMatch{Key: key, Val: val}
		}
	}

	return nil
}

type EnvMatch struct {
	Key string
	Val Value
}
