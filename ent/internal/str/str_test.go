package str_test

import (
	"testing"

	"github.com/gnames/gnparser/ent/internal/str"
	"github.com/stretchr/testify/assert"
)

func TestStringTools(t *testing.T) {
	t.Run("ToASCII", func(t *testing.T) {
		data := []struct {
			msg string
			in  string
			out string
			tbl map[rune]string
		}{
			{"Döringina", "Döringina", "Doeringina", str.Transliterations},
			{"Aëtosaurus", "Aëtosaurus", "Aetosaurus", str.Transliterations},
			{"thomæ", "thomæ", "thomae", str.Transliterations},
			{"many ö", "ööö", "oeoeoe", str.Transliterations},
			{"’", "’", "'", str.GlobalTransliterations},
			{"‘", "‘", "'", str.GlobalTransliterations},
			{"’", "’", "", str.Transliterations},
			{"‘", "‘", "", str.Transliterations},
		}
		for _, v := range data {
			res, _ := str.ToASCII([]byte(v.in), v.tbl)
			assert.Equal(t, string(res), v.out, v.msg)
		}
	})

	t.Run("NumToStr", func(t *testing.T) {
		data := []struct {
			msg string
			in  string
			out string
		}{
			{"1", "1", "uni"},
			{"2", "2", "bi"},
			{"3", "3", "tri"},
			{"4", "4", "quadri"},
			{"5", "5", "quinque"},
			{"6", "6", "sex"},
			{"7", "7", "septem"},
			{"8", "8", "octo"},
			{"9", "9", "novem"},
			{"10", "10", "decem"},
			{"11", "11", "undecim"},
			{"12", "12", "duodecim"},
			{"13", "13", "tredecim"},
			{"14", "14", "quatuordecim"},
			{"15", "15", "quindecim"},
			{"16", "16", "sedecim"},
			{"17", "17", "septendecim"},
			{"18", "18", "octodecim"},
			{"19", "19", "novemdecim"},
			{"20", "20", "viginti"},
			{"21", "21", "vigintiuno"},
			{"22", "22", "vigintiduo"},
			{"23", "23", "vigintitre"},
			{"24", "24", "vigintiquatuor"},
			{"25", "25", "vigintiquinque"},
			{"26", "26", "vigintisex"},
			{"27", "27", "vigintiseptem"},
			{"28", "28", "vigintiocto"},
			{"30", "30", "triginta"},
			{"31", "31", "trigintauno"},
			{"32", "32", "trigintaduo"},
			{"38", "38", "trigintaocto"},
			{"40", "40", "quadraginta"},
			{"400", "400", "400"},
			{"something", "something", "something"},
		}
		for _, v := range data {
			res := str.NumToStr(v.in)
			assert.Equal(t, res, v.out, v.msg)
		}
	})

	t.Run("FixAllCaps", func(t *testing.T) {
		data := []struct {
			msg string
			in  string
			out string
		}{
			{"KURNAKOV", "KURNAKOV", "Kurnakov"},
			{"GÓMEZ-BOLEA", "GÓMEZ-BOLEA", "Gómez-Bolea"},
			{"hello", "hello", "hello"},
		}
		for _, v := range data {
			res := str.FixAllCaps(v.in)
			assert.Equal(t, res, v.out, v.msg)
		}
	})
}
