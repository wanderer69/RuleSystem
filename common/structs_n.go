package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

type Variable struct {
	Name  string
	Value string
}

type Fact struct {
	ID         string
	Name       string
	State      int					// 0 - открыт для поиска 1 - захвачен 2 - удален
	Attributes []*Attribute          // перечень атрибутов факта
	Data       map[string]*Attribute // словарь значений атрибутов, где ключ имя атрибута а значение - значение атрибута
}

type Attribute struct {
	Type     string // константа, переменная, атрибут константы, атрибут переменной
	Const    string
	Variable string
	Value    string
	Code     int
}

func NewAttribute(type_ string, name string, value string) *Attribute {
	a := Attribute{}
	a.Type = type_
	switch a.Type {
	case "const":
		a.Const = name
	case "variable":
		a.Variable = name
	case "variable_value":
		a.Variable = name
		a.Value = value
	case "const_value":
		a.Const = name
		a.Value = value
	case "const_variable":
		a.Const = name
		a.Variable = value
	case "const_variable_value":
		a.Const = name
		a.Variable = value
	case "special":
		v, _ := strconv.ParseInt(value, 10, 64)
		a.Code = int(v)
	default:
		//fmt.Printf("Bad type %v\r\n", a.Type)
		panic(fmt.Sprintf("Bad type %v\r\n", a.Type))
	}
	return &a
}

func (a *Attribute) SetAttributeValue(value string) {
	switch a.Type {
	case "variable_value":
		a.Value = value
	case "const_value":
		a.Value = value
	case "const_variable":
	case "const_variable_value":
		a.Value = value
	}
}

func PrintAttribute(a *Attribute) string {
	result := ""
	if a != nil {
		switch a.Type {
		case "const":
			result = fmt.Sprintf("%v", a.Const)
		case "variable":
			result = fmt.Sprintf("?%v", a.Variable)
		case "variable_value":
			result = fmt.Sprintf("?%v:%v", a.Variable, a.Value)
		case "const_value":
			result = fmt.Sprintf("%v:%v", a.Const, a.Value)
		case "const_variable":
			result = fmt.Sprintf("%v:?%v", a.Const, a.Variable)
		case "special":
			result = fmt.Sprintf("%v", a.Code)
		}
	}
	return result
}

func CompareAttributes_f(a1 *Attribute, a2 *Attribute) (*Attribute, bool) {
	var result *Attribute
	if a1 != nil && a1 != nil {
		switch a1.Type {
		case "const":
			switch a2.Type {
			case "const":
				if a1.Const != a2.Const {
					return nil, false
				}
			case "variable":
				result = NewAttribute("variable_value", a2.Variable, a1.Const)
			case "variable_value":
				if a1.Const != a2.Value {
					return nil, false
				}
			case "const_value":
				if a1.Const != a2.Const {
					return nil, false
				}
			case "const_variable":
				if a1.Const != a2.Const {
					return nil, false
				}
			case "const_variable_value":
			}
		case "variable":
			switch a2.Type {
			case "const":
				result = NewAttribute("variable_value", a1.Variable, a2.Const)
			case "variable":
				if a1.Variable != a2.Variable {
					return nil, false
				}
			case "variable_value":
				if a1.Variable != a2.Variable {
					return nil, false
				}
			case "const_value":
				return nil, false
			case "const_variable":
				return nil, false
			case "const_variable_value":
			}
		case "variable_value":
			switch a2.Type {
			case "const":
				return nil, false
			case "variable":
				if a1.Variable != a2.Variable {
					return nil, false
				}
			case "variable_value":
				if a1.Variable != a2.Variable {
					return nil, false
				}
			case "const_value":
				if a1.Value != a2.Value {
					return nil, false
				}
				result = NewAttribute("variable_value", a1.Variable, a2.Const)
			case "const_variable":
				return nil, false
			case "const_variable_value":
			}
		case "const_value":
			switch a2.Type {
			case "const":
				if a1.Const != a2.Const {
					return nil, false
				}
			case "variable":
				return nil, false
			case "variable_value":
				if a1.Value != a2.Value {
					return nil, false
				}
				result = NewAttribute("variable_value", a2.Variable, a1.Const)
			case "const_value":
				if a1.Const != a2.Const {
					return nil, false
				}
				if a1.Value != a2.Value {
					return nil, false
				}
			case "const_variable":
				result = NewAttribute("variable_value", a1.Variable, a2.Const)
			case "const_variable_value":
			}
		case "const_variable":
			switch a2.Type {
			case "const":
				if a1.Const != a2.Const {
					return nil, false
				}
			case "variable":
				return nil, false
			case "variable_value":
				// have two results
				result = NewAttribute("variable_value", a1.Variable, a1.Value)
			case "const_value":
				result = NewAttribute("variable_value", a2.Variable, a1.Const)
			case "const_variable":
				return nil, false
			case "const_variable_value":
			}
		case "const_variable_value":
		}
	}
	return result, true
}

