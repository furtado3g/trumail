package verifier

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/mail"
	"strings"
)

// Address stores all information about an email Address
type Address struct {
	Address  string `json:"address" xml:"address"`
	Username string `json:"username" xml:"username"`
	Domain   string `json:"domain" xml:"domain"`
	MD5Hash  string `json:"md5Hash" xml:"md5Hash"`
}

// ParseAddress attempts to parse an email address and return it in the form
// of an Address struct pointer - domain case insensitive
func ParseAddress(email string) (*Address, error) {
	// Parses the address with the internal go mail address parser
	a, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}

	// Find the last occurrence of an @ sign
	index := strings.LastIndex(a.Address, "@")

	// Parse the username, domain and case unique address
	username := a.Address[:index]
	domain := strings.ToLower(a.Address[index+1:])
	address := fmt.Sprintf("%s@%s", username, domain)

	// Hash the address
	hashBytes := md5.Sum([]byte(address))
	md5Hash := hex.EncodeToString(hashBytes[:])

	// Returns the Address with the username and domain split out
	return &Address{address, username, domain, md5Hash}, nil
}
