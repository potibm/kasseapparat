package initializer

import (
	"os"
	"sync"

	sumupSevice "github.com/potibm/kasseapparat/internal/app/service/sumup"
	"github.com/sumup/sumup-go"
	"github.com/sumup/sumup-go/client"
)

var (
	instance *sumupSevice.Service
	once     sync.Once
)

func InitializeSumup() *sumupSevice.Service {
	once.Do(func() {
		options := client.New()
		clientOptions := options.WithAPIKey(os.Getenv("SUMUP_API_KEY"))

		client := sumup.NewClient(clientOptions)

		instance = sumupSevice.NewService(client, os.Getenv("SUMUP_MERCHANT_CODE"))
	})

	return instance
}

func GetSumupService() *sumupSevice.Service {
	if instance == nil {
		panic("SumUp service is not initialized. Call InitializeSumup first.")
	}

	return instance
}
