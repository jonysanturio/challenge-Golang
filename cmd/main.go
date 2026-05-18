package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal" // IMPORTANTE: Faltaba para el Graceful Shutdown
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
	// 1. Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: Archivo .env no encontrado. Usando variables del sistema.")
	}

	// 2. Inicializar BD
	config.Init() // Si tienes esta función en config
	db := config.GetDB() // o GetDB()
	if db == nil {
		log.Fatal("Error al conectar a la base de datos")
	}

	// 3. Migraciones
	// Se usa db.AutoMigrate en vez de dbAutoMigrate
	if err := db.AutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.ProductCategory{},
		&models.ProductHistory{},
	); err != nil {
		log.Fatalf("Error al migrar los models: %v", err)
	}

	// 4. Correr Seeders
	if err := seeders.RunSeeders(db); err != nil {
		log.Fatalf("Error al ejecutar los seeders: %v", err)
	}

	// 5. Inicializar WebSockets
	hub := websocket.NewHub()
	go hub.Run()

	// 6. Inyección de Dependencias
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo, hub)
	productHandler := handlers.NewProductHandler(productService)

	// 7. Configurar Router
	router := gin.Default() // gin.Default ya incluye Logger() y Recovery()

	// --- RUTAS DE WEBSOCKETS ---
	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request) // Asegúrate que el método se llame ServeWs (o ServeWS)
	})

	// --- RUTAS PÚBLICAS ---
	public := router.Group("/api")
	{
		categories := public.Group("/categories")
		{
			// Nota del Tech Lead: Si refactorizas categorias, aquí usarías categoryHandler.GetCategories
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

	// --- RUTAS PROTEGIDAS (Requieren Token y Rol) ---
	protected := router.Group("/api")
	// Aquí deberías agregar tu middleware que verifica que el JWT sea válido:
	// protected.Use(middleware.AuthMiddleware()) 
	{
		// Rutas exclusivas para ADMIN
		admin := protected.Group("/")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			// Productos CRUD (Usamos el handler con Inyección de Dependencias)
			products := admin.Group("/products")
			{
				products.POST("", productHandler.CreateProduct) 
				products.PUT("/:id", productHandler.Update) // Este es el que refactorizamos
				products.DELETE("/:id", productHandler.DeleteProduct)
			}

			// Categorias CRUD
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
			// Aquí irían acciones que un cliente autenticado puede hacer
		}
	}

	// 8. Iniciar el Servidor y Graceful Shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Ejecutamos el servidor en una Goroutine para que no bloquee el hilo principal
	go func() {
		log.Printf("Servidor corriendo en el puerto %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar señal de apagado del sistema operativo (Ctrl+C, Docker stop, etc.)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit 
	
	log.Println("Señal de apagado recibida. Cerrando conexiones...")

	// Damos 5 segundos para que terminen las peticiones HTTP activas
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("El servidor forzó el apagado:", err)
	}

	log.Println("Servidor apagado correctamente. ¡Adiós!")
}