package main

import (
    "net/http"
    "sync"
    "github.com/go-playground/validator/v10"
    "strconv" 
    "errors"
    "github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
    ID     string   `json:"id" binding:"required"`
    Title  string   `json:"title" binding:"required"`
    Artist string   `json:"artist" binding:"required"`
    Price  float64  `json:"price" binding:"required,gt=0"`
    Tags   []string `json:"tags" binding:"omitempty,dive,required"`
}

// albums slice to seed record album data.
var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99, Tags: []string{"jazz", "saxophone"}},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99, Tags: []string{"jazz", "baritone saxophone"}},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99, Tags: []string{"jazz", "vocal"}},
}

// mutex to protect albums slice for concurrent access.
var albumsMutex sync.RWMutex

func main() {
    router := gin.Default()
    router.GET("/albums", getAlbums)
    router.GET("/albums/:id", getAlbumByID)
    router.POST("/albums", postAlbums)

    router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
    albumsMutex.RLock()
    defer albumsMutex.RUnlock()
    c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
    var newAlbum album

    if err := c.ShouldBindJSON(&newAlbum); err != nil {
        c.JSON(http.StatusBadRequest, validationErrorResponse(err))
        return
    }

    // Дополнительная проверка: пустые строки в tags
    for i, tag := range newAlbum.Tags {
        if tag == "" {
            c.JSON(http.StatusBadRequest, gin.H{
                "Ошибка валидации": map[string]string{
                    "Теги[" + strconv.Itoa(i) + "]": "не должны быть пустыми",
                },
            })
            return
        }
    }

    albumsMutex.Lock()
    defer albumsMutex.Unlock()

    // Проверка на уникальность ID
    for _, a := range albums {
        if a.ID == newAlbum.ID {
            c.JSON(http.StatusConflict, gin.H{"Ошибка": "альбом с таким ID уже существует"})
            return
        }
    }

    albums = append(albums, newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
}



// getAlbumByID locates the album whose ID value matches the id parameter sent by the client.
func getAlbumByID(c *gin.Context) {
    id := c.Param("id")

    albumsMutex.RLock()
    defer albumsMutex.RUnlock()

    for _, a := range albums {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"Ошибка": "альбом не найден"})
}


func validationErrorResponse(err error) gin.H {
    var ve validator.ValidationErrors
    if errors.As(err, &ve) {
        out := make(map[string]string)
        for _, fe := range ve {
            field := fe.Field()
            switch fe.Tag() {
            case "required":
                out[field] = "обязателен"
            case "gt":
                out[field] = "должен быть больше чем " + fe.Param()
            default:
                out[field] = "валиден"
            }
        }
        return gin.H{"Ошибка валидации": out}
    }
    return gin.H{"Ошибка": err.Error()}
}