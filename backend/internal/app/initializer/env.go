package initializer

import "github.com/joho/godotenv"

func InitializeDotenv() {
	_ = godotenv.Load()
}
