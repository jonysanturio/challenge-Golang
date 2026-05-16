package main

import (
	"quisur-challenge/config"
	"quisur-challenge/handlers"
	"quisur-challenge/middleware"
	"quisur-challenge/models"
	"quisur-challenge/seeders"
	"quisur-challenge/websocket"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil{
		log.Printf("ERROR: file not found")
	}

	// Inicializacion de configuracion
	config.Init()
	
	// Conexion a BD
	db = config.GetDB()
	if db == nil {
		log.Fatal("Error al conectar a la base de datos")
	}

	// Migracion de los models
	if err := dbAutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.ProductCategory{},
		&models.ProductHistory{},
	); err != nil {
		log.Fatal("Error al migrar los models: %v", err)
	}

	// Correr los seeders
	if err := seeders.RunSeeders(db); err != nil {
		log.Fatal("Error al ejecutar los seeders: %v", err)
	}

	// Configuracion del router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// WebSockets
	hub := websocket.NewHub()
	go hub.Run()

	// Llamado del router
	public := router.Group("/api")
	{
		// Llamado de las categorias
		categories := public.Group("/categories")
		{
			categories.GET("", handlers.GetCategories)
			categories.GET("/:id", handlers.GetCategoryByID)
		}
		// Llamado de los productos
		products := public.Group("/products")
		{
			products.GET("", handlers.GetProducts)
			products.GET("/:id", handlers.GetProductByID)
			products.GET("/:id/history", handlers.GetProduct)
		}
		// Busquedas
		search := public.Group("/search")
		{
			search.GET("", handlers.Search)
		}
	}
}

// Autenticacion
protected := router.Group("/api")
protected.Use()
{
	admin := protected.Group("/")
	admin.Use(middleware.RoleMiddleware("admin"))
	{
		// Productos CRUD
		products := admin.Group("/products")
		{
			products.POST("", handlers.CreateProduct)
			products.PUT("/:id", handlers.UpdateProduct)
			products.DELETE("/:id", handlers.DeleteProduct)	
		}

		// Categorias CRUD
		categories := admin.Group("/categories")
		{
			categories.POST("", handlers.CreateCategory)
			categories.PUT("/:id", handlers.UpdateCategory)
			categories.DELETE("/:id", handlers.DeleteCategory)	
		}
	}
	// Router Cliente
	client := protected.Group("/")
	client.Use(middleware.RoleMiddleware("client", "admin"))
	{
		//
	}
}
// webSocket con los endpoints
router.GET("/ws", func(c *gin.Context){
	webSocket.ServeWs(hub, c.Writer, c.Request)
})

// Inicializar el servidor
port := os.Getenv("PORT")
if port == ""{
	port = "8080"
}
log.Printf("Servidor corriendo en el puerto %s", port)
if err := router.Run(":" + port); err != nil{
	log.Fatal("Error al iniciar el servidor: %v", err)
}