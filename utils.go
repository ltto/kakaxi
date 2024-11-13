package kakaxi

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"unsafe"

	"github.com/shopspring/decimal"
)

func FormatBytes(size float64) string {
	units := []string{" B", " KB", " MB", " GB", " TB"}
	i := 0
	for i = 0; size >= 1024 && i < 4; i++ {
		res := decimal.NewFromFloat(size).Div(decimal.NewFromInt(1024))
		val, _ := res.Float64()
		size = val
	}

	result := strings.ReplaceAll(fmt.Sprintf("%.2f", FloatRound(size, 2)), "0000", "") + units[i]
	return result
}

func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func FileExist(savePath string) bool {
	_ = os.MkdirAll(path.Dir(savePath), 0777)
	//判断文件是否存在
	if _, err := os.Stat(savePath); err == nil {
		return true
	}
	return false
}

func Bytes2string(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func String2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
