package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"unicode/utf8"
)

func LevelShift(tab int) string {
	res := ""
	for i := 0; i < tab; i += 1 {
		res = res + "\t"
	}
	return res
}

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

type Item struct {
	Type string
	Data string
}

//type ( [ {
func Load_level(level_attribute string, pos_begin int, pos_end int, lent int, text string, level int, flag string, debug int) ([]Item, int, int, int) {
	prev := pos_begin
	i_prev := pos_begin
	if debug > 39 {
		fmt.Printf("<-- text '%v'\r\n", text[pos_begin:pos_end])
	}
	flag_s := false
	g_s_flag := false
	items_a := []Item{}
	for i, w := pos_begin, 0; i < pos_end; {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		if debug > 42 {
			fmt.Printf("%v%#U starts at byte position %d %v level %v\n", LevelShift(level), runeValue, i, string(runeValue), level)
		}
		w = width
		s1 := string(runeValue)
		if !g_s_flag {
			// если это скобка - ищем ответную
			if s1 == "(" {
				if debug > 40 {
					fmt.Printf("i %v\r\n", i)
				}
				if prev < i {
					ci := Item{"symbols", text[prev:i]}
					items_a = append(items_a, ci)
				}
				level_g := 1
				prev = i
				i = i + w
				for {
					runeValue, width := utf8.DecodeRuneInString(text[i:])
					if debug > 42 {
						fmt.Printf("!%#U starts at byte position %d %v\n", runeValue, i, string(runeValue))
					}
					w = width
					s1 := string(runeValue)
					if s1 == "(" {
						level_g = level_g + 1
					} else {
						if s1 == ")" {
							level_g = level_g - 1
							if level_g > 0 {

							} else {
								break
							}
						}
					}
					i = i + w
					if i > len(text) {
						fmt.Printf("Error i %v len(text) %v\r\n", i, len(text))
						panic("Error!")
					}
				}
				if debug > 41 {
					fmt.Printf("text '%v' prev %v i %v w %v\r\n", text[prev:i+w], prev, i, w)
					fmt.Printf("prev %v, i+w %v, i-prev %v, text %v, level+1 %v\r\n", prev, i+w, i+w-prev, text, level+1)
				}
				ci := Item{"(", text[prev+1 : i]}
				items_a = append(items_a, ci)
				i = i + w
				prev = i
				i_prev = i
			} else if (s1 == " ") || (s1 == "\t") || (s1 == "\r") || (s1 == "\n") {
				// это разделитель! если до этого были отличные символы - строим строку.
				if debug > 40 {
					fmt.Printf("space flag_s %v\r\n", flag_s)
				}
				if flag_s {
					flag_s = false
					if debug > 41 {
						fmt.Printf("i_prev %v, i %v, w %v\r\n", i_prev, i, w)
					}
					if debug > 40 {
						fmt.Printf("text separator %v\r\n", text[i_prev:i])
					}
					if i_prev < i {
						ci := Item{"symbols", text[i_prev:i]}
						items_a = append(items_a, ci)
					}
				} else {
					// нет ничего это повторный разделитель
				}
				i = i + w
				i_prev = i
				prev = i
			} else if s1 == ")" {
				// скобка закрывающая завершаем работу и выходим
				if debug > 40 {
					fmt.Printf("close bracket flag_s %v\r\n", flag_s)
				}
				if flag_s {
					flag_s = false
					if debug > 41 {
						fmt.Printf("i_prev %v, i %v\r\n", i_prev, i)
					}
					if debug > 40 {
						fmt.Printf("text %v\r\n", text[i_prev:i])
					}
					ci := Item{")", text[i_prev:i]}
					items_a = append(items_a, ci)
				}
				return items_a, i + 1, pos_end, 0
			} else if s1 == "[" {
				if debug > 41 {
					fmt.Printf("i %v\r\n", i)
				}
				if prev < i {
					ci := Item{"symbols", text[prev:i]}
					items_a = append(items_a, ci)
				}
				level_g := 1
				prev = i
				i = i + w
				for {
					runeValue, width := utf8.DecodeRuneInString(text[i:])
					if debug > 43 {
						fmt.Printf("[>> !%#U starts at byte position %d %v\r\n", runeValue, i, string(runeValue))
					}
					w = width
					s1 := string(runeValue)
					if s1 == "[" {
						level_g = level_g + 1
					} else {
						if s1 == "]" {
							level_g = level_g - 1
							if level_g > 0 {

							} else {
								break
							}
						}
					}
					i = i + w
					if i > len(text) {
						fmt.Printf("Error i %v len(text) %v\r\n", i, len(text))
						panic("Error!")
					}
				}
				if debug > 41 {
					fmt.Printf("text [ '%v' prev %v i %v w %v\r\n", text[prev:i+w], prev, i, w)
					fmt.Printf("prev [ %v, i+w %v, i-prev %v, text %v, level+1 %v\r\n", prev, i+w, i+w-prev, text, level+1)
				}
				ci := Item{"[", text[prev+1 : i]}
				items_a = append(items_a, ci)
				i = i + w
				prev = i
				i_prev = i
			} else if s1 == "]" {
				// скобка закрывающая завершаем работу и выходим
				if debug > 40 {
					fmt.Printf("close bracket flag_s %v\r\n", flag_s)
				}
				if flag_s {
					if debug > 41 {
						fmt.Printf("i_prev %v, i %v\r\n", i_prev, i)
					}
					if debug > 40 {
						fmt.Printf("text ] %v\r\n", text[i_prev:i])
					}
					ci := Item{"]", text[i_prev:i]}
					items_a = append(items_a, ci)
				}
				return items_a, i + 1, pos_end, 0
			} else if s1 == "{" {
				if debug > 41 {
					fmt.Printf("i %v\r\n", i)
				}
				if prev < i {
					ci := Item{"symbols", text[prev:i]}
					items_a = append(items_a, ci)
				}
				level_g := 1
				prev = i
				i = i + w
				for {
					runeValue, width := utf8.DecodeRuneInString(text[i:])
					if debug > 40 {
						fmt.Printf("[>>> !%#U starts at byte position %d %v\r\n", runeValue, i, string(runeValue))
					}
					w = width
					s1 := string(runeValue)
					if s1 == "{" {
						level_g = level_g + 1
					} else {
						if s1 == "}" {
							level_g = level_g - 1
							if level_g > 0 {

							} else {
								break
							}
						}
					}
					i = i + w
					if i > len(text) {
						fmt.Printf("Error i %v len(text) %v\r\n", i, len(text))
						panic("Error!")
					}
				}
				if debug > 41 {
					fmt.Printf("text { '%v' prev %v i %v w %v\r\n", text[prev:i+w], prev, i, w)
					fmt.Printf("prev { %v, i+w %v, i-prev %v, text %v, level+1 %v\r\n", prev, i+w, i+w-prev, text, level+1)
				}
				ci := Item{"{", text[prev+1 : i]}
				items_a = append(items_a, ci)
				i = i + w
				prev = i
				i_prev = i
			} else if s1 == "}" {
				// скобка закрывающая завершаем работу и выходим
				if debug > 40 {
					fmt.Printf("close bracket flag_s %v\r\n", flag_s)
				}
				if flag_s {
					if debug > 41 {
						fmt.Printf("i_prev %v, i %v\r\n", i_prev, i)
					}
					if debug > 40 {
						fmt.Printf("text ] %v\r\n", text[i_prev:i])
					}
					ci := Item{"}", text[i_prev:i]}
					items_a = append(items_a, ci)
				}
				if debug > 40 {
				}
				return items_a, i + 1, pos_end, 0
			} else if s1 == "\"" {
				i_prev = i // + w
				if flag_s {
					flag_s = false
				} else {
					flag_s = true
				}
				g_s_flag = true
				i = i + w
			} else {
				// это символ!!
				if !flag_s {
					i_prev = i // + w
					flag_s = true
				} else {
				}
				i = i + w
			}
		} else {
			if s1 == "\"" {
				// это разделитель! если до этого были отличные символы - строим строку.
				if debug > 40 {
					fmt.Printf("space 2 flag_s %v\r\n", flag_s)
				}
				i = i + w
				if flag_s {
					flag_s = false
					if debug > 41 {
						fmt.Printf("i_prev %v, i %v, w %v\r\n", i_prev, i, w)
					}
					if debug > 40 {
						fmt.Printf("text %v\r\n", text[i_prev:i])
					}
					ci := Item{"string", text[i_prev:i]}
					items_a = append(items_a, ci)
				} else {
					// нет ничего это повторный разделитель
				}
				g_s_flag = false
				i_prev = i // + w
			} else {
				i = i + w
				flag_s = true
			}
		}
		if s1 == level_attribute {
			if !g_s_flag {
				if i_prev < i-w {
					ci := Item{"symbols", text[i_prev : i-w]}
					items_a = append(items_a, ci)
				}
				// возврат значения
				return items_a, i + 1, pos_end, 0
			}
		}
		if i >= pos_end {
			// все закончилось....
			if flag_s {
				flag_s = false
				if debug > 41 {
					fmt.Printf("i_prev %v, i %v\r\n", i_prev, i)
				}
				if debug > 40 {
					fmt.Printf("text all %v\r\n", text[i_prev:i])
				}
				ci := Item{"symbol", text[i_prev:i]}
				items_a = append(items_a, ci)
			}
		}
	}
	return items_a, 0, 0, 0
}

