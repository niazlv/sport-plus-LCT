package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/niazlv/sport-plus-LCT/internal/routes"
	swagger "github.com/num30/gin-swagger-ui"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

func main() {
	app := gin.Default()

	app.Use(gin.Logger())
	app.Use(cors.Default())

	// Create a new Fizz instance from the Gin engine.
	f := fizz.NewFromEngine(app)

	// Создаем спецификацию безопасности и добавляем её к генератору OpenAPI.
	f.Generator().SetSecuritySchemes(map[string]*openapi.SecuritySchemeOrRef{
		"BearerAuth": {
			SecurityScheme: &openapi.SecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
	})

	// Add Open API description
	infos := &openapi.Info{
		Title:       "Sport Plus",
		Description: "Sport plus API, for LCT Hackaton 2024",
		Version:     "0.2",
	}

	// Create an endpoint for openapi.json file
	f.GET("/openapi.json", nil, f.OpenAPI(infos, "json"))

	// Now add a UI handler
	swagger.AddOpenApiUIHandler(app, "swagger", "/openapi.json")

	// Add handler
	// Second parameter is additional open API info. It's not required and can be nil
	// Third parameter is handler function. It should be a tonic.Handler in order for it to appear on OpenAPI spec
	// f.GET("/hello/:name", []fizz.OperationOption{fizz.Summary("Get a greeting")}, tonic.Handler(func(c *gin.Context, req *GetRequest) (*GetResponse, error) {
	// 	return &GetResponse{Result: "Hello " + req.Name}, nil
	// }, http.StatusOK))

	routes.Setup(f)
	app.Static("/uploads", "./uploads")
	// run our server
	f.Engine().Run(":8000")

	// app.Run(":8000")
}
