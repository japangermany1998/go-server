package auth

import (
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	listPasswords := []string{"123456", "nhatduc123"}

	for _, testCase := range listPasswords {
		hashed, err := HashPassword(testCase)
		if err != nil {
			t.FailNow()
		}

		if err = CheckPasswordHash(testCase, hashed); err != nil {
			t.FailNow()
		}
		if err = CheckPasswordHash(generateRandomString(), hashed); err == nil {
			t.FailNow()
		}
	}
}

func generateRandomString() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	charset := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	res := make([]byte, 4+rand.Intn(60))

	for i := 0; i < len(res); i++ {
		res[i] = charset[rand.Intn(len(charset))]
	}
	return string(res)
}

func TestValidateJWT(t *testing.T) {
	type test struct {
		userId              string
		tokenSecretMake     string
		tokenSecretValidate string
		expires             time.Duration
		sleep               time.Duration
		expectedErr         bool
	}
	cases := []test{{
		userId:              "123456",
		tokenSecretMake:     "abc",
		tokenSecretValidate: "abc",
		expires:             time.Second * 5,
		sleep:               time.Second * 2,
		expectedErr:         false,
	}}

	for _, testCase := range cases {
		signedString, err := MakeJWT(testCase.userId, testCase.tokenSecretMake, testCase.expires)
		if err != nil {
			t.FailNow()
		}
		time.Sleep(testCase.sleep)
		userId, err := ValidateJWT(signedString, testCase.tokenSecretValidate)
		if testCase.expectedErr != (err != nil) {
			t.FailNow()
		}
		if userId != testCase.userId {
			t.FailNow()
		}
	}
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}

	for i := 0; i < 10; i++ {
		token1 := generateRandomString()
		if i%2 == 0 {
			header.Add("Authorization", "Bearer "+token1)
		} else {
			header.Del("Authorization")
		}
		token2, err := GetBearerToken(header)
		if (i%2 == 0) == (err != nil) {
			t.FailNow()
		}
		if (i%2 == 0) == (token1 != token2) {
			t.FailNow()
		}
	}
}
