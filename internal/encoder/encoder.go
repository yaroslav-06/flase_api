package encoder

import "crypto/sha256"

func Enc(val string) string {
	h := sha256.New()
	h.Write([]byte(val))
	return string(h.Sum(nil))
}
