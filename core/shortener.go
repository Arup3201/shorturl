package core

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
	"time"
)

var B62Mapping = map[uint8]byte{
	0: '0', 1: '1', 2: '2', 3: '3', 4: '4', 5: '5', 6: '6', 7: '7', 8: '8', 9: '9',
	10: 'A', 11: 'B', 12: 'C', 13: 'D', 14: 'E', 15: 'F', 16: 'G', 17: 'H', 18: 'I', 19: 'J',
	20: 'K', 21: 'L', 22: 'M', 23: 'N', 24: 'O', 25: 'P', 26: 'Q', 27: 'R', 28: 'S', 29: 'T',
	30: 'U', 31: 'V', 32: 'W', 33: 'X', 34: 'Y', 35: 'Z',
	36: 'a', 37: 'b', 38: 'c', 39: 'd', 40: 'e', 41: 'f', 42: 'g', 43: 'h', 44: 'i', 45: 'j',
	46: 'k', 47: 'l', 48: 'm', 49: 'n', 50: 'o', 51: 'p', 52: 'q', 53: 'r', 54: 's', 55: 't',
	56: 'u', 57: 'v', 58: 'w', 59: 'x', 60: 'y', 61: 'z',
}

func genRandom10Digits() uint64 {
	a, _ := rand.Int(rand.Reader, big.NewInt(9999999998))
	return 1 + a.Uint64()
}

func getBase62(n uint64) string {
	x := n
	var r uint8
	var reversed []byte
	for x != 0 {
		r = uint8(x % 62)
		reversed = append(reversed, B62Mapping[r])
		x /= 62
	}

	var res strings.Builder
	for i := len(reversed) - 1; i >= 0; i-- {
		res.WriteByte(reversed[i])
	}

	return res.String()
}

func ShortenURL(ctx context.Context,
	longURL string,
	db *PostgresDB) (string, error) {

	random10Digit := genRandom10Digits()
	base62 := getBase62(random10Digit)
	shortUrl := "https://short.ly/" + base62

	err := db.Save(ctx, longURL, shortUrl, "active", time.Now().Add(365*24*time.Hour))
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}
