package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"
)

func gen() {
	timestamp := time.Now().Unix()

	s := sha1.New()
	fmt.Println(fmt.Sprintf("%s:%s:%d", "bitrise", "bitrise-shared-sso-secret", timestamp))
	s.Write([]byte(fmt.Sprintf("%s:%s:%d", "bitrise", "bitrise-shared-sso-secret", timestamp)))
	token := hex.EncodeToString(s.Sum(nil))
	fmt.Println(fmt.Sprintf("%d", timestamp))
	fmt.Println(token)
}

func main() {
	gen()
}
