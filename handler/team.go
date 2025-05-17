package handler

import (
	"EMBECK/controller"
	"EMBECK/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllTeams(c *fiber.Ctx) error {
	team, err := controller.GetAllTeams(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.StatusInternalServerError,
			"message": "error nih servernya bre!",
			"data":    nil,
		})
}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"message": "success ngambil data bre!",
		"data":    team,
	})
}	

func GetTeamByID(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := controller.GetTeamByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Team gak ketemu bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Success ambil 1 team!",
		"data":    team,
	})
}

func CreateTeam(c *fiber.Ctx) error {
	var team model.Team

	if err := c.BodyParser(&team); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data gak valid bre!",
			"data":    nil,
		})
	}

	id, err := controller.CreateTeam(c.Context(), team)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal nambahin data team bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Team berhasil ditambahkan bre!",
		"data":    id,
	})
}

func UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	var team model.Team

	if err := c.BodyParser(&team); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Format data salah bre!",
			"data":    nil,
		})
	}

	if err := controller.UpdateTeam(c.Context(), id, team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal update team bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Team berhasil di-update bre!",
	})
}

func DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := controller.DeleteTeam(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Gagal hapus team bre!",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Team berhasil dihapus bre!",
	})
}
