package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AmadoJunior/Gipitty/config"
	"github.com/AmadoJunior/Gipitty/controllers"
	"github.com/AmadoJunior/Gipitty/repos"
	"github.com/AmadoJunior/Gipitty/routes"
	"github.com/AmadoJunior/Gipitty/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	server      *gin.Engine
	ctx         context.Context
	redisClient *redis.Client
	mongoCLient *mongo.Client

	userRepository repos.IUserRepo

	userService services.UserService
	authService services.AuthService

	AuthController controllers.AuthController
	UserController controllers.UserController

	AuthRouteController routes.AuthRouteController
	UserRouteController routes.UserRouteController
)

func init() {
	//Load ENV
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed Loading ENV", err)
	}

	//Context
	ctx = context.Background()

	if err != nil {
		panic(err)
	}

	fmt.Println("user repo succesfully initiated...")

	//Connect to MongoDB
	mongoConn := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)

	if err != nil {
		panic(err)
	}

	fmt.Println("mongodb successfully connected...")

	//Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisUri,
		Password: "",
		DB:       0,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	//Test Redis
	err = redisClient.Set(ctx, "test", "ok", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("redis successfully connected...")

	//Init User Repo
	userRepository = repos.NewUserRepo(ctx)
	err = userRepository.InitRepository(mongoClient, "Gipitty", "users")
	if err != nil {
		panic(err)
	}

	//Auth
	userService = services.NewUserServiceImpl(userRepository, ctx)
	authService = services.NewAuthServiceImpl(userRepository, ctx)

	AuthController = controllers.NewAuthController(authService, userService)
	UserController = controllers.NewUserController(userService)

	AuthRouteController = routes.NewAuthRouteController(AuthController)
	UserRouteController = routes.NewRouteUserController(UserController, userService)

	//Gin Server
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed Loading ENV", err)
	}

	defer userRepository.DeinitRepository()

	value, err := redisClient.Get(ctx, "test").Result()
	if err == redis.Nil {
		fmt.Println("Redis Key Doesn't Exist...")
	} else if err != nil {
		panic(err)
	}

	//Cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true
	server.Use(cors.New(corsConfig))

	//Static
	server.Use(static.Serve("/", static.LocalFile("/home/amado/Documents/Gipitty/public", true)))

	//API
	router := server.Group("/api")
	router.GET("/healthChecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)

	log.Fatal(server.Run(":" + config.Port))
}
