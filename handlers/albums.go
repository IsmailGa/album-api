package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "example/web-service-gin/models"
    "example/web-service-gin/services"
    "errors"
)

type AlbumHandler struct {
    service *services.AlbumService
}

func NewAlbumHandler(s *services.AlbumService) *AlbumHandler {
    return &AlbumHandler{service: s}
}

func (h *AlbumHandler) GetAlbums(c *gin.Context) {
    albums := h.service.GetAlbums()
    c.JSON(http.StatusOK, albums)
}

func (h *AlbumHandler) GetAlbumByID(c *gin.Context) {
    id := c.Param("id")
    album, err := h.service.GetAlbumByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "альбом не найден"})
        return
    }
    c.JSON(http.StatusOK, album)
}

func (h *AlbumHandler) PostAlbum(c *gin.Context) {
    var newAlbum models.Album

    if err := c.ShouldBindJSON(&newAlbum); err != nil {
        c.JSON(http.StatusBadRequest, validationErrorResponse(err))
        return
    }

    for i, tag := range newAlbum.Tags {
        if tag == "" {
            c.JSON(http.StatusBadRequest, gin.H{
                "validation_error": map[string]string{
                    "tags[" + strconv.Itoa(i) + "]": "не должны быть пустыми",
                },
            })
            return
        }
    }

    err := h.service.AddAlbum(newAlbum)
    if err != nil {
        if err.Error() == "duplicate" {
            c.JSON(http.StatusConflict, gin.H{"error": "альбом с таким ID уже существует"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при сохранении альбома"})
        }
        return
    }

    c.JSON(http.StatusCreated, newAlbum)
}

func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {
    id := c.Param("id")
    var updatedAlbum models.Album

    if err := c.ShouldBindJSON(&updatedAlbum); err != nil {
        c.JSON(http.StatusBadRequest, validationErrorResponse(err))
        return
    }

    for i, tag := range updatedAlbum.Tags {
        if tag == "" {
            c.JSON(http.StatusBadRequest, gin.H{
                "validation_error": map[string]string{
                    "tags[" + strconv.Itoa(i) + "]": "не должны быть пустыми",
                },
            })
            return
        }
    }

    err := h.service.UpdateAlbum(id, updatedAlbum)
    if err != nil {
        if err.Error() == "not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "альбом не найден"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при обновлении альбома"})
        }
        return
    }

    c.JSON(http.StatusOK, updatedAlbum)
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
                out[field] = "некорректен"
            }
        }
        return gin.H{"validation_error": out}
    }
    return gin.H{"error": err.Error()}
}
