package new_strings

import (
	//	"fmt"
	"strings"
	"unicode/utf8"
)

func GetSlice(text string, begin_pos int, end_pos int) string {
	result := ""
	for i, w := 0, 0; i < len(text); i += w {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		w = width
		s1 := string(runeValue)
		if i >= begin_pos {
			if i < end_pos {
				result = result + s1
			}
		}
	}
	return result
}

func ParseStringBySignList(text string, sign_list []string) []string {
	// result := ""
	flag_s := false
	prev_i := 0
	ssl := []string{}
	ii := 0
	//	fmt.Printf("text %v\r\n", text)
	for i, w := 0, 0; i < len(text); i += w {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		w = width
		sl := string(runeValue)
		if flag_s {
			if sl == "\"" {
				flag_s = false
				ss := GetSlice(text, prev_i, i+w) //
				ss = strings.Trim(ss, " \r\n")
				ssl = append(ssl, ss)
				prev_i = i + w
			}
		} else {
			if sl == "\"" {
				if (i - prev_i) > 1 {

				} else {
					prev_i = i
				}
				//ss := GetSlice(text, prev_i, i)
				//ss = strings.Trim(ss, " \r\n")
				//ssl = append(ssl, ss)

				flag_s = true
			} else {
				for j, _ := range sign_list {
					if sl == sign_list[j] {
						if i-prev_i > 0 {
							ss := GetSlice(text, prev_i, i)
							ss = strings.Trim(ss, " \r\n")
							ssl = append(ssl, ss)
						}
						ssl = append(ssl, sl)
						prev_i = i + w
						break
					}
				}
			}
		}
		ii = i + w
	}
	//	fmt.Printf("- prev_i %v, ii %v\r\n", prev_i, ii)
	if (ii - prev_i) > 0 {
		//fmt.Printf("prev_i %v, ii %v\r\n", prev_i, ii)
		ss := GetSlice(text, prev_i, ii)
		ss = strings.Trim(ss, " \r\n")
		ssl = append(ssl, ss)
	}
	//	fmt.Printf("ssl %v\r\n", ssl)
	return ssl
}
