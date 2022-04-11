/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"

	"github.com/guiPython/codepix/domain/model"
	"github.com/guiPython/codepix/infrastructure/database"
	repository "github.com/guiPython/codepix/infrastructure/repository/pixKey"
	"github.com/spf13/cobra"
)

// fixturesCmd represents the fixtures command
var fixturesCmd = &cobra.Command{
	Use:   "fixtures",
	Short: "Run fixtures for fake data generation",

	Run: func(cmd *cobra.Command, args []string) {
		database := database.ConnectDB(os.Getenv("env"))
		defer database.Close()
		log.Printf("Env: %s\n", os.Getenv("env"))
		pixKeyRepository := repository.NewPixKeyRepository(database)

		santander, _ := model.NewBank("001", "Santander")
		bradesco, _ := model.NewBank("002", "Bradesco")

		pixKeyRepository.AddBank(santander)
		pixKeyRepository.AddBank(bradesco)

		account, _ := model.NewAccount("Marcos", "1111", santander)
		account.ID = "6e4635ce-88d1-4e58-9597-d13fc446ee47"
		pixKeyRepository.AddAccount(account)

		account, _ = model.NewAccount("Jorge", "2222", bradesco)
		account.ID = "51a720b2-5144-4d7f-921d-57023b1e24c1"
		pixKeyRepository.AddAccount(account)

		account, _ = model.NewAccount("Marcia", "3333", bradesco)
		account.ID = "103cc632-78e7-4476-ab63-d5ad3a562d26"
		pixKeyRepository.AddAccount(account)

		account, _ = model.NewAccount("Paula", "4444", santander)
		account.ID = "463b1b2a-b5fa-4b88-9c31-e5c894a20ae3"
		pixKeyRepository.AddAccount(account)
	},
}

func init() {
	rootCmd.AddCommand(fixturesCmd)
}