type GrammaticItem struct {
	Type       string
	Mod        string // "", "[0]" - нулевой элемент
	Attribute  int    // 0 none, 1 ==,
	Value      string
	GR_ID_List []string // список идентификатор граматического правила
	Ender      string   // здесь лежит окончатель
}
type GrammaticRule struct {
	ID                string // идентификатор
	GrammaticItemList []GrammaticItem
}

type CurrentEnv struct {
	Pi_cnt          int
	State           int
	Next_state      int
	Result_generate string
	Result          string
	I               int
	IntVars         map[string]int
	StringVars      map[string]string
}

type Env struct {
	map_gr           map[string]GrammaticRule
	base_gr_array    []GrammaticRule
	high_level_array []string
	expr_array       []string
	parse_func_dict  map[string]func(pi ParseItem, env *Env, level int) (string, error)
	ErrorsList       []string
	CE               *CurrentEnv
	Struct           interface{}
}

type ParseFuncDict struct {
	Dict map[string]ParseFunc
}

type ParseFunc func(pi ParseItem, env *Env, level int) (string, error)

func NewEnv() *Env {
	env := Env{}
	env.parse_func_dict = make(map[string]func(pi ParseItem, env *Env, level int) (string, error))
	env.map_gr = make(map[string]GrammaticRule)
	return &env
}

