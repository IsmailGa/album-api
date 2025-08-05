package services

import (
    "example/web-service-gin/models"
    "example/web-service-gin/storage"
)

type AlbumService struct {
    store *storage.AlbumStore
}

func NewAlbumService(s *storage.AlbumStore) *AlbumService {
    return &AlbumService{store: s}
}

func (s *AlbumService) GetAlbums() []models.Album {
    return s.store.GetAll()
}

func (s *AlbumService) GetAlbumByID(id string) (models.Album, error) {
    return s.store.GetByID(id)
}

func (s *AlbumService) AddAlbum(album models.Album) error {
    return s.store.Add(album)
}
