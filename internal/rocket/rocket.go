package rocket

import (
	"fmt"
	"github.com/badkaktus/gorocket"
	"log"
)

func NewRocket() {
	client := gorocket.NewClient("https://rocketchat-intensive-msk.21-school.ru")

	// login as the main admin user
	login := gorocket.LoginPayload{
		User:     "airi@student.21-school.ru",
		Password: "Sasha27122001",
	}

	lg, err := client.Login(&login)
	for i := 0; i < 200; i++ {
		user, err := client.UsersInfo(&gorocket.SimpleUserRequest{
			Username: "chatayap",
		})
		if !user.Success {
			log.Println(i, err)
			break
		}
	}
	if err != nil {
		fmt.Printf("Error: %+v", err)
	}
	fmt.Printf("I'm %s", lg.Data.Me.Username)
}
