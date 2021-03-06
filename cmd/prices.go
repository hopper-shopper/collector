package cmd

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/cobra"
	"github.com/steschwa/hopper-analytics-collector/coingecko"
	"github.com/steschwa/hopper-analytics-collector/constants"
	"github.com/steschwa/hopper-analytics-collector/models"
	db "github.com/steschwa/hopper-analytics-collector/mongo"
)

func RegisterPricesCmd(root *cobra.Command) {
	root.AddCommand(pricesCommand)
}

var pricesCommand = &cobra.Command{
	Use:   "prices",
	Short: "Load and save current crypto prices",
	Run: func(cmd *cobra.Command, args []string) {
		dbClient := GetMongo()
		defer dbClient.Disconnect()
		coinGeckoClient := coingecko.NewCoinGeckoClient()

		ids := []constants.CoinGeckoId{
			constants.COINGECKO_AVAX,
			constants.COINGECKO_FLY,
		}
		currencies := []constants.CoinGeckoCurrency{
			constants.COINGECKO_USD,
			constants.COINGECKO_EUR,
		}

		prices, err := coinGeckoClient.CurrentPrice(ids, currencies)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatalln(err)
		}

		priceDocuments := []models.PriceDocument{}
		for coin, priceData := range prices {
			for currency, price := range priceData {
				priceDocuments = append(priceDocuments, models.PriceDocument{
					Coin:      coin,
					Currency:  currency,
					Price:     price,
					Timestamp: time.Now(),
				})
			}
		}

		pricesCollection := &db.PricesCollection{
			Client: dbClient,
		}

		err = pricesCollection.InsertMany(priceDocuments)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatalln(err)
		}
	},
}
