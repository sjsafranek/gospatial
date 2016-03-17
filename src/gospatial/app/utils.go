package app

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// http://play.golang.org/p/4FkNSiUDMg
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(crand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// func UUID() string {
// 	b := make([]byte, 16)
// 	_, err := crand.Read(b)
// 	if err != nil {
// 		log.Fatalln("uuid error: ", Error.Error())
// 		return ""
// 	}
// 	b[6] = (b[6] & 0x0f) | 0x40
// 	b[8] = (b[8] & 0x3f) | 0x80
// 	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
// }

const letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewAPIKey(n int) string {
	s := ""
	for i := 1; i <= n; i++ {
		s += string(letters[rand.Intn(len(letters))])
		// Info.Println(letters[rand.Intn(len(letters))])
	}
	return s
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func sliceIndex(value string, slice []string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