func (env *Env) SetEnv(mgr map[string]GrammaticRule, gr []GrammaticRule, hla []string, ea []string) {
	env.map_gr = mgr
	env.base_gr_array = gr
	env.high_level_array = hla
	env.expr_array = ea
}

func (env *Env) SetBGRAEnv() {
	for _, v := range env.map_gr {
		env.base_gr_array = append(env.base_gr_array, v)
	}
}

func (env *Env) SetHLAEnv(hla []string) {
	env.high_level_array = hla
}

func (env *Env) SetEAEnv(ea []string) {
	env.expr_array = ea
}

func (gr *GrammaticRule) AddRuleHandler(pf func(pi ParseItem, env *Env, level int) (string, error), env *Env) {
	_, ok := env.map_gr[gr.ID]
	if !ok {
		panic(fmt.Sprintf("Rule %v not exist", gr.ID))
	}
	env.parse_func_dict[gr.ID] = pf
	return
}

func MakeRule(rule_name string, env *Env) *GrammaticRule {
	gr, ok := env.map_gr[rule_name]
	if ok {
		panic(fmt.Sprintf("Rule name %v exist", rule_name))
	}
	gr = GrammaticRule{ID: rule_name, GrammaticItemList: []GrammaticItem{}}
	env.map_gr[rule_name] = gr
	return &gr
}

func (gri *GrammaticRule) AddItemToRule(Type string, Mod string, Attribute int, Value string, Ender string, GR_ID_List []string, env *Env) {
	gr, ok := env.map_gr[gri.ID]
	if !ok {
		panic(fmt.Sprintf("Rule %v not exist", gr.ID))
	}
	for i, _ := range GR_ID_List {
		if GR_ID_List[i] == gr.ID {
			// panic(fmt.Sprintf("Rule %v exist in GR_ID_List. Recurrent not support", gr.ID))
		}
	}
	gi := GrammaticItem{Type: Type, Mod: Mod, Attribute: Attribute, Value: Value, Ender: Ender, GR_ID_List: GR_ID_List}
	gr.GrammaticItemList = append(gr.GrammaticItemList, gi)
	env.map_gr[gr.ID] = gr
	return
}

func LoadGrammaticRule(env *Env, name string) (map[string]GrammaticRule, []GrammaticRule, error) {
	data, err := ioutil.ReadFile(name) // "settings.json"
	if err != nil {
		fmt.Print(err)
		return nil, nil, err
	}
	var gra []GrammaticRule
	err = json.Unmarshal(data, &gra)
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil, err
	}
	map_gr := make(map[string]GrammaticRule)
	for i, _ := range gra {
		map_gr[gra[i].ID] = gra[i]
	}
	env.map_gr = map_gr
	env.base_gr_array = gra
	return map_gr, gra, nil
}

