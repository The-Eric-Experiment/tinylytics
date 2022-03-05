package event

import (
	"time"

	"github.com/google/uuid"
)

func GetSessionUserIdent(item *ClientInfo) string {
	nm, _ := uuid.FromBytes([]byte("fdff1df9-3a7f-4451-8556-06c51d5d4fd1"))
	return uuid.NewSHA1(nm, []byte(item.UserAgent+item.Domain+item.HostName+item.IP)).String()
}

func GetSessionId(item *ClientInfo, startTime time.Time) string {
	nm, _ := uuid.FromBytes([]byte("6ddc50f6-86e9-4dde-a9c2-fa33f01d141d"))
	return uuid.NewSHA1(nm, []byte(item.UserAgent+item.Domain+item.HostName+item.IP+startTime.String())).String()
}
