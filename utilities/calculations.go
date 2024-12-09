package utilities

func XorChecksum(message []byte) byte {
	checksum := message[0]
	for i := 1; i < len(message); i++ {
		checksum = checksum ^ message[i]
	}
	return checksum
}
