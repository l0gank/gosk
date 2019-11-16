package secure

import (
	"fmt"
	"hash"
	"strconv"
	"time"
)

// New initalizes security service
func New(h hash.Hash) *Service {
	return &Service{h: h}
}

// Service holds security related methods
type Service struct {
	h        hash.Hash
}

// Token generates new unique token
func (s *Service) Token(str string) string {
	s.h.Reset()
	fmt.Fprintf(s.h, "%s%s", str, strconv.Itoa(time.Now().Nanosecond()))
	return fmt.Sprintf("%x", s.h.Sum(nil))
}
