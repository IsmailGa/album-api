package main

import (
    "github.com/gin-gonic/gin"
    "example/web-service-gin/storage"
    "example/web-service-gin/services"
    "example/web-service-gin/handlers"
)

func main() {
    r := gin.Default()

    store := storage.NewAlbumStore()
    service := services.NewAlbumService(store)
    handler := handlers.NewAlbumHandler(service)

    r.GET("/albums", handler.GetAlbums)
    r.GET("/albums/:id", handler.GetAlbumByID)
    r.POST("/albums", handler.PostAlbum)
    r.PUT("/albums/:id", handler.UpdateAlbum)

    r.Run(":8080")
}
