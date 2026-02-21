package swagger

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

func init() {
	SwaggerInfo.Version = "1.0"
	SwaggerInfo.Title = "GymApp API"
	SwaggerInfo.Description = "API for managing gym training plans and recipes"
	SwaggerInfo.BasePath = "/"
	SwaggerInfo.Schemes = []string{"http", "https"}
}
