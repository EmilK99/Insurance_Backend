package payments

// https://developer.paypal.com/docs/api/payments/v2/#definition-capture_status_details
type CaptureStatusDetails struct {
	Reason string `json:"reason,omitempty"`
}

// PurchaseUnitAmount struct
type PurchaseUnitAmount struct {
	Currency  string                       `json:"currency_code"`
	Value     string                       `json:"value"`
	Breakdown *PurchaseUnitAmountBreakdown `json:"breakdown,omitempty"`
}

// PurchaseUnitAmountBreakdown struct
type PurchaseUnitAmountBreakdown struct {
	ItemTotal        *Money `json:"item_total,omitempty"`
	Shipping         *Money `json:"shipping,omitempty"`
	Handling         *Money `json:"handling,omitempty"`
	TaxTotal         *Money `json:"tax_total,omitempty"`
	Insurance        *Money `json:"insurance,omitempty"`
	ShippingDiscount *Money `json:"shipping_discount,omitempty"`
	Discount         *Money `json:"discount,omitempty"`
}

// Money struct
//
// https://developer.paypal.com/docs/api/orders/v2/#definition-money
type Money struct {
	Currency string `json:"currency_code"`
	Value    string `json:"value"`
}

type SellerProtection struct {
	Status            string   `json:"status,omitempty"`
	DisputeCategories []string `json:"dispute_categories,omitempty"`
}
