package config

import (
	simple "github.com/bitly/go-simplejson"
	jsoniter "github.com/json-iterator/go"
)

type Values struct {
	sj *simple.Json
}

func (val *Values) Scan(v interface{}) error {
	b, err := val.sj.MarshalJSON()
	if err != nil {
		return err
	}

	if string(b) == "null" {
		return nil
	}

	return jsoniter.Unmarshal(b, v)
}

func (val *Values) Set(key []string, value interface{}) {
	val.sj.SetPath(key, value)
}

func (val *Values) Get(path ...string) Values {
	v := val.sj.GetPath(path...)
	return Values{v}
}

func newValues(data map[string]any) Values {
	sj := simple.New()
	sj.SetPath(nil, data)
	return Values{sj: sj}
}
