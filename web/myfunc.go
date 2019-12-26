package web

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func MyMd5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return fmt.Sprintf(hex.EncodeToString(m.Sum(nil)))
	/*
	   m := md5.Sum([]byte (s))
	   return fmt.Sprintf(hex.EncodeToString(m[:]))
	*/
	/*
	   m := md5.Sum([]byte(s))
	   fmt.Printf("%x", m)
	   return fmt.Sprintf()
	*/
	/*
	   m := md5.New()
	   io.WriteString(m, s)
	   return fmt.Sprintf(hex.EncodeToString(m.Sum(nil)))
	*/
}
