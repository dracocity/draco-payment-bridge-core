package request

// CreateOrder
type CreateOrder struct {
	CreateOrderPayload
	OrderID         string  `json:"order_id,omitempty"`         // Merchant's custom order ID. We recommend using a unique order ID. Example: CGORDER-12345.
	PriceAmount     float64 `json:"price_amount"`               // The price set by the merchant. Example: 1050.99.
	PriceCurrency   string  `json:"price_currency"`             // ISO 4217 currency code which defines the currency in which you wish to price your merchandise used to define price parameter. Supported currencies.
	ReceiveCurrency string  `json:"receive_currency,omitempty"` // ISO 4217 currency code specifying the currency in which you want to receive settlements. Currency conversions are handled by CoinGate. See Supported Settlement Currencies for available options, or use DO_NOT_CONVERT to keep the payment in its original currency. For setup instructions, see How to Configure Settlement Currency.
	Title           string  `json:"title"`                      // Min 3 - Max 150 characters. Example: product title (Apple iPhone 6), order id (MyShop Order #12345), cart id (Cart #00004335).
	Description     string  `json:"description"`                // More details about this order. Min 3 - Max 500 characters. It can be cart items, product details or other information. Example: 1 x Apple iPhone 6, 1 x Apple MacBook Air.
	CallbackURL     string  `json:"callback_url,omitempty"`     // Send an automated message to Merchant URL when order status is changed. For testing you can use requestcatcher.com tool. URL must be direct without redirection.
	CancelURL       string  `json:"cancel_url,omitempty"`       // Redirect to Merchant URL when buyer cancels the order.
	SuccessURL      string  `json:"success_url,omitempty"`      // Redirect to Merchant URL after successful payment.
}

type CreateOrderPayload struct {
	Token   string              `json:"token,omitempty"`   // Your custom token to validate payment callback (notification).
	Shopper *CreateOrderShopper `json:"shopper,omitempty"` // Optional object to enhance the shopper’s experience and prefill the Travel Rule form on the checkout page. Learn more about the Travel Rule. All fields are optional.
}

type CreateOrderShopper struct {
	Type                string                     `json:"type,omitempty"`                  // must be 'business' or 'personal'
	IPAddress           string                     `json:"ip_address,omitempty"`            // Shopper’s IP address
	Email               string                     `json:"email,omitempty"`                 // Shoppers e-mail address
	FirstName           string                     `json:"first_name,omitempty"`            // Shopper’s first name. For business type, provide the representative’s first name
	LastName            string                     `json:"last_name,omitempty"`             // Shopper’s last name. For business type, provide the representative’s last name
	DateOfBirth         string                     `json:"date_of_birth,omitempty"`         // Shopper’s date of birth. For business type, provide the representative’s date of birth. Must be a valid date in the format YYYY-MM-DD
	ResidenceAddress    string                     `json:"residence_address,omitempty"`     // Shopper's residence address. Can be omitted if the shopper is a business
	ResidencePostalCode string                     `json:"residence_postal_code,omitempty"` // Shopper's residence postal code. Can be omitted if the shopper is a business
	ResidenceCity       string                     `json:"residence_city,omitempty"`        // Shopper’s residence city. Can be omitted if the shopper is a business
	ResidenceCountry    string                     `json:"residence_country,omitempty"`     // Shopper's country of residence. Must be a valid Alpha-2 country code. May be omitted if the shopper is a business.
	CompanyDetails      *CreateOrderCompanyDetails `json:"company_details,omitempty"`       // Company details must be provided when the shopper type is 'business'.
}

type CreateOrderCompanyDetails struct {
	Name                 string `json:"name,omitempty"`                  // Company's legal name.
	Code                 string `json:"code,omitempty"`                  // Company code (or business ID number) assigned by the authority under which the company is registered.
	IncorporationDate    string `json:"incorporation_date,omitempty"`    // Date of company's incorporation. Must be a valid date in the format YYYY-MM-DD.
	IncorporationCountry string `json:"incorporation_country,omitempty"` // Country where the company was incorporated. Must be a valid Alpha-2 country code.
	Address              string `json:"address,omitempty"`               // Company’s registered address.
	PostalCode           string `json:"postal_code,omitempty"`           // Postal code of the company’s registered address.
	City                 string `json:"city,omitempty"`                  // Company’s registered address city.
	Country              string `json:"country,omitempty"`               // Company's country of registration. Must be a valid Alpha-2 country code.
}
