package guess_epsides

import (
	"fmt"
	"testing"
)

func TestMapOrders(t *testing.T) {
	orders, err := MapOrders([]string{
		"01 01",
		"01 02",
		"02 03",
		"02 04",
	})
	if err != nil {
		t.Error(err)
	}
	for k, v := range orders {
		if v != nil {
			fmt.Println(k, *v)
		} else {
			fmt.Println(k, "N")
		}
	}
}
