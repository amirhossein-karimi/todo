package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()

		r.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Hello World",
			})
		})

		if err := r.Run(":8080"); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
