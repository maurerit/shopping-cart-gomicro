package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/maurerit/shopping-cart-gomicro/handling"
	"github.com/maurerit/shopping-cart-gomicro/proto"
	"github.com/maurerit/shopping-cart-gomicro/repository"
	"github.com/micro/go-micro"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	Application struct {
		Namespace string
		Name      string
		Database  struct {
			Driver         string
			Protocol       string
			Host           string
			Port           uint
			Database       string
			Username       string
			Password       string
			Options        string
			MaxConnections int `yaml:"maxConnections"`
		}
	}
}

func main() {
	configuration := Configuration{}

	data, err := ioutil.ReadFile("application.yaml")

	//<editor-fold desc="Configuration Launch Sanity check">
	//TODO: Candidate for shared function for all my services.
	if err != nil {
		panic("Could not read application.yaml for configuration data.")
	}

	err = yaml.Unmarshal([]byte(data), &configuration)

	if err != nil {
		panic("Could not unmarshall application.yaml as yaml")
	}
	//</editor-fold>

	//TODO: Candidate for shared function for all my services
	service := micro.NewService(
		micro.Name(configuration.Application.Namespace + "." + configuration.Application.Name),
	)

	//<editor-fold desc="DB Initialization">
	//TODO: Candidate for shared function for all my services.
	connectionString := configuration.Application.Database.Username +
		":" +
		configuration.Application.Database.Password +
		"@" +
		configuration.Application.Database.Protocol +
		"(" +
		configuration.Application.Database.Host +
		":" +
		fmt.Sprintf("%d", configuration.Application.Database.Port) +
		")/" +
		configuration.Application.Database.Database

	repository.DB, err = sql.Open(configuration.Application.Database.Driver, connectionString)

	if err != nil {
		panic("Could not open database connection using: " + connectionString)
	}
	repository.DB.Ping()
	repository.DB.SetMaxOpenConns(configuration.Application.Database.MaxConnections)
	defer repository.DB.Close()
	//</editor-fold>

	service.Init()
	cartservice.RegisterCartServiceHandler(service.Server(),
		&handling.CartService{service.Client()})
	service.Run()
}
