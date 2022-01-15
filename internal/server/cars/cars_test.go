package cars

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCarGeneration(test *testing.T) {
	for i := 0; i < 100000; i++ {
		seed := int64(time.Now().UTC().UnixNano())
		mockCar := GenerateNewCar(seed)
		mockCar1 := GenerateNewCar(seed)
		car1, err := json.MarshalIndent(mockCar, "", "")
		if err != nil {
			test.Errorf("Failed to create json from car1:[%s]", err)
			test.Fail()
		}
		car2, err := json.MarshalIndent(mockCar1, "", "")
		if err != nil {
			test.Errorf("Failed to create json from car2:[%s]", err)
			test.Fail()
		}
		if !assert.Equal(test, mockCar, mockCar1, "Cars shoud be equal") {
			test.Errorf("Cars are not equal: car1 \n [%s] \n and car2 [%s]", string(car1), string(car2))
			test.Fail()
		} else {
			test.Logf("Cars are equal: car1 \n [%s] \n and car2 [%s]", string(car1), string(car2))
		}
	}

}
