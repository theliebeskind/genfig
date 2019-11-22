// Code generated by genfig (config built by merging 'default.yml' and 'default.yml'); DO NOT EDIT.

package config

func init() {
	Envs.Default = Config{
		Apis: ConfigApis{
			Google: ConfigApisGoogle{
				Uri: "google.com",
			},
		},
		Db: ConfigDb{
			Pass: "",
			Uri:  "mongdb://localhos:27017/db",
			User: "",
		},
		EmptyArray: []interface{}{},
		List:       []map[interface{}]interface{}{map[interface{}]interface{}{"a": 1, "b": 2}, map[interface{}]interface{}{"a": 3, "b": 4}},
		LongDesc: ConfigLongDesc{
			De: "Lange Beschreibung",
			En: "Long description",
		},
		Project: "genfig",
		Randomizer: ConfigRandomizer{
			Threshold: 0.75,
		},
		Secrets: []string{""},
		Server: ConfigServer{
			Host: "localhost",
			Port: 1234,
		},
		Version: "0.1.0",
		Wip:     true,
	}
	Envs.Default._map = map[string]interface{}{"apis": map[interface{}]interface{}{"google": map[interface{}]interface{}{"uri": "google.com"}}, "db": map[interface{}]interface{}{"pass": "", "uri": "mongdb://localhos:27017/db", "user": ""}, "emptyArray": []interface{}{}, "list": []interface{}{map[interface{}]interface{}{"a": 1, "b": 2}, map[interface{}]interface{}{"a": 3, "b": 4}}, "longDesc": map[interface{}]interface{}{"de": "Lange Beschreibung", "en": "Long description"}, "project": "genfig", "randomizer": map[interface{}]interface{}{"threshold": 0.75}, "secrets": []interface{}{""}, "server": map[interface{}]interface{}{"host": "localhost", "port": 1234}, "version": "0.1.0", "wip": true}
}
