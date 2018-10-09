package dgraph

import (
	"fmt"
)

// Helper implementation for testing dgraph persistence.
func main() {
	fmt.Println("Hello World")
	Open("127.0.0.1:9080")
	err := CreateSchema()
	if err != nil {
		fmt.Println("Error while creating schema ", err)
	}

	uid, err := GetUId(Client, "default:Pod2", "isPod")
	if err != nil {
		fmt.Println("Error while fetching uid ", err)
	}
	fmt.Println("Uid is " + uid)
}
