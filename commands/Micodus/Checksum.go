package micodus

func MicodusChecksum(message []byte) []byte {
	checksum := message[0]
	for i := 1; i < len(message); i++ {
		checksum = checksum ^ message[i]
	}
	return []byte{checksum}
}
