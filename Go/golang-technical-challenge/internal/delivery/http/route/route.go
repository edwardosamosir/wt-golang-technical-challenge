package route

import (
	"golang-technical-challenge/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	InvoiceController *http.InvoiceController
}

func (c *RouteConfig) Setup() {
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Post("/api/invoices/import", c.InvoiceController.Import)
	c.App.Get("/api/invoices", c.InvoiceController.GetInvoices)
	c.App.Post("/api/invoices", c.InvoiceController.Create)
	c.App.Put("/api/invoices/:invoiceNo", c.InvoiceController.Update)
	c.App.Delete("/api/invoices/:invoiceNo", c.InvoiceController.Delete)
}
