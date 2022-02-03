package main

import (
	Types "gin_demo/router"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	Types.RegisterRouter(router)
	router.Run(":80")

}
