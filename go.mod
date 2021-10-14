// +heroku goVersion go1.17

module flight_app

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/pgx/v4 v4.13.0
	github.com/plutov/paypal v2.0.5+incompatible // indirect
	github.com/plutov/paypal/v4 v4.4.0
	github.com/portto/solana-go-sdk v1.8.1 // indirect
	github.com/rs/cors v1.8.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stripe/stripe-go v70.15.0+incompatible
)
