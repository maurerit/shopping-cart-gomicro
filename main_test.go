package main

import (
	"gopkg.in/yaml.v2"
	"testing"
)

var configurationString = `application:
	  namespace: your.namespace.to.register
	  name: yourappname
	  database:
	    driver: mysql
	    protocol: tcp
	    host: hostname
	    port: 3306
	    database: database
	    username: username
	    password: password
	    options: someoptions using Pear DB like options or nothing
	    maxConnections: 5`

// Really only here to ensure that the above string is properly parsed into the expected data structure
func TestConfigurationLoading(t *testing.T) {
	configuration := Configuration{}

	err := yaml.Unmarshal([]byte(configurationString), &configuration)

	if err != nil {
		panic("Could not unmarshall application.yaml as yaml")
	}

	if configuration.Application.Name != "yourappname" {
		t.Error("Wrong application name")
	}
}