func CompareAttributes(a1 *Attribute, a2 *Attribute) (*Attribute, bool) {
	var result *Attribute
	if a1 != nil && a1 != nil {
		// сравниваем
		// если обе константы равны то истина
		// если константа и переменная то истина и у переменной атрибут константа
		// ecли константа с атрибутом равна константе с атрибутом то истина
		// если константа с атрибутом равна константе с переменной то истина и у переменной атрибут со значением атрибута константы
		switch a1.Type {
		case "const":
			switch a2.Type {
			case "const":
				if a1.Const != a2.Const {
					return nil, false
				}
			case "variable":
				result = NewAttribute("variable_value", a2.Variable, a1.Const)
			case "variable_value":
				return nil, false
			case "const_value":
				return nil, false
			case "const_variable":
				return nil, false
			case "const_variable_value":
				return nil, false
			}
		case "variable":
			switch a2.Type {
			case "const":
				result = NewAttribute("variable_value", a1.Variable, a2.Const)
			case "variable":
				return nil, false
			case "variable_value":
				return nil, false
			case "const_value":
				return nil, false
			case "const_variable":
				return nil, false
			case "const_variable_value":
				return nil, false
			}
		case "variable_value":
			switch a2.Type {
			case "const":
				return nil, false
			case "variable":
				return nil, false
			case "variable_value":
				return nil, false
			case "const_value":
				return nil, false
			case "const_variable":
				return nil, false
			case "const_variable_value":
				return nil, false
			}
		case "const_value":
			switch a2.Type {
			case "const":
				return nil, false
			case "variable":
				return nil, false
			case "variable_value":
				return nil, false
			case "const_value":
				if a1.Const != a2.Const {
					return nil, false
				}
				if a1.Value != a2.Value {
					return nil, false
				}
			case "const_variable":
				result = NewAttribute("const_variable_value", a2.Const, a1.Variable)
				result.SetAttributeValue(a2.Value)
			case "const_variable_value":
				if a1.Const != a2.Const {
					return nil, false
				}
				if a1.Value != a2.Value {
					return nil, false
				}
			}
		case "const_variable":
			switch a2.Type {
			case "const":
				return nil, false
			case "variable":
				return nil, false
			case "variable_value":
				return nil, false
			case "const_value":
				result = NewAttribute("const_variable_value", a1.Const, a2.Variable)
				result.SetAttributeValue(a1.Value)
			case "const_variable":
				return nil, false
			case "const_variable_value":
				result = NewAttribute("const_variable_value", a2.Const, a1.Variable)
				result.SetAttributeValue(a2.Value)
			}
		case "const_variable_value":
			switch a2.Type {
			case "const":
				return nil, false
			case "variable":
				return nil, false
			case "variable_value":
				return nil, false
			case "const_value":
				if a1.Const != a2.Const {
					return nil, false
				}
				if a1.Value != a2.Value {
					return nil, false
				}
			case "const_variable":
				result = NewAttribute("const_variable_value", a1.Const, a2.Variable)
				result.SetAttributeValue(a1.Value)
			case "const_variable_value":
				if a1.Const != a2.Const {
					return nil, false
				}
				if a1.Value != a2.Value {
					return nil, false
				}
			}
		}
	}
	return result, true
}

