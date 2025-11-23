package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/silentmol/avito-backend-trainee/internal/controller/http"
)

func getRouter(handle *http.Handle, appName string) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: appName,
	})

	app.Post("/team/add", handle.AddTeam)
	app.Get("/team/get", handle.GetTeam)

	app.Post("/users/setIsActive", handle.SetIsActive)
	app.Get("/users/getReview", handle.GetReview)

	app.Post("/pullRequest/create", handle.CreatePR)
	app.Post("/pullRequest/merge", handle.MergePR)
	app.Post("/pullRequest/reassign", handle.ReassignPR)

	return app
}
