package handler

import (
	"EMBECK/controller"
	"EMBECK/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllPlayers(c *fiber.Ctx) error {
	players, err := controller.GetAllPlayers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal ambil data players!",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil ambil data players!",
		"data":    players,
	})
}


func GetPlayerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	player, err := controller.GetPlayerByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Player tidak ditemukan!",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil ambil player!",
		"data":    player,
	})
}

func CreatePlayer(c *fiber.Ctx) error {
	var player model.Player
	if err := c.BodyParser(&player); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal parse data!",
			"error":   err.Error(),
		})
	}

	err := controller.CreatePlayer(c.Context(), player)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal simpan player!",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Player berhasil ditambahkan!",
		"data":    player,
	})
}

func UpdatePlayer(c *fiber.Ctx) error {
	id := c.Params("id")

	var updatedPlayer model.Player
	if err := c.BodyParser(&updatedPlayer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal parse data!",
			"error":   err.Error(),
		})
	}

	// Pastikan ID tidak berubah ke kosong/null
	updatedPlayer.ID = id

	err := controller.UpdatePlayer(c.Context(), id, updatedPlayer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal update player!",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player berhasil diupdate!",
		"data":    updatedPlayer,
	})
}


func DeletePlayer(c *fiber.Ctx) error {
	id := c.Params("id")

	err := controller.DeletePlayer(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hapus player!",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Player berhasil dihapus!",
	})
}

