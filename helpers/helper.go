package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"math"
)

func FuncAdd(i int) int {
	i += 1
	return i
}

func FuncMinus(i int) int {
	if i > 0 {
		i -= 1
		return i
	}
	return 0
}

func FuncRaw(s string) interface{} {
	return template.HTML(s)
}

func Md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Pager(page int, total int) (int, int, int, int) {
	limit := 10
	start := 0

	maxPage := int(math.Ceil(float64(total) / float64(limit)))

	if page <= 0 {
		page = 1
	} else if page > maxPage {
		page = maxPage
	}

	start = (page - 1) * limit
	end := start + limit

	if end > total {
		end = total
	}

	return start, end, page, maxPage
}
