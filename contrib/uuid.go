package contrib

// I found this code snippet on http://play.golang.org/p/4FkNSiUDMg
// Unfortunately I have no idea about who wrote it so I can't give
// them any credits etc (nor have I any idea of which license it
// is licensed under, but given that the go website in general
// have all code licensed under BSD license I'm assuming that
// that is the case for this as well :-S). Please let me know
// if you're the author (or know if the code is under another
// license and I'll take the appropriate actions).

import (
	"crypto/rand"
	"fmt"
	"io"
)

// newUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