func CompareAttributeConst(a *Attribute, b string) (*Attribute, bool) {
	var result *Attribute
	if a != nil {
		// сравниваем
		// если константа и значение равны то истина
		// если переменная и значение равны то истина и у переменной атрибут значение
		switch a.Type {
		case "const":
			if a.Const != b {
				return nil, false
			}
		case "variable":
			result = NewAttribute("variable_value", a.Variable, b)
		case "variable_value":
			return nil, false
		case "const_value":
			return nil, false
		case "const_variable":
			return nil, false
		case "const_variable_value":
			return nil, false
		}
	}
	return result, true
}

func SetVariable(a1 *Attribute, a2 *Attribute) (*Attribute, bool) {
	var result *Attribute
	if a1 != nil && a1 != nil {
		// подставляем переменную - первый атрибут исходный, второй - список переменных
		// если переменная и переменная со значением - подставляем константу из значения и истинно
		// если константа с переменной и константа с переменной и атрибутом то истина подставляем константу с атрибутом
		switch a1.Type {
		case "const":
			return nil, false
		case "variable":
			switch a2.Type {
			case "const":
				return nil, false
			case "variable":
				return nil, false
			case "variable_value":
				result = NewAttribute("const", a2.Value, "")
			case "const_value":
				return nil, false
			case "const_variable":
				return nil, false
			case "const_variable_value":
				return nil, false
			}
		case "variable_value":
			return nil, false
		case "const_value":
			return nil, false
		case "const_variable":
			switch a2.Type {
			case "const":
				return nil, false
			case "variable":
				return nil, false
			case "variable_value":
				return nil, false
			case "const_value":
				return nil, false
			case "const_variable":
				return nil, false
			case "const_variable_value":
				result = NewAttribute("const_value", a1.Const, a2.Variable)
			}
		case "const_variable_value":
			return nil, false
		}
	}
	return result, true
}

/*
add - add fact to memory  (в первом аргументе имя переменной хранящей уникальный идентификатор факта, далее имя факта, далее список его атрибутов)
delete - fact from memory (в первом аргументе имя переменной хранящей идентификатор факта)
print - print string or fact (в первом аргументе строка или имя переменной)
call -
match - предикат првый аргумент которого имя переменной хранящей факт
*/

type Operator struct {
	Name       string
	Attributes []*Attribute          // перечень атрибутов факта
	Data       map[string]*Attribute // словарь значений атрибутов, где ключ имя атрибута а значение - значение атрибута
}

func NewOperator() Operator {
	o := Operator{}
	o.Data = make(map[string]*Attribute)
	return o
}

func PrintOperator(o Operator) string {
	result := fmt.Sprintf("%v (", o.Name)
	for i, _ := range o.Attributes {
		result = result + fmt.Sprintf("%v ", PrintAttribute(o.Attributes[i]))
	}
	result = result + fmt.Sprintf(")")
	return result
}

func PrintFact(f *Fact) string {
	result := fmt.Sprintf("%v (", f.ID)
	for i, _ := range f.Attributes {
		result = result + fmt.Sprintf("%v ", PrintAttribute(f.Attributes[i]))
	}
	result = result + fmt.Sprintf(")")
	return result
}

type Rule struct {
	ID           string      // идентификатор
	Name         string      // имя правила уникальное?
	Enabled      bool        //
	Conditions   []*Operator // список условий срабатывания правила
	Consequences []*Operator // список действий
}

func PrintRule(r Rule) string {
	result := fmt.Sprintf("Rule %v {\r\n", r.Name)
	for i, _ := range r.Conditions {
		c := r.Conditions[i]
		result = result + fmt.Sprintf("%v\r\n", PrintOperator(*c))
	}
	result = result + fmt.Sprintf(" => ")
	for i, _ := range r.Consequences {
		c := r.Consequences[i]
		result = result + fmt.Sprintf("%v\r\n", PrintOperator(*c))
	}
	result = result + fmt.Sprintf("}")
	return result
}

func NewFact() *Fact {
	f := Fact{}
	f.ID = "fact_" + Unique_Value(7)
	f.Data = make(map[string]*Attribute)
	f.State = 0
	return &f
}

type FactMemory struct {
	Facts    []*Fact
	FactDict map[string]*Fact
}