func SaveGrammaticRule(env *Env, name string) error {
	gr := env.base_gr_array
	data_1, err2_ := json.MarshalIndent(&gr, "", "  ")
	if err2_ != nil {
		fmt.Println("error:", err2_)
		return err2_
	}
	_ = ioutil.WriteFile(name, data_1, 0644)
	return nil
}

type FuncDictItem struct {
	Name string
	Args int
}

type Variable struct {
	Name  string
	Value string
}
type ParseItem struct {
	Items     []Item
	Gr        *GrammaticRule
	PI        [][]ParseItem
	Variables []Variable
}

type Condition_mask struct {
	Base   string
	Length int
	Type   int
}

func ParseArgList(si string, debug int) []string {
	ender := ","

        s_in := si
	s_error := []string{}
	level := 0
	res_out := true
	ia := []string{}

	for {
		// читаем по грамматическим элементам
		l_1, pos_beg, pos_end, err1 := Load_level(ender, 0, len(s_in), len(s_in), s_in, 0, "", debug)
		if err1 != 0 {
			fmt.Printf("err %v\r\n", err1)
			s_error = append(s_error, fmt.Sprintf("error %v in level %v in process %v", err1, level, s_in))
			break
		}
		// ищем подходящий грамматический элемент
		if len(l_1) > 0 {
			if debug > 20 {
				fmt.Printf("l_1 %v\r\n", l_1)
			}
			for i, _ := range l_1 {
				ia = append(ia, l_1[i].Data)
			}
			if !res_out {
				break
			}
			if pos_beg >= pos_end {
				break
			}
			s_in = s_in[pos_beg:]
			s_in = strings.Trim(s_in, " \r\n\t")
			if len(s_in) == 0 {
				break
			}
		}
	}
	return ia
}

