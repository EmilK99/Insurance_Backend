package payments

import (
	"context"
	"flight_app/app/store"
	"fmt"
	"github.com/plutov/paypal/v4"
	"strconv"
)

type Client struct {
	Client   *paypal.Client
	ClientID string
	SecretID string
}

func (c *Client) Initialize(ctx context.Context) error {
	// Initialize client
	var err error
	c.Client, err = paypal.NewClient(c.ClientID, c.SecretID, paypal.APIBaseSandBox)
	if err != nil {
		return err
	}

	// Retrieve access token
	_, err = c.Client.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreatePayout(ctx context.Context, contractID int, userEmail string, ticketprice float32) error {

	// Set payout item with Venmo wallet
	payout := paypal.Payout{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			SenderBatchID: "Payouts_" + fmt.Sprint(contractID),
			EmailSubject:  "You have a payout!",
			EmailMessage:  "You have received a payout! Thanks for using our service!",
		},
		Items: []paypal.PayoutItem{
			{
				RecipientType:   "EMAIL",
				RecipientWallet: paypal.PaypalRecipientWallet,
				Receiver:        userEmail,
				Amount: &paypal.AmountPayout{
					Value:    fmt.Sprint(ticketprice),
					Currency: "USD",
				},
				Note:         "Thanks for your patronage!",
				SenderItemID: "201403140001",
			},
		},
	}

	_, err := c.Client.CreatePayout(ctx, payout)
	if err != nil {
		return err
	}
	//fmt.Println(*res)
	return nil
}

func (c *Client) CreateOrder(ctx context.Context, contract store.Contract, returnURL, cancelURL string) (string, error) {
	//fmt.Println(returnURL, cancelURL)
	value := fmt.Sprintf("%.2f", contract.Fee)
	order, err := c.Client.CreateOrder(ctx,
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{
				ReferenceID: strconv.Itoa(contract.ID),
				Amount: &paypal.PurchaseUnitAmount{
					Value:    value,
					Currency: "USD",
				},
			},
		},
		nil,
		&paypal.ApplicationContext{
			ReturnURL: returnURL,
			CancelURL: cancelURL,
		},
	)
	if err != nil {
		return "", err
	}

	//fmt.Println(*order)
	return order.Links[1].Href, nil
}
