// @title Фильмотека API
// @description API фильмотеки
// @version 1.0
// @host localhost:3000
// @BasePath /
package main

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"vk-inter-test-go/internal/bl"
	"vk-inter-test-go/internal/config"
	"vk-inter-test-go/internal/db"
	"vk-inter-test-go/internal/io"
	"vk-inter-test-go/internal/io/http/handlers"
)

func main() {
	configSrv, err := config.InitConfServ()
	if err != nil {
		log.Fatal("cannot initialize config")
	}

	dbRepo := db.NewDBRepo(configSrv)
	defer dbRepo.Close()
	blInst := bl.NewBL(dbRepo, configSrv.Logger.Named("bl"))
	controller := handlers.NewController(blInst, configSrv.Logger.Named("io"))

	mux := io.SetupRoutes(controller)
	srvStr := fmt.Sprintf("%s:%s", configSrv.Options.Host, configSrv.Options.Port)
	configSrv.Logger.Info("srv", zap.String("addr: ", fmt.Sprintf("http://%s", srvStr)))
	log.Fatal(http.ListenAndServe(srvStr, mux))

}
