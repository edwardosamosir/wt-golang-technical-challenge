package config

import (
	"golang-technical-challenge/internal/delivery/http"
	"golang-technical-challenge/internal/delivery/http/route"
	"golang-technical-challenge/internal/repository"
	"golang-technical-challenge/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// add repository setup here
	invoiceRepository := repository.NewInvoiceRepository(config.Log)

	// add usecase setup here
	invoiceUseCase := usecase.NewInvoiceUseCase(config.DB, config.Log, config.Validate, invoiceRepository)

	// add controller here
	invoiceController := http.NewInvoiceController(invoiceUseCase, config.Log)

	routeConfig := route.RouteConfig{
		App:               config.App,
		InvoiceController: invoiceController,
	}
	routeConfig.Setup()
}
