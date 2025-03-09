package intl

import "encoding/json"

const reallyDefaultLang = "en"

var DefaultLang = reallyDefaultLang

func ResetDefaultLang() {
	DefaultLang = reallyDefaultLang
}

type String map[string]string

type Locales []string

func (this String) Use(langs Locales) string {

	if len(langs) == 0 {
		return this.Default()
	}

	for _, lang := range langs {
		if val, has := this[lang]; has {
			return val
		}
	}

	return this.Default()
}

func (this String) Default() string {

	if val, has := this[DefaultLang]; has {
		return val
	}

	return "[nil intl string]"
}

func (this String) MustMarshall() []byte {
	val, err := json.Marshal(this)
	if err != nil {
		panic("failed to marshall intl.String: " + err.Error())
	}
	return val
}

func MustUnmarshall(data []byte) String {
	var val String
	if err := json.Unmarshal(data, &val); err != nil {
		panic("failed to unmarshall intl.String: " + err.Error())
	}
	return val
}

type Paragraph []String

func (this Paragraph) Use(langs Locales) []string {

	if len(langs) == 0 {
		return this.Default()
	}

	result := make([]string, len(this))
	for idx, val := range this {
		result[idx] = val.Use(langs)
	}

	return result
}

func (this Paragraph) Default() []string {

	result := make([]string, len(this))
	for idx, val := range this {
		result[idx] = val.Default()
	}

	return result
}
