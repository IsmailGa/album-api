package storage

import (
    "sync"
    "example/web-service-gin/models"
    "errors"
)

type AlbumStore struct {
    mu     sync.RWMutex
    albums []models.Album
}

func NewAlbumStore() *AlbumStore {
    return &AlbumStore{
        albums: []models.Album{
            {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99, Tags: []string{"jazz", "saxophone"}},
            {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99, Tags: []string{"jazz", "baritone saxophone"}},
        },
    }
}

func (s *AlbumStore) GetAll() []models.Album {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.albums
}

func (s *AlbumStore) GetByID(id string) (models.Album, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, a := range s.albums {
        if a.ID == id {
            return a, nil
        }
    }
    return models.Album{}, errors.New("not found")
}

func (s *AlbumStore) Add(album models.Album) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    for _, a := range s.albums {
        if a.ID == album.ID {
            return errors.New("duplicate")
        }
    }
    s.albums = append(s.albums, album)
    return nil
}
