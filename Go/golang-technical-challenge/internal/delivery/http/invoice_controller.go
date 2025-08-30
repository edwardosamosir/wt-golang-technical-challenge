package http

import (
	"golang-technical-challenge/internal/model"
	"golang-technical-challenge/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type InvoiceController struct {
	UseCase *usecase.InvoiceUseCase
	Log     *logrus.Logger
}

func NewInvoiceController(useCase *usecase.InvoiceUseCase, log *logrus.Logger) *InvoiceController {
	return &InvoiceController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *InvoiceController) Import(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		c.Log.WithError(err).Error("Failed to retrieve file from form-data")
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}

	results, err := c.UseCase.ImportInvoices(ctx.UserContext(), fileHeader)
	if err != nil {
		c.Log.WithError(err).Error("Failed to process invoice import")
		return fiber.ErrInternalServerError
	}

	// Handle import errors
	if errs, ok := results.([]model.ImportError); ok {
		return ctx.JSON(model.WebResponse[[]model.ImportError]{Data: errs})
	}

	// Success response
	return ctx.JSON(model.WebResponse[any]{Data: results})
}

func (c *InvoiceController) GetInvoices(ctx *fiber.Ctx) error {
	date := ctx.Query("date")
	page := ctx.QueryInt("page", 1)
	size := ctx.QueryInt("size", 10)

	response, err := c.UseCase.GetInvoices(ctx.UserContext(), date, page, size)
	if err != nil {
		c.Log.WithError(err).Error("Failed to get invoices")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.InvoiceListResponse]{
		Data: response,
	})
}

func (c *InvoiceController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateInvoiceRequest)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Warn("Invalid JSON format for create invoice")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to create invoice")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.InvoiceResponse]{
		Data: response,
	})
}

func (c *InvoiceController) Update(ctx *fiber.Ctx) error {
	invoiceNo := ctx.Params("invoiceNo")

	request := new(model.UpdateInvoiceRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Warn("Invalid JSON format for update invoice")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	response, err := c.UseCase.Update(ctx.UserContext(), invoiceNo, request)
	if err != nil {
		c.Log.WithError(err).WithField("invoice_no", invoiceNo).Error("Failed to update invoice")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.InvoiceResponse]{
		Data: response,
	})
}

func (c *InvoiceController) Delete(ctx *fiber.Ctx) error {
	invoiceNo := ctx.Params("invoiceNo")

	request := &model.DeleteInvoiceRequest{
		InvoiceNo: invoiceNo,
	}

	if err := c.UseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).WithField("invoice_no", invoiceNo).Error("Failed to delete invoice")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{
		Data: true,
	})
}
