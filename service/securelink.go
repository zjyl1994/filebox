package service

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"time"

	"github.com/zjyl1994/filebox/vars"
)

const SecureLinkExpireSecond = 300

func GenSecureLinkStr(path string) string {
	expireAt := time.Now().Add(time.Second * SecureLinkExpireSecond).Unix()
	sign := GenSecureLink(path, expireAt)
	return "s=" + sign + "&e=" + strconv.FormatInt(expireAt, 10)
}

func GenSecureLink(path string, expire int64) string {
	expireBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(expireBytes, uint64(expire))

	h := md5.New()
	h.Write([]byte(path))
	h.Write(expireBytes)
	h.Write([]byte(vars.Username))
	h.Write([]byte(vars.Password))
	result := h.Sum(nil)

	return base64.URLEncoding.EncodeToString(result)
}

func CheckSecureLink(path, token, expire string) bool {
	t := getTimeFromUnixtimestampString(expire)
	if time.Now().After(t) {
		return false
	}
	if GenSecureLink(path, t.Unix()) != token {
		return false
	}
	return true
}

func getTimeFromUnixtimestampString(t string) time.Time {
	ts, _ := strconv.ParseInt(t, 10, 64)
	return time.Unix(ts, 0)
}
