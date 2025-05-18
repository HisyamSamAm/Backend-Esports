package handler

import (
	"EMBECK/controller"
	"EMBECK/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllTournaments(c *fiber.Ctx) error {
	tournaments, err := controller.GetAllTournaments(c.Context())
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
		"data":    tournaments,
	})
}

func GetTournamentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	tournament, err := controller.GetTournamentByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Tournament gak ketemu bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Success ambil 1 tournament!",
		"data":    tournament,
	})
}

func CreateTournament(c *fiber.Ctx) error {
	var tournament model.Tournament

	if err := c.BodyParser(&tournament); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	err := controller.CreateTournament(c.Context(), tournament)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal nambahin data tournament bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Berhasil nambahin data tournament!",
	})
}

func UpdateTournament(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedData model.Tournament

	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	err := controller.UpdateTournament(c.Context(), id, updatedData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal update data tournament bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil update data tournament!",
	})
}

func DeleteTournament(c *fiber.Ctx) error {
	id := c.Params("id")

	err := controller.DeleteTournament(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal hapus data tournament bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Berhasil hapus data tournament!",
	})
}