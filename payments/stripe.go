package payments

import (
	"bytes"
	"encoding/json"
	"flight_app/app/store"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func HandleCreatePaymentIntent(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	stripe.Key = "sk_test_51JSJMMEK56Gl43G4cACjvRgB6Cq0BUQTJooNmt456I0Qb9I52eQ0ibAzlMbxj01A2LbRs8U5RyZvVR1jFcopRjU500ILfJHteJ"

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID       string `json:"user_id"`
		FlightNumber string `json:"flight_number"`
		FlightDate   int    `json:"flight_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	//TODO: get unique app id from header
	var newContract = store.Contract{UserID: req.UserID, FlightNumber: req.FlightNumber, FlightDate: req.FlightDate}

	//err := newContract.CreateContract(pool)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	log.Printf("Unable to create contract: %v", err)
	//	return
	//}
	newContract.Fee = rand.Float32() * 10
	newContract.ID = 1
	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(int64(newContract.Fee * 100)),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String(fmt.Sprint(newContract.ID)),
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("pi.New: %v", err)
		return
	}

	writeJSON(w, struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: pi.ClientSecret,
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}

func HandleStripeWebhook(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("PaymentIntent was successful!")

		err = store.VerifyPayment(pool, paymentIntent.Description)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error verifying payment: %v\n", err)
		}
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("PaymentMethod was attached to a Customer!")
	// ... handle other event types
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)

}
