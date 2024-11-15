package environment

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	RedocFolderPath *string = flag.String("PATH_REDOC_FOLDER", "/docs/swagger.json", "Swagger docs folder")

	localDev = flag.String("localDev", "false", "local development")

	singleton *Environment
)

const (
	QRCodeGatewayRootURL          = "QR_CODE_GATEWAY_ROOT_URL"
	QRCodeGatewayToken            = "QR_CODE_GATEWAY_TOKEN"
	WebhookMercadoLivrePaymentURL = "WEBHOOK_MERCADO_LIVRE_PAYMENT"
	DBHost                        = "DB_HOST"
	DBUser                        = "POSTGRES_USER"
	DBPassword                    = "POSTGRES_PASSWORD"
	DBPort                        = "DB_PORT"
	DBName                        = "POSTGRES_DB"
	Region                        = "AWS_REGION"
	CustomerRootAPI               = "CUSTOMER_ROOT_API"
)

type Environment struct {
	qrCodeGatewayRootURL          string
	qrCodeGatewayToken            string
	webhookMercadoLivrePaymentURL string
	dbHost                        string
	dbPort                        string
	dbName                        string
	dbUser                        string
	dbPassword                    string
	region                        string
	customerRootAPI               string
}

func LoadEnvironmentVariables() {
	flag.Parse()

	if localFlag := *localDev; localFlag != "false" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file", err.Error())
		}
	}

	qrCodeGatewayRootURL := getEnvironmentVariable(QRCodeGatewayRootURL)
	qrCodeGatewayToken := getEnvironmentVariable(QRCodeGatewayToken)
	webhookMercadoLivrePaymentURL := getEnvironmentVariable(WebhookMercadoLivrePaymentURL)
	dbHost := getEnvironmentVariable(DBHost)
	dbPort := getEnvironmentVariable(DBPort)
	dbUser := getEnvironmentVariable(DBUser)
	dbPassword := getEnvironmentVariable(DBPassword)
	dbName := getEnvironmentVariable(DBName)
	region := getEnvironmentVariable(Region)
	customerRootAPI := getEnvironmentVariable(CustomerRootAPI)

	once := &sync.Once{}

	once.Do(func() {
		singleton = &Environment{
			qrCodeGatewayRootURL:          qrCodeGatewayRootURL,
			qrCodeGatewayToken:            qrCodeGatewayToken,
			dbHost:                        dbHost,
			dbPort:                        dbPort,
			dbUser:                        dbUser,
			dbPassword:                    dbPassword,
			dbName:                        dbName,
			webhookMercadoLivrePaymentURL: webhookMercadoLivrePaymentURL,
			region:                        region,
			customerRootAPI:               customerRootAPI,
		}
	})
}

func getEnvironmentVariable(key string) string {
	value, hashKey := os.LookupEnv(key)

	if !hashKey {
		log.Fatalf("There is no %v environment variable", key)
	}

	return value
}

func GetWebhookMercadoLivrePaymentURL() string {
	return singleton.webhookMercadoLivrePaymentURL
}

func GetQRCodeGatewayRootURL() string {
	return singleton.qrCodeGatewayRootURL
}

func GetQRCodeGatewayToken() string {
	return singleton.qrCodeGatewayToken
}

func GetDBHost() string {
	return singleton.dbHost
}

func GetDBPort() string {
	return singleton.dbPort
}

func GetDBName() string {
	return singleton.dbName
}

func GetDBUser() string {
	return singleton.dbUser
}

func GetDBPassword() string {
	return singleton.dbPassword
}

func GetRegion() string {
	return singleton.region
}

func GetCustomerRootAPI() string {
	return singleton.customerRootAPI
}