func NewFactMemory() *FactMemory {
	fm := FactMemory{}
	fm.FactDict = make(map[string]*Fact)
	return &fm
}

type Environment struct {
	Variables  map[string]*Attribute
	FactMemory *FactMemory
	Rules      []*Rule
	RulesDict  map[string]*Rule

	//MatchedFacts []*MatchedFact
}

func NewEnvironment() *Environment {
	env := Environment{}
	env.Variables = make(map[string]*Attribute)
	env.FactMemory = NewFactMemory()
	return &env
}

func (env *Environment) AddRule(r Rule) {
	r.Enabled = true
	env.Rules = append(env.Rules, &r)
}

func (env *Environment) AddFact(f *Fact) {
	for i, _ := range f.Attributes {
		switch f.Attributes[i].Type {
		case "const":
			f.Data[f.Attributes[i].Const] = f.Attributes[i]
		case "variable":
			f.Data[f.Attributes[i].Variable] = f.Attributes[i]
		case "variable_value":
			f.Data[f.Attributes[i].Variable] = f.Attributes[i]
		case "const_value":
			f.Data[f.Attributes[i].Const] = f.Attributes[i]
		case "const_variable":
			f.Data[f.Attributes[i].Const] = f.Attributes[i]
		case "const_variable_value":
			f.Data[f.Attributes[i].Const] = f.Attributes[i]
		case "special":
		}
	}
	env.FactMemory.Facts = append(env.FactMemory.Facts, f)
	env.FactMemory.FactDict[f.ID] = f
}

type Programm struct {
	Variables []*Variable
	Rules     []*Rule
	Facts     []*Fact
}

func Init_Unique_Value() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Unique_Value(len_n int) string {
	var bytes_array []byte

	for i := 0; i < len_n; i++ {
		bytes := rand.Intn(35)
		if bytes > 9 {
			bytes = bytes + 7
		}
		bytes_array = append(bytes_array, byte(bytes+16*3))
	}
	str := string(bytes_array)
	return str
}

func Programm2Environment(p *Programm) *Environment {
	env := NewEnvironment()
	for i, _ := range p.Rules {
		env.Rules = append(env.Rules, p.Rules[i])
	}

	for i, _ := range p.Facts {
		//		env.FactMemory.Facts[p.Facts[i].ID] = p.Facts[i]
		env.FactMemory.Facts = append(env.FactMemory.Facts, p.Facts[i])
	}
	return env
}

func LoadProgramm(file_name string) (*Programm, error) {
	data, err := ioutil.ReadFile(file_name)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	var p Programm

	err = json.Unmarshal(data, &p)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	return &p, nil
}

func SaveProgramm(p *Programm, file_name string) error {
	data_1, err2_ := json.MarshalIndent(&p, "", "  ")
	if err2_ != nil {
		fmt.Println("error:", err2_)
		return err2_
	}
	_ = ioutil.WriteFile(file_name, data_1, 0644)
	return nil
}

type MatchedFact struct {
	Fact       *Fact
	Attributes []*Attribute
}

type MatchContext struct {
	Current_fact int
	Facts        []*Fact
	Attributes   []*Attribute
}

func create_match_context(facts []*Fact, attributes []*Attribute) *MatchContext {
	facts_n := []*Fact{}
	for i, _ := range facts {
		if facts[i].State == 0 {
			facts_n = append(facts_n, facts[i])
		}
	}
	mc := MatchContext{0, facts_n, attributes}
	return &mc
}

