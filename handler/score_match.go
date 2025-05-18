package handler

import (
	"EMBECK/controller"
	"EMBECK/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllScoreMatches(c *fiber.Ctx) error {
	scoreMatches, err := controller.GetAllScoreMatches(c.Context())
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
		"data":    scoreMatches,
	})
}

func GetScoreMatchByID(c *fiber.Ctx) error {
	id := c.Params("id")

	scoreMatch, err := controller.GetAllScoreMatchesByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "ScoreMatch gak ketemu bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Success ambil 1 score match!",
		"data":    scoreMatch,
	})
}

func CreateScoreMatch(c *fiber.Ctx) error {
	var scoreMatch model.Score

	if err := c.BodyParser(&scoreMatch); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	err := controller.CreateScoreMatch(c.Context(), scoreMatch)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal nambahin data score match bre!",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Berhasil nambahin data score match bre!",
	})
}

func UpdateScoreMatch(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedData model.Score

	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	err := controller.UpdateScoreMatch(c.Context(), id, updatedData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal update data score match bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil update data score match bre!",
	})
}

func DeleteScoreMatch(c *fiber.Ctx) error {
	id := c.Params("id")

	err := controller.DeleteScoreMatch(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal hapus data score match bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil hapus data score match bre!",
	})
}