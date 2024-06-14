package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type PasswordHasher interface {
	Hash(password string) string
	Verify(password, hashedPassword string) bool
}

type Hasher struct {
	secret string
}

func NewHasher(secret string) *Hasher {
	return &Hasher{secret: secret}
}

func (h *Hasher) Hash(password string) string {
	salt := hex.EncodeToString([]byte(uuid.NewString()))
	res := sha256.Sum256([]byte(salt + h.secret + password))
	return fmt.Sprintf("%x:%s", res, salt)
}

func (h *Hasher) Verify(password, hashedPassword string) bool {
	data := strings.Split(hashedPassword, ":")
	key, salt := data[0], data[1]
	res := sha256.Sum256([]byte(salt + h.secret + password))
	return key == fmt.Sprintf("%x", res)
}
