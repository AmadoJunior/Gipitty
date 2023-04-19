package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AmadoJunior/Gipitty/config"
	"github.com/AmadoJunior/Gipitty/controllers"
	"github.com/AmadoJunior/Gipitty/routes"
	"github.com/AmadoJunior/Gipitty/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	redisClient *redis.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

func init() {
	//Load ENV
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed Loading ENV", err)
	}

	//Context
	ctx = context.Background()

	//Connect to MongoDB
	mongoConn := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)

	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB Successfully Connected...")

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

	fmt.Println("Redis Successfully Connected...")

	//Collection
	authCollection = mongoClient.Database("Gipitty").Collection("users")

	//Auth
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthServiceImpl(authCollection, ctx)

	AuthController = controllers.NewAuthController(authService, userService)
	UserController = controllers.NewUserController(userService)

	AuthRouteController = routes.NewAuthRouteController(AuthController)
	UserRouteController = routes.NewRouteUserController(UserController)

	//Gin Server
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Failed Loading ENV", err)
	}

	defer mongoClient.Disconnect(ctx)

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

	router := server.Group("/api")
	router.GET("/healthChecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)

	log.Fatal(server.Run(":" + config.Port))
}
