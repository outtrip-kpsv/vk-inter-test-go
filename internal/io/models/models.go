package models

import (
	"vk-inter-test-go/internal/db/repo"
)

type MovieIo struct {
	Movie  repo.Movie   `json:"movie"`
	Actors []repo.Actor `json:"actors"`
}

type ActorIo struct {
	Actor  repo.Actor   `json:"actor"`
	Movies []repo.Movie `json:"movies"`
}
