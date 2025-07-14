package handler

import (
	"EMBECK/model"
	repo "EMBECK/repository"

	"github.com/gofiber/fiber/v2"
)

func GetAllOrders(c *fiber.Ctx) error {
	orders, err := repo.GetAllOrders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "error nih servernya bre!",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success ngambil data bre!",
		"data":    orders,
	})
}

func GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")

	order, err := repo.GetOrderByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Order gak ketemu bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Success ambil 1 order!",
		"data":    order,
	})
}

func CreateOrder(c *fiber.Ctx) error {
	var order model.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	err := repo.CreateOrder(c.Context(), order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal nambahin data order bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Berhasil nambahin data order bre!",
	})
}
func UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order model.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	err := repo.UpdateOrder(c.Context(), id, order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal update data order bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil update data order bre!",
	})
}
func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repo.DeleteOrder(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal hapus data order bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil hapus data order bre!",
	})
}