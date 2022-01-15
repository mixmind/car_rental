package cars

import (
	"car-rental/internal/server/domain"
	"math/rand"
)

/*
Generate new car with random parameters information
Car Company name:"Mercedes", "BMW", "Fiat", "Seat", "Kia", "Hyundai", "Renault", "Peugeot", "Suzuki"
Doors:3-5
Big Luggage:0-4
Small Luggage:0-5
Adult places:2-7
Air confitioner:true/false
Minimum Drivers Age:24-65
Price:5-500
Available Locations: "Jerusalem",		"Tel Aviv",		"Haifa",		"Ashdod",		"Rishon LeZiyyon",		"Petah Tikva",		"Beersheba",
					"Netanya",		"Holon",		"Bnei Brak",		"Rehovot",		"Bat Yam"
Car Group:1-5
Car Random decription:"Brand new car",		"Best choice for big family",		"Best choice for couples",		"Best choice for rich people",		"Small gas consumption",
							"Best choice for wild drivers",		"Best choice for new drivers",		"Best choice for city driving",		"Best choice for international trips"
*/
func GenerateNewCar(seed int64) domain.Car {
	rand.Seed(seed)

	mockCar := domain.Car{
		CarCompanyName:     domain.CarCompaniesList[randInt(0, len(domain.CarCompaniesList))],
		Doors:              randIntWithMax(3, 5),
		BigLuggage:         randIntWithMax(0, 4),
		SmallLuggage:       randIntWithMax(0, 5),
		AdultPlaces:        randIntWithMax(2, 7),
		AirConditioner:     rand.Int()%2 == 0,
		MinimumAge:         randIntWithMax(24, 65),
		Price:              randIntWithMax(5, 500),
		AvailableLocations: generateAdditionalFeatures(randInt(1, len(domain.CitiesList)), len(domain.CitiesList), domain.CitiesList),
		CarGroup:           randIntWithMax(1, 5),
		Description:        generateDescription(),
	}
	return mockCar
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randIntWithMax(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

/*
Generatate complex car information as array
*/
func generateAdditionalFeatures(numberOfObjects int, limitOfObjects int, objectsList []string) []string {
	var objectsArray []string
	objectsMap := make(map[string]bool, numberOfObjects)
	for i := 0; i < numberOfObjects; i++ {
		objectNumber := randInt(0, limitOfObjects)
		if _, exists := objectsMap[objectsList[objectNumber]]; !exists {
			objectsArray = append(objectsArray, objectsList[objectNumber])
			objectsMap[objectsList[objectNumber]] = true
		}
	}
	return objectsArray
}

/*
Generatate car description
*/
func generateDescription() string {
	return domain.CarDescriptionList[randInt(0, 8)]
}
