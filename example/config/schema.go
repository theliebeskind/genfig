// Code generated by genfig (schema built from 'default.yml') on 2019-11-19T13:33:01+01:00; DO NOT EDIT.

package config

type Config struct {
	Secrets    []string
	Wip        bool
	List       []map[interface{}]interface{}
	Project    string
	Server     ConfigServer
	Db         ConfigDb
	Randomizer ConfigRandomizer
	Apis       ConfigApis
	LongDesc   ConfigLongDesc
	EmptyArray []interface{}
	Version    string
	_map       map[string]interface{}
}

type ConfigApis struct {
	Google ConfigApisGoogle
}

type ConfigApisGoogle struct {
	Uri string
}

type ConfigDb struct {
	User string
	Pass string
	Uri  string
}

type ConfigLongDesc struct {
	En string
	De string
}

type ConfigRandomizer struct {
	Threshold float64
}

type ConfigServer struct {
	Port int64
	Host string
}