func match_facts(mc *MatchContext, mfa []*MatchedFact) (*MatchedFact, error) {
	// нужно учитывать те переменные которые уже были заполнены.
	// алгоритм - проходим по списку и выбираем те факты которые совпали
	// строим список атрибутов на основании
	attributes := []*Attribute{}
	// надо проверять что запись уже была ассоциирована

	for i, _ := range mc.Attributes {
		res := false
		var a *Attribute
		for k, _ := range mfa {
			if len(mfa) > 0 {
				for j, _ := range mfa[k].Attributes {
					if mfa[k].Attributes[j] != nil {
						a, res = SetVariable(mc.Attributes[i], mfa[k].Attributes[j])
						if res {
							break
						}
					}
				}
				if res {
					break
				}
			}
		}
		if res {
			attributes = append(attributes, a)
		} else {
			attributes = append(attributes, mc.Attributes[i])
		}
	}
	mf := MatchedFact{}
	for {
		// выбираем текущий факт.
		// сравниваем атрибуты
		// все наоборот! цикл по списку атрибутов для сравнения!
		flag := true
		mal := []*Attribute{}
		if len(mc.Facts) == 0 {
                        flag = false
		        break
		}
		if mc.Facts[mc.Current_fact].State == 0 {
			if len(mc.Facts[mc.Current_fact].Attributes) == len(attributes) {
				for j, _ := range mc.Facts[mc.Current_fact].Attributes {
					// ищем список атрибутов по фактам
					a, res := CompareAttributes(mc.Facts[mc.Current_fact].Attributes[j], attributes[j])
					if res {
						mal = append(mal, a)
					} else {
						flag = false
						break
					}
				}
			}
		} else {
			flag = false
		}
		if flag {
			mf.Fact = mc.Facts[mc.Current_fact]
			mf.Fact.State = 1
			mf.Attributes = mal
			break
		}

		if mc.Current_fact < (len(mc.Facts) - 1) {
			mc.Current_fact = mc.Current_fact + 1
		} else {
			return nil, errors.New("Not matched")
		}
	}
	return &mf, nil
}

// исполнение
// выбираем правило init - выполняется один раз
// выполнение оператора
func ExecuteOperator(o *Operator, env *Environment, mfa []*MatchedFact, r *Rule) (bool, bool) {
	state := false
	result := false
	switch o.Name {
	case "match":
		result = true
	case "add":
		// добавление факта
		if o.Attributes[0].Type != "variable" {
			// ошибка ожидаем переменную

		} else {
			f := NewFact()
			attrs := []*Attribute{}
			al := o.Attributes[1:]
			for i, _ := range al {
				switch al[i].Type {
				case "const":
					attrs = append(attrs, al[i])
				case "variable":
					// ищем переменную
					flag := true
					for k, _ := range mfa {
						for j, _ := range mfa[k].Attributes {
							if mfa[k].Attributes[j] != nil {
								switch mfa[k].Attributes[j].Type {
								case "const":
								case "variable":
								case "variable_value":
									if al[i].Variable == mfa[k].Attributes[j].Variable {
										a := NewAttribute("const", mfa[k].Attributes[j].Value, "")
										attrs = append(attrs, a)
										flag = false
										break
									}
								case "const_value":
								case "const_variable":
								case "const_variable_value":
								case "special":
								}
								if !flag {
									break
								}
							}
						}
					}
					if flag {
						attrs = append(attrs, al[i])
					}
				case "variable_value":
					a := NewAttribute("const", al[i].Value, "")
					attrs = append(attrs, a)
				case "const_value":
					attrs = append(attrs, al[i])
				case "const_variable":
					// ищем переменную
					flag := true
					for k, _ := range mfa {
						for j, _ := range mfa[k].Attributes {
							if mfa[k].Attributes[j] != nil {
								switch mfa[k].Attributes[j].Type {
								case "const":
								case "variable":
								case "variable_value":
								case "const_value":
								case "const_variable":
								case "const_variable_value":
									if al[i].Const == mfa[k].Attributes[j].Const {
										if al[i].Variable == mfa[k].Attributes[j].Variable {
											a := NewAttribute("const_value", al[i].Const, mfa[k].Attributes[j].Value)
											attrs = append(attrs, a)
											flag = false
											break
										}
									}
								case "special":
								}
								if !flag {
									break
								}
							}
						}
					}
					if flag {
						attrs = append(attrs, al[i])
					}
				case "const_variable_value":
					attrs = append(attrs, al[i])
				case "special":
				}
			}
			ss := ""
			for i, _ := range attrs {
				ss = ss + fmt.Sprintf("%v ", PrintAttribute(attrs[i]))
			}

			f.Attributes = append(f.Attributes, attrs...)
			env.AddFact(f)
			o.Attributes[0].Value = f.ID
			env.Variables[o.Attributes[0].Variable] = o.Attributes[0]
		}
		result = true
	case "delete":
		if o.Attributes[0].Type == "variable" || o.Attributes[0].Type == "const" {
			fact_id := ""
			if o.Attributes[0].Type == "const" {
				if "self" == o.Attributes[0].Const {
					// надо взять текущие факты
				}
				fact_id = o.Attributes[0].Const
			} else {
				fact_id = o.Attributes[0].Value
			}
			f, ok := env.FactMemory.FactDict[fact_id]
			if !ok {
				// ошибка! нет такого факта
			} else {
				f.State = 2
			}
		} else {
			// ошибка ожидаем переменную
		}
		result = true
	case "print":
		if o.Attributes[0].Type == "variable" || o.Attributes[0].Type == "const" {
			if o.Attributes[0].Type == "const" {
				fmt.Printf("%v", o.Attributes[0].Const)
			} else {
			}
		}
		result = true
	case "call":
		result = true
	case "eq":
		result = true
	case "branch_if_false":
		result = true
	case "disable":
		if o.Attributes[0].Type == "variable" || o.Attributes[0].Type == "const" {
			if o.Attributes[0].Type == "const" {
				fmt.Printf("%v", o.Attributes[0].Const)
				if o.Attributes[0].Const == "self" {
					r.Enabled = false
				}
			} else {
			}
		}
		result = true
	case "enable":
		result = true
	case "quit":
		result = true
		state = true
	}
	return state, result
}

