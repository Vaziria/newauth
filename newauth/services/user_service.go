package services

import "time"

type UserService struct{}

type PayloadVerif struct {
	UserID    uint      `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Email     string    `json:"email"`
}

func (sr *UserService) CreateVerifLink() {

}

func (sr *UserService) Serialize()   {}
func (sr *UserService) Deserialize() {}

func (sr *UserService) AcceptVerifLink() {

}
