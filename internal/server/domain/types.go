package domain

type (
	RestResponse struct {
		ResponseMessage interface{} `json:"responseMessage,omitempty"`
		ResponseError   interface{} `json:"responseError,omitempty"`
	}

	Car struct {
		CarID              int      `json:"carID"`
		CarCompanyName     string   `json:"carCompanyName"`
		Doors              int      `json:"doors"`
		BigLuggage         int      `json:"bigLuggage,omitempty"`
		SmallLuggage       int      `json:"smallLuggage,omitempty"`
		AdultPlaces        int      `json:"adultPlaces"`
		AirConditioner     bool     `json:"airConditioner"`
		MinimumAge         int      `json:"minimumAge"`
		Price              int      `json:"price"`
		AvailableLocations []string `json:"availableLocations,omitempty"`
		CarGroup           int      `json:"carGroup"`
		Description        string   `json:"description"`
	}

	RentInfo struct {
		RentID          int      `json:"rentID"`
		CarID           int      `json:"carID"`
		FromDate        string   `json:"fromDate"`
		ToDate          string   `json:"toDate"`
		Location        string   `json:"location"`
		AvailableExtras []string `json:"availableExtras,omitempty"`
		Discounts       []string `json:"discounts,omitempty"`
		CarDetails      string   `json:"carDetails"`
		AgeGroup        string   `json:"ageGroup,omitempty"`
		CarGroup        int      `json:"carGroup,omitempty"`
	}

	CombinedRentInfo struct {
		Car
		RentID          int32    `json:"rentID,omitempty"`
		FromDate        string   `json:"fromDate,omitempty"`
		ToDate          string   `json:"toDate,omitempty"`
		Location        string   `json:"location,omitempty"`
		AvailableExtras []string `json:"availableExtras,omitempty"`
		Discounts       []string `json:"discounts,omitempty"`
		CarDetails      string   `json:"carDetails,omitempty"`
	}

	SearchParams struct {
		DateFilter    string
		DefaultFilter string
	}
)
