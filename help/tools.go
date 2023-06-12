package help

import (
	"crypto/md5"
	"fmt"
)

type Tools struct {
}

// 生成md5
func (t *Tools) GenerateMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