// выполнение правила идет в два этапа - во-первых определяется что правило применимо
// во-вторых - правило применяется
// если правил подходящих несколько то применение идет в цикле
func ExecuteRule(r *Rule, env *Environment) (bool, error) {
	// проверяем что правило разрешено?
	if !r.Enabled {
		return false, nil
	}
	// проверяем что есть подходящие факты
	flag := true
	//env.MatchedFacts = []*MatchedFact{}
	// для каждого оператора Match формируем отдельный контекст
	mca := []*MatchContext{}
	for i, _ := range r.Conditions {
		// формируем контекст получения списка
		mc := create_match_context(env.FactMemory.Facts, r.Conditions[i].Attributes)
		mca = append(mca, mc)
	}
	// делаем матч для каждого match. если не срабатывает - делаем возврат.
	// цель - найти такой вариант который удовлетворит всех
	for {
		current_match := 0
		mfa := []*MatchedFact{}

		if len(mca) > 0 {
			for {
			        // fmt.Printf("current_match %v mca %v\r\n", current_match, mca)
				mf, err := match_facts(mca[current_match], mfa)
				if err != nil {
					// спискок фактов закончился но сравнение завершено не было
					// если число найденных фактов не равно длине списка, то либо возвращаемся если
					// в current_match > 0 либо выходим с решением что фактов соответствующих условию нет
					if current_match == 0 {
						// первое правило не выполняется
						flag = false
						break
					} else {						
						current_match = current_match - 1
					}
				} else {
					// добавляем найденный факт и список сопоставленных переменных в список
					mfa = append(mfa, mf)
					if current_match < (len(r.Conditions) - 1) {
						current_match = current_match + 1
					} else {
						// результат достигнут
						break
					}
				}
			}
		} else {

		}
		if flag {
			// запускаем пометку фактов
			for i, _ := range mca {
				mca[i].Facts[mca[i].Current_fact].State = 2
			}
			// правило подходит. исполняем.
			for i, _ := range r.Consequences {
				state, res := ExecuteOperator(r.Consequences[i], env, mfa, r)
				if !res {
					// ошибка исполнения
					return false, errors.New(fmt.Sprintf("Error execution rule %v\r\n", r.Name))
				}
				if state {
					return false, nil
				}
			}
			if len(mca) == 0 {
				break
			} else {
				// меняем указатель на факт, так как факт отработан
			}
		} else {
			// запускаем распометку фактов
			for i, _ := range mca {
				if mca[i].Facts[mca[i].Current_fact].State == 1 {
					mca[i].Facts[mca[i].Current_fact].State = 0
				}
			}
			break
		}
	}
	return true, nil
}
