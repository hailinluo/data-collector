package api

import "github.com/gin-gonic/gin"

func InitApi() error {
	g := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	grp := g.Group("/trigger")
	grp.GET("/fund-companies", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "OK",
		})
	})
	if err := g.Run(":8824"); err != nil {
		// TODO log exit
		return err
	}
}
