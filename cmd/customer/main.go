package main

import (
	"fmt"
	"net/http"

	"github.com/aidosgal/transline-test/pkg/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println("Running service:", cfg.Service.Name)
	fmt.Println("App port:", cfg.Service.Port)
	fmt.Println("Customer service:", cfg.CustomerService.URL, cfg.CustomerService.Port)
	fmt.Println("Postgres:", cfg.Postgres.Host, cfg.Postgres.DBName)

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.CustomerService.Port), nil)
}
