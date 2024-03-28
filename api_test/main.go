package main

import (
	"EXAM3/api-gateway/api_test/handlers"
	"EXAM3/api-gateway/api_test/storage/kv"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	kv.Init(kv.NewInMemoryInst())

	router := gin.New()

	router.POST("/user/register", handlers.RegisterUser)
	router.POST("/user/verify/:code", handlers.Verify)
	router.GET("/user/get", handlers.GetUser)
	router.POST("/user/create", handlers.CreateUser)
	router.DELETE("/user/delete", handlers.DeleteUser)
	router.GET("/users", handlers.ListUsers)

	router.GET("/product/get", handlers.GetProduct)
	router.POST("/product/create", handlers.CreateProduct)
	router.DELETE("/product/delete", handlers.DeleteProduct)
	router.GET("/products", handlers.ListProducts)

	log.Fatal(http.ListenAndServe(":9999", router))
}