func load_items(in_file_name string, out_file_name string, env *Env) (string, error) {
	data, err := ioutil.ReadFile(in_file_name)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	s := string(data)

	// делим на элементы
	ll := strings.Split(s, "\n") // \r
	ll_n := []string{}
	for j, _ := range ll {
		// fmt.Printf("%v '%v'\r\n", j, ll[j])
		ls := strings.Trim(ll[j], " \r\n\t")
		if len(ls) > 0 {
			if (ls[0] == ';' && ls[1] == ';') || (ls[0] == '#') {
				// это комментарий
				// по идее можно добавить комментарий как отдельный оператор
			} else {
				// ищем в строке подстроку
				ls_l := strings.Split(ls, "//")
				if len(ls_l) > 1 {

					ls = ls_l[0]
				}
				// по идее надо добавлть признак строки а еще признак остановки при пошаговой отладке
				ll_n = append(ll_n, ls)
			}
		}
	}
	s = strings.Join(ll_n, "\r\n")

	check_rule := func(f_item []Item, gr_name_array []string, debug int) (*GrammaticRule, bool) {
		var res *GrammaticRule
		flag := false
		for _, gr_item_name := range gr_name_array {
			gr_item := env.map_gr[gr_item_name]
			if len(gr_item.GrammaticItemList) == len(f_item) {
				n := 0
				if debug > 2 {
					fmt.Printf("gr_item.ID %v gr_item %v\r\n", gr_item.ID, gr_item)
				}
				for i, _ := range gr_item.GrammaticItemList {
					gi := gr_item.GrammaticItemList[i]
					item := f_item[i]
					if gi.Type == item.Type {
						if debug > 2 {
							fmt.Printf("--- gi %v item %v\r\n", gi, item)
						}
						switch gi.Attribute {
						case 0:
							n = n + 1
						case 1:
							switch gi.Mod {
							case "":
								if gi.Value == item.Data {
									n = n + 1
								}
							case "[0]":
								if gi.Value == item.Data[0:1] {
									n = n + 1
								}
							}
						}
					} else {
						break
					}
				}
				if n == len(gr_item.GrammaticItemList) {
					res = &gr_item
					flag = true
					break
				}
			}
		}
		return res, flag
	}

	var print_pi func(pi ParseItem, level int)

	print_pi = func(pi ParseItem, level int) {
		pi_cnt := 0
		st := ""
		for k := 0; k < level; k++ {
			st = st + "\t"
		}
		fmt.Printf("%v%v:\r\n", st, pi.Gr.ID)
		for i, _ := range pi.Items {
			if len(pi.Gr.GrammaticItemList) > 0 {
				if len(pi.Gr.GrammaticItemList[i].GR_ID_List) > 0 {
					for j, _ := range pi.PI[pi_cnt] {
						print_pi(pi.PI[pi_cnt][j], level+1)
					}
					pi_cnt = pi_cnt + 1
				} else {
					fmt.Printf("%v\t%v \r\n", st, pi.Items[i])
				}
			} else {
				fmt.Printf("%v\t %v\r\n", st, pi.Items[i])
			}
		}
		fmt.Printf("\r\n")
	}

	var translate func(ender string, s_in string, gr_list []string, level int, debug int) ([]ParseItem, []string, bool)

	var generate_pi func(pi ParseItem, env_i *Env, level int) (string, bool, error)

	generate_pi = func(pi ParseItem, env_i *Env, level int) (string, bool, error) {
		result := ""
		st := ""
		for k := 0; k < level; k++ {
			st = st + "\t"
		}
		state := 0
		fn, ok := env_i.parse_func_dict[pi.Gr.ID]
		if !ok {
			return "", false, errors.New(fmt.Sprintf("handler for rule %v not defined", pi.Gr.ID))
		}
		flag := false
		result0 := ""
		for {
			switch state {
			case 0:
				env_i.CE.State = state
				res, err := fn(pi, env_i, level)
				if err != nil {
					return "", false, errors.New(fmt.Sprintf("Error when translate %v: %v", pi.Gr.ID, err))
				}
				result0 = res
				state = env_i.CE.State
			case 100:
				res_t := ""
				for j, _ := range pi.PI[env_i.CE.Pi_cnt] {
					ce_old := env.CE
					ce := CurrentEnv{}
					ce.Pi_cnt = env_i.CE.Pi_cnt // 0
					ce.State = -1
					ce.Next_state = -1
					ce.Result_generate = ""
					ce.Result = ""
					ce.I = 0
					ce.IntVars = make(map[string]int)
					ce.StringVars = make(map[string]string)
					env_i.CE = &ce
					res, status, err := generate_pi(pi.PI[env_i.CE.Pi_cnt][j], env_i, level+1)
					env_i.CE = ce_old
					if err != nil {
						return "", false, err
					}
					if status {
						res_t = res_t + res
					}
				}
				state = env_i.CE.Next_state
				env_i.CE.Result_generate = res_t
			case 200:
				// список аргументов
				pi_cnt := env_i.CE.Pi_cnt
				arg_lst := ""
				pi1 := pi.PI[pi_cnt][0]
				for _, it := range pi1.Items {
					res := it.Data
					if len(arg_lst) > 0 {
						arg_lst = arg_lst + ", " + res
					} else {
						arg_lst = arg_lst + res
					}
				}
				item_r := strings.Trim(arg_lst, " \r\n\t") + ","
				pia, e_list, res_ := translate(",", item_r, env_i.expr_array, level, 0)
				if res_ {
					if false {
						for i, _ := range pia {
							print_pi(pia[i], 0)
						}
					}
					result_ := ""
					for i, _ := range pia {
						ce_old := env.CE
						ce := CurrentEnv{}
						ce.Pi_cnt = 0
						ce.State = -1
						ce.Next_state = -1
						ce.Result_generate = ""
						ce.Result = ""
						ce.I = 0
						ce.IntVars = make(map[string]int)
						ce.StringVars = make(map[string]string)
						env.CE = &ce
						res_out, status, err := generate_pi(pia[i], env_i, 0)
						env.CE = ce_old
						if err != nil {
							fmt.Printf("Error while generate %v\r\n", pia[i])
							return "", false, err
						}
						if status {
							result_ = result_ + " " + res_out
						}
					}
					env_i.CE.Result_generate = result_
				} else {
					fmt.Printf("Error while translate %v\r\n", item_r)
					env_i.ErrorsList = append(env_i.ErrorsList, e_list...)
					env_i.ErrorsList = append(env_i.ErrorsList, []string{fmt.Sprintf("", "Error while translate %v\r\n", item_r)}...)
					return "", false, errors.New(fmt.Sprintf("Error while translate %v\r\n", item_r))
				}
				state = env_i.CE.Next_state
			case 1000:
				if len(result0) > 0 {
					result = result0
				} else {
				}
				flag = true
			default:
				env_i.CE.State = state
				res, err := fn(pi, env_i, level)
				if err != nil {
					return "", false, errors.New(fmt.Sprintf("Error when translate %v: %v", pi.Gr.ID, err))
				}
				result0 = ""
				result = res
				state = env_i.CE.State
			}
			if flag {
				break
			}
		}
		return result, true, nil
	}

	translate = func(ender string, s_in string, gr_list []string, level int, debug int) ([]ParseItem, []string, bool) {
		res_out := true
		s_error := []string{}
		pia := []ParseItem{}
		for {
			// читаем по грамматическим элементам
			l_1, pos_beg, pos_end, err1 := Load_level(ender, 0, len(s_in), len(s_in), s_in, 0, "", debug)
			if err1 != 0 {
				fmt.Printf("err %v\r\n", err1)
				s_error = append(s_error, fmt.Sprintf("error %v in level %v in process %v", err1, level, s_in))
				break
			}
			// ищем подходящий грамматический элемент
			if len(l_1) > 0 {
				if debug > 20 {
					fmt.Printf("l_1 %v\r\n", l_1)
				}
				gr, res := check_rule(l_1, gr_list, debug)
				if res {
					if debug > 20 {
						fmt.Printf("ID %v\r\n", gr.ID)
					}
					pi := ParseItem{Items: l_1, Gr: gr}
					if len(gr.GrammaticItemList) > 0 {
						for i, gr_i := range gr.GrammaticItemList {
							if len(gr_i.GR_ID_List) > 0 {
								s_ := l_1[i].Data
								flag_gr := false
								if len(gr_i.GR_ID_List) == 1 {
									ia := []Item{}
									id := ""
									switch gr_i.GR_ID_List[0] {
									case "список":
										item := Item{"symbols", strings.Trim(s_, " \r\n\t")}
										ia = append(ia, item)
										id = gr_i.GR_ID_List[0]
										flag_gr = true
									case "список_аргументов":
                                                                                al := ParseArgList(s_, debug)
										// al := strings.Split(s_, ",")
										for _, v := range al {
											item := Item{"symbols", strings.Trim(v, " \r\n\t")}
											ia = append(ia, item)
										}
										id = gr_i.GR_ID_List[0]
										flag_gr = true
									}
									if flag_gr {
										pi__ := ParseItem{Items: ia, Gr: &GrammaticRule{ID: id}}
										pi.PI = append(pi.PI, []ParseItem{pi__})
									}
								}
								if !flag_gr {
									pi__, e_list, res := translate(gr_i.Ender, s_, gr_i.GR_ID_List, level+1, debug)
									if res {
										pi.PI = append(pi.PI, pi__)
									} else {
										res_out = false
										s_error = append(s_error, e_list...)
										break
									}
								}
							}
						}
					}
					pia = append(pia, pi)
				} else {
					s_error = append(s_error, fmt.Sprintf("error %v in level %v in process %v", err1, level, l_1))
				}
				if !res_out {
					break
				}
				if pos_beg >= pos_end {
					break
				}
				s_in = s_in[pos_beg:]
				s_in = strings.Trim(s_in, " \r\n\t")
				if len(s_in) == 0 {
					break
				}
			} else {
				s_error = append(s_error, fmt.Sprintf("error %v in level %v in process %v", err1, level, s_in))
				res_out = false
				break
			}
		}
		return pia, s_error, res_out
	}
	pia, e_list, res := translate(";", s, env.high_level_array, 0, 0)
	if res {
		if true {
			for i, _ := range pia {
				print_pi(pia[i], 0)
			}
		}
		result := ""

		for i, _ := range pia {
			ce := CurrentEnv{}
			ce.Pi_cnt = 0
			ce.State = -1
			ce.Next_state = -1
			ce.Result_generate = ""
			ce.Result = ""
			ce.I = 0
			ce.IntVars = make(map[string]int)
			ce.StringVars = make(map[string]string)
			env.CE = &ce
			res_out, status, err := generate_pi(pia[i], env, 0)
			if err != nil {
				fmt.Printf("Error while translate %v\r\n", err)
				return "", err
			}
			if status {
				result = result + "\r\n" + res_out
			}
		}
		err = ioutil.WriteFile(out_file_name, []byte(result), 0644)
		if err != nil {
			panic(err)
		}
	} else {
		ss := ""
		for i, _ := range e_list {
			s := fmt.Sprintf("%v\r\n", e_list[i])
			ss = ss + s
		}
		return "", errors.New(ss)
	}
	return "", nil
}

func (env *Env) ParseFile(in_file_name string, out_file_name string) (string, error) {
	return load_items(in_file_name, out_file_name, env)
}
