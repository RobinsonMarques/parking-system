package input

import "github.com/RobinsonMarques/parking-system/database"

type CreateUserInput struct {
	Person   database.Person
	Document string `json:"Document" binding:"required"`
}

type UpdateUserInput struct {
	Person     database.Person
	Document   string     `json:"Document"`
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type UpdateAdminInput struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type UpdateTrafficWarden struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type UpdateVehicle struct {
	LicensePlate    string     `json:"LicensePlate" binding:"required"`
	NewLicensePlate string     `json:"NewLicensePlate" binding:"required"`
	VehicleModel    string     `json:"VehicleModel"`
	VehicleType     string     `json:"VehicleType"`
	LoginInput      LoginInput `json:"Login" binding:"required"`
}

type CreateAdminInput struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type CreateTrafficWarden struct {
	Person     database.Person
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type UpdateVehicleOwner struct {
	NewUserID  uint       `json:"NewUserID" binding:"required"`
	LoginInput LoginInput `json:"Login" binding:"required"`
}

type CreateParkingTicket struct {
	Login       LoginInput
	Location    string `json:"Location"`
	ParkingTime int    `json:"ParkingTime"`
	VehicleID   uint   `json:"VehicleID"`
}

type CreateVehicle struct {
	LicensePlate string `json:"LicensePlate" binding:"required"`
	VehicleModel string `json:"VehicleModel" binding:"required"`
	VehicleType  string `json:"VehicleType" binding:"required"`
	UserID       uint   `json:"UserID" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

type CreateRecharge struct {
	Value       int64      `json:"Value" binding:"required"`
	PaymentType string     `json:"PaymentType" binding:"required"`
	LoginInput  LoginInput `json:"Login" binding:"required"`
}

type Payment struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Amount        int64 `json:"amount"`
	BilletDetails struct {
		BankAccount   string `json:"bankAccount"`
		BarcodeNumber string `json:"barcodeNumber"`
		OurNumber     string `json:"ourNumber"`
		Portfolio     string `json:"portfolio"`
	} `json:"billetDetails"`
	CheckoutURL     string `json:"checkoutUrl"`
	Code            int64  `json:"code"`
	DueDate         string `json:"dueDate"`
	ID              string `json:"id"`
	InstallmentLink string `json:"installmentLink"`
	Link            string `json:"link"`
	PayNumber       string `json:"payNumber"`
	Payments        []struct {
		Amount        int64       `json:"amount"`
		ChargeID      string      `json:"chargeId"`
		Date          string      `json:"date"`
		FailReason    interface{} `json:"failReason"`
		Fee           float64     `json:"fee"`
		ID            string      `json:"id"`
		ReleaseDate   string      `json:"releaseDate"`
		Status        string      `json:"status"`
		TransactionID interface{} `json:"transactionId"`
		Type          string      `json:"type"`
	} `json:"payments"`
	Reference string `json:"reference"`
	Status    string `json:"status"`
}

type ValidateAccessToken struct {
	Details []struct {
		ErrorCode string `json:"errorCode"`
		Message   string `json:"message"`
	} `json:"details"`
	Error     string `json:"error"`
	Path      string `json:"path"`
	Status    int64  `json:"status"`
	Timestamp string `json:"timestamp"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Jti         string `json:"jti"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	UserName    string `json:"user_name"`
}
type Billing struct {
	Name     string `json:"name"`
	Document string `json:"document"`
	Email    string `json:"email"`
	Notify   bool   `json:"notify"`
}
type Charge struct {
	Description  string `json:"description"`
	Amount       int64  `json:"amount"`
	PaymentTypes string `json:"paymentTypes"`
}

type Recharge struct {
	Charge  Charge  `json:"charge"`
	Billing Billing `json:"billing"`
}

type Response struct {
	Embedded struct {
		Charges []struct {
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
			Amount        float64 `json:"amount"`
			BilletDetails struct {
				BankAccount   string `json:"bankAccount"`
				BarcodeNumber string `json:"barcodeNumber"`
				OurNumber     string `json:"ourNumber"`
				Portfolio     string `json:"portfolio"`
			} `json:"billetDetails"`
			CheckoutURL     string `json:"checkoutUrl"`
			Code            int64  `json:"code"`
			DueDate         string `json:"dueDate"`
			ID              string `json:"id"`
			InstallmentLink string `json:"installmentLink"`
			Link            string `json:"link"`
			PayNumber       string `json:"payNumber"`
			Reference       string `json:"reference"`
			Status          string `json:"status"`
		} `json:"charges"`
	} `json:"_embedded"`
}
