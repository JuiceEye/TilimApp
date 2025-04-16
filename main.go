package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hashed := "$2a$10$7205ZNa2CYI6NZu9S3BM.OpcURURmcUHCd7DqIyNIBZTIdIZXpiGy"
	input := "11111111111" // сюда введённый пароль

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
	if err != nil {
		fmt.Println("❌ Пароль НЕ совпадает")
	} else {
		fmt.Println("✅ Пароль совпадает!")
	}
}
