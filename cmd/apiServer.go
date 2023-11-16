/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alfredomagalhaes/fiber-user-api/repositories"
	"github.com/alfredomagalhaes/fiber-user-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var apiPort int

// apiServerCmd represents the apiServer command
var apiServerCmd = &cobra.Command{
	Use:   "apiServer",
	Short: "Initialize API server",
	Long: `Run an API server to expose users endpoints to validate
	the integration with AWS Cognito to authorize the application`,
	Run: func(cmd *cobra.Command, args []string) {
		var repoConfig repositories.MySqlRepoConfig
		//Load env file
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("No .env file found")
		}

		//Create log for the application
		file, err := os.OpenFile(os.Getenv("LOG_FILENAME"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
		repoConfig.Host = os.Getenv("DB_HOST")
		repoConfig.Port = os.Getenv("DB_PORT")
		repoConfig.User = os.Getenv("DB_USER")
		repoConfig.Pwd = os.Getenv("DB_PASSWORD")
		repoConfig.DbName = os.Getenv("DB_NAME")
		repo, err := repositories.NewMySqlRepository(repoConfig)

		if err != nil {
			log.Fatalf("error while initializing database connection...\n%v", err)
		}

		//--
		//Initialize the fiber application
		//--
		app := fiber.New()
		app.Use(logger.New(logger.Config{
			Output: file,
		}))

		initializeRoutes(app, repo)

		go func() {
			runPort := fmt.Sprintf(":%d", apiPort)
			if err := app.Listen(runPort); err != nil {
				log.Panic(err)
			}
		}()

		cancelChan := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
		signal.Notify(cancelChan, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

		<-cancelChan // This blocks the main thread until an interrupt is received
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()

		fmt.Println("Running cleanup tasks...")

		fmt.Println("Application was successful shutdown.")

	},
}

func init() {
	rootCmd.AddCommand(apiServerCmd)

	// Here you will define your flags and configuration settings.
	apiServerCmd.Flags().IntVarP(&apiPort, "apiPort", "p", 3000, "define the port that the application will run")
}

func initializeRoutes(app *fiber.App, repo repositories.UserRepository) {

	// give response when at /
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the endpoint 😉 🇧🇷",
		})
	})

	//Create a group to version the api
	apiV1 := app.Group("/api/v1")

	routes.LoginRoutes(apiV1)
	routes.UserRoutes(apiV1, repo)
}
