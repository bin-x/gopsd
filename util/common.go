package util

import (
	"encoding/binary"
	"errors"
	"os"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

var (
	BlendModeKeys = map[string]string{
		"pass": "Pass through", "norm": "Normal", "diss": "Dissolve",
		"dark": "Darken", "mul": "Multiply", "idiv": "Color burn",
		"lbrn": "Linear burn", "dkCl": "Darker color", "lite": "Lighten",
		"scrn": "Screen", "div": "Color dodge", "lddg": "Linear dodge",
		"lgCl": "Lighter color", "over": "Overlay", "sLit": "Soft light",
		"hLit": "Hard light", "vLit": "Vivid light", "lLit": "Linear light",
		"pLit": "Pin light", "hMix": "Hard mix", "diff": "Difference",
		"smud": "Exclusion", "fsub": "Subtract", "fdiv": "Divide",
		"hue": "Hue", "sat": "Saturation", "colr": "Color", "lum": "Luminosity",
	}
	ColorModes = map[int16]string{
		0: "Bitmap", 1: "Grayscale", 2: "Indexed", 3: "RGB",
		4: "CMYK", 7: "Multichannel", 8: "Duotune", 9: "Lab",
	}
)

func IsDocumentValid(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	b := make([]byte, 4)
	_, err = f.Read(b)
	if err != nil {
		return false, err
	}
	if string(b) != "8BPS" {
		return false, errors.New("Wrong document signature.")
	}

	b = make([]byte, 2)
	_, err = f.Read(b)
	if err != nil {
		return false, err
	}
	ver := binary.BigEndian.Uint16(b)
	if (strings.HasSuffix(path, "psd") && ver != 1) ||
		(strings.HasSuffix(path, "psb") && ver != 2) {
		return false, errors.New("Wrong document version.")
	}
	return true, nil
}

func InRange(i interface{}, min, max int) bool {
	val := getInteger(i)
	if val >= min && val <= max {
		return true
	}

	return false
}

func ValueIs(i interface{}, numbers ...int) bool {
	val := getInteger(i)
	for n := range numbers {
		if val == numbers[n] {
			return true
		}
	}
	return false
}

func StringValueIs(value string, values ...string) bool {
	for i := range values {
		if value == values[i] {
			return true
		}
	}
	return false
}

func getInteger(unk interface{}) int {
	switch i := unk.(type) {
	case int64:
		return int(i)
	case int32:
		return int(i)
	case int16:
		return int(i)
	case byte:
		return int(i)
	case int:
		return i
	default:
		return 0
	}
}

func UnpackRLEBits(data []int8, length int) []int8 {
	result := make([]int8, length)
	wPos, rPos := 0, 0
	for rPos < len(data) {
		n := data[rPos]
		if n > 0 {
			count := int(n) + 1
			for j := 0; j < count; j++ {
				result[wPos] = data[rPos]
				wPos++
				rPos++
			}
		} else {
			b := data[rPos]
			rPos++
			count := int(-n) + 1
			for j := 0; j < count; j++ {
				result[wPos] = b
				wPos++
			}
		}
	}
	return result
}

func BytesToUTF16(b []byte, o binary.ByteOrder) string {
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}
	return string(utf16.Decode(utf))
}
