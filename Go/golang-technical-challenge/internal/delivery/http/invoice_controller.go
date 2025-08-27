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

func (c *InvoiceController) GetInvoices(ctx *fiber.Ctx) error {
	date := ctx.Query("date")
	page := ctx.QueryInt("page", 1)
	size := ctx.QueryInt("size", 10)

	response, err := c.UseCase.GetInvoices(ctx.UserContext(), date, page, size)
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.InvoiceListResponse]{
		Data: response,
	})
}

func (c *InvoiceController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateInvoiceRequest)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
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
		c.Log.WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	response, err := c.UseCase.Update(ctx.UserContext(), invoiceNo, request)
	if err != nil {
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
		c.Log.WithError(err).Error("error deleting invoice")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{
		Data: true,
	})
}
