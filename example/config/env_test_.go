// Code generated by genfig (config built by merging 'default.yml' and 'test.json') on 2019-07-23T23:30:42+02:00; DO NOT EDIT.

package config

func init() {
	Envs.Test = Config{
		Randomizer: ConfigRandomizer{
			Threshold: 0.12345,
		},
		LongDesc: ConfigLongDesc{
			En: "Long description",
			De: "Lange Beschreibung",
		},
		Project: "genfig",
		Server: ConfigServer{
			Port: 1234,
			Host: "localhost",
		},
		Db: ConfigDb{
			User: "chuck",
			Pass: "norris",
			Uri:  "mongdb://${db.user}:${db.pass}@remotedb:27018/proddb",
		},
		Wip:     true,
		Version: "1-test",
		Secrets: []string{"ChuckNorriscanwinagameofConnectFourinonlythreemoves"},
		Apis: ConfigApis{
			Google: ConfigApisGoogle{
				Uri: "google.com",
			},
		},
	}
	Envs.Test._map = map[string]interface{}{"apis": map[interface{}]interface{}{"google": map[interface{}]interface{}{"uri": "google.com"}}, "db": map[interface{}]interface{}{"pass": "norris", "uri": "mongdb://${db.user}:${db.pass}@remotedb:27018/proddb", "user": "chuck"}, "longDesc": map[interface{}]interface{}{"de": "Lange Beschreibung", "en": "Long description"}, "project": "genfig", "randomizer": map[interface{}]interface{}{"threshold": 0.12345}, "secrets": []interface{}{"ChuckNorriscanwinagameofConnectFourinonlythreemoves"}, "server": map[interface{}]interface{}{"host": "localhost", "port": 1234}, "version": "1-test", "wip": true}
}
