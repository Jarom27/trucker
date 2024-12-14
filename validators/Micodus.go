package validators

import (
	"fmt"
	"log"
	"trucker/utilities"
)

type MicodusValidator struct {
}

func (validator MicodusValidator) Validate(message []byte) bool {
	validationResult := true
	if !validateChecksum(message[1 : len(message)-1]) {
		fmt.Println("Invalid checksum")
		return true
	}
	return validationResult
}
func validateChecksum(message_without_flags []byte) bool {
	log.SetPrefix("ValidateChecksum(): ")

	checksum := message_without_flags[len(message_without_flags)-1]
	calculatedChecksum := utilities.XorChecksum(message_without_flags[:len(message_without_flags)-1])

	log.Println(checksum, calculatedChecksum)

	return checksum == calculatedChecksum
}
