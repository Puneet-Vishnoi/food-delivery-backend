package apigateway
import "github.com/gin-gonic/gin"

// var services map[string]*models.ServiceConfig

// func init() {
// 	services = make(map[string]*models.ServiceConfig)
// }

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/restaurents")
	}
}
