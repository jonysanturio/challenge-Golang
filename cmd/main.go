package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/jonysanturio/challenge-golang/websocket"
	"github.com/jonysanturio/challenge-golang/config"
	"github.com/jonysanturio/challenge-golang/seeders"

	"github.com/jonysanturio/challenge-golang/internal/handlers"
	"github.com/jonysanturio/challenge-golang/internal/middleware"
	"github.com/jonysanturio/challenge-golang/internal/models"
	"github.com/jonysanturio/challenge-golang/internal/repositories"
	"github.com/jonysanturio/challenge-golang/internal/services"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: Archivo .env no encontrado. Usando variables del sistema.")
	}

	// Inicializar BD
	config.Init() 
	db := config.GetDB() 
	if db == nil {
		log.Fatal("Error al conectar a la base de datos")
	}

	// Migraciones
	if err := db.AutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.ProductCategory{},
		&models.ProductHistory{},
	); err != nil {
		log.Fatalf("Error al migrar los models: %v", err)
	}

	// Correr Seeders
	if err := seeders.RunSeeders(db); err != nil {
		log.Fatalf("Error al ejecutar los seeders: %v", err)
	}

	// Inicializar WebSockets
	hub := websocket.NewHub()
	go hub.Run()

	// Inyección de Dependencias
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo, hub)
	productHandler := handlers.NewProductHandler(productService)

	// Configurar Router
	router := gin.Default() 

	// Rutas de WebSockets
	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	// Rutas Públicas
	public := router.Group("/api")
	{
		categories := public.Group("/categories")
		{
			categories.GET("", handlers.GetCategories)
			categories.GET("/:id", handlers.GetCategoryByID)
		}
		
		products := public.Group("/products")
		{
			products.GET("", handlers.GetProducts)
			products.GET("/:id", handlers.GetProductByID)
			products.GET("/:id/history", handlers.GetProduct) // ¿Quizás debería llamarse GetProductHistory?
		}
		
		search := public.Group("/search")
		{
			search.GET("", handlers.Search)
		}
	}

	// Rutas Protegidas (Requieren Token y Rol)
	protected := router.Group("/api")
	{
		admin := protected.Group("/")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			products := admin.Group("/products")
			{
				products.POST("", productHandler.CreateProduct) 
				products.PUT("/:id", productHandler.Update)
				products.DELETE("/:id", productHandler.DeleteProduct)
			}

			categories := admin.Group("/categories")
			{
				categories.POST("", handlers.CreateCategory)
				categories.PUT("/:id", handlers.UpdateCategory)
				categories.DELETE("/:id", handlers.DeleteCategory)
			}
		}

		// Rutas para CLIENTE y ADMIN
		client := protected.Group("/")
		client.Use(middleware.RoleMiddleware("client", "admin"))
		{
			// Acciones de clientes
		}
	}

	// Iniciar el Servidor y Graceful Shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Servidor corriendo en el puerto %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit 
	
	log.Println("Señal de apagado recibida. Cerrando conexiones...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("El servidor forzó el apagado:", err)
	}

	log.Println("Servidor apagado correctamente. ¡Adiós!")
}