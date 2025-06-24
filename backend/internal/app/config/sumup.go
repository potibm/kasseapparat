package config

type SumupConfig struct {
	ApiKey            string
	MerchantCode      string
	CurrencyCode      string
	CurrencyMinorUnit int
	AffiliateKey      string
	ApplicationId     string
	PublicUrl         string
}

func loadSumupConfig() SumupConfig {
	return SumupConfig{
		ApiKey:            getEnv("SUMUP_API_KEY", ""),
		MerchantCode:      getEnv("SUMUP_MERCHANT_CODE", ""),
		CurrencyCode:      getCurrencyCode(),
		CurrencyMinorUnit: getCurrencyMinorUnit(),
		AffiliateKey:      getEnv("SUMUP_AFFILIATE_KEY", ""),
		ApplicationId:     getEnv("SUMUP_APPLICATION_ID", ""),
		PublicUrl:         getEnv("SUMUP_PUBLIC_URL", ""),
	}
}
