// package main

// import (
//     "github.com/gin-gonic/gin"
//     "net/http"
// )

// func main() {
//     // Create a new Gin router
//     router := gin.Default()

//     // Define a route handler
//     router.GET("/", func(c *gin.Context) {
//         c.String(http.StatusOK, "Hello World!")
//     })

//	    // Run the server on port 8080
//	    router.Run(":8080")
//	}
package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// "log"
	// "time"
)

//// Using Built-in Middleware
// func LoggerMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()
// 		c.Next()
// 		duration := time.Since(start)
// 		log.Printf("Request - Method: %s | Status: %d | Duration: %v", c.Request.Method, c.Writer.Status(), duration)
// 	}
// }

// func main() {
// 	router := gin.Default()

// 	// Use our custom logger middleware
// 	router.Use(LoggerMiddleware())

// 	router.GET("/", func(c *gin.Context) {
// 		c.String(200, "Hello, World!")
// 	})

// 	router.Run(":8080")
// }

//// Creating Custom Middleware
// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		apiKey := c.GetHeader("X-API-Key")
// 		if apiKey == "" {
// 			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
// 			return
// 		}
// 		c.Next()
// 	}
// }

// func main() {
// 	router := gin.Default()

// 	// Use our custom authentication middleware for a specific group of routes
// 	authGroup := router.Group("/api")
// 	authGroup.Use(AuthMiddleware())
// 	{
// 		authGroup.GET("/data", func(c *gin.Context) {
// 			c.JSON(200, gin.H{"message": "Authenticated and authorized!"})
// 		})
// 	}

// 	router.Run(":8080")
// }

//// Basic Routing
// func main() {
// 	router := gin.Default()

// 	// Basic route
// 	router.GET("/", func(c *gin.Context) {
// 		c.String(200, "Hello, World!")
// 	})

// 	// Route with URL parameters
// 	router.GET("/users/:id", func(c *gin.Context) {
// 		id := c.Param("id")
// 		c.String(200, "User ID: "+id)
// 	})

// 	// Route with query parameters
// 	router.GET("/search", func(c *gin.Context) {
// 		query := c.DefaultQuery("q", "default-value")
// 		c.String(200, "Search query: "+query)
// 	})

// 	router.Run(":8080")
// }

//// Route Groups
// func main() {
// 	router := gin.Default()

// 	// Public routes (no authentication required)
// 	public := router.Group("/public")
// 	{
// 		public.GET("/info", func(c *gin.Context) {
// 			c.String(200, "Public information")
// 		})
// 		public.GET("/products", func(c *gin.Context) {
// 			c.String(200, "Public product list")
// 		})
// 	}

// // Private routes (require authentication)
// private := router.Group("/private")
// private.Use(AuthMiddleware())
// {
// 	private.GET("/data", func(c *gin.Context) {
// 		c.String(200, "Private data accessible after authentication")
// 	})
// 	private.POST("/create", func(c *gin.Context) {
// 		c.String(200, "Create a new resource")
// 	})
// }

// 	router.Run(":8080")
// }

//// Controllers and Handlers
// type UserController struct{}

// // GetUserInfo is a controller method to get user information
// func (uc *UserController) GetUserInfo(c *gin.Context) {
// 	userID := c.Param("id")
// 	// Fetch user information from the database or other data source
// 	// For simplicity, we'll just return a JSON response.
// 	c.JSON(200, gin.H{"id": userID, "name": "John Doe", "email": "john@example.com"})
// }

// func main() {
// 	router := gin.Default()

// 	userController := &UserController{}

// 	// Route using the UserController
// 	router.GET("/users/:id", userController.GetUserInfo)

// 	router.Run(":8080")
// }

// // To-do app
type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	router := gin.Default()

	// Connect to the SQLite database
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the Todo model to create the table
	db.AutoMigrate(&Todo{})

	// Route to create a new Todo
	router.POST("/todos", func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		// Save the Todo to the database
		db.Create(&todo)

		c.JSON(200, todo)
	})

	// Route to get all Todos
	router.GET("/todos", func(c *gin.Context) {
		var todos []Todo

		// Retrieve all Todos from the database
		db.Find(&todos)

		c.JSON(200, todos)
	})

	// Route to get a specific Todo by ID
	router.GET("/todos/:id", func(c *gin.Context) {
		var todo Todo
		todoID := c.Param("id")

		result := db.First(&todo, todoID)
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}

		c.JSON(200, todo)
	})

	// Route to update a Todo by ID
	router.PUT("/todos/:id", func(c *gin.Context) {
		var todo Todo
		todoID := c.Param("id")

		result := db.First(&todo, todoID)
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}

		var updatedTodo Todo
		if err := c.ShouldBindJSON(&updatedTodo); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON data"})
			return
		}

		todo.Title = updatedTodo.Title
		todo.Description = updatedTodo.Description
		db.Save(&todo)

		c.JSON(200, todo)
	})

	// Route to delete a Todo by ID
	router.DELETE("/todos/:id", func(c *gin.Context) {
		var todo Todo
		todoID := c.Param("id")

		result := db.First(&todo, todoID)
		if result.Error != nil {
			c.JSON(404, gin.H{"error": "Todo not found"})
			return
		}

		db.Delete(&todo)
		c.JSON(200, gin.H{"message": fmt.Sprintf("Todo with ID %s deleted", todoID)})
	})

	router.Run(":8080")
}
