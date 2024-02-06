package random

import "math/rand"

func GenerateRandomID() string {
	const roomIDLength = 4
	const charset = "ABCDEFG0123456789"
	roomID := make([]byte, roomIDLength)
	for i := range roomID {
		roomID[i] = charset[rand.Intn(len(charset))]
	}
	return string(roomID)
}
