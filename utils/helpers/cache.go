package helper

import (
	"encoding/base64"
	"fmt"
	"myself_framwork/configs"
)

func RemoveCacheBalance(customerId string) bool {
	encodedString := base64.StdEncoding.EncodeToString([]byte(customerId))
	key_first := fmt.Sprintf("%s-res-%s", "balance", encodedString)
	key_two := fmt.Sprintf("%s-exp-%s", "balance", encodedString)
	fmt.Printf("remove %s\n", customerId)
	fmt.Printf("remove %s\n", encodedString)

	if configs.HasKey(key_first) {
		configs.RemoveKey(key_first)
	}

	if configs.HasKey(key_two) {
		configs.RemoveKey(key_two)
	}

	return true
}
