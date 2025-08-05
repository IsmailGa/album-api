package models

type Album struct {
    ID     string   `json:"id" binding:"required"`
    Title  string   `json:"title" binding:"required"`
    Artist string   `json:"artist" binding:"required"`
    Price  float64  `json:"price" binding:"required,gt=0"`
    Tags   []string `json:"tags" binding:"omitempty,dive,required"`
}