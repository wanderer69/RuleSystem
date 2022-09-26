package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	rs "github.com/wanderer69/RuleSystem/common"
	. "github.com/wanderer69/RuleSystem/parser"
)

// грамматика !!!
/* верхнеуровневые элементы
язык правил
define rule:
defrule <rule name> {
# comment
    # condition
    condition {
    # match(<fact name>, <attribute name 1>:<value>, <attribute name 2>:<value>, ...)
    # match select from
    match(сказуемое, вариант:?вариант, значение:?значение, часть_речи:глагол, склонение:?склонение, лицо:?лицо)
    # can be selected many facts
    };
    # =>
    #consequence
    consequence {
    # циклическая обработка для каждого факта из найденного списка
    # текущий факт в переменной ?current_fact
    # обращение к атрибутам текущего факта:
    # 	уникальный идентификатор - ?current_fact:ID
    # 	имя факта ?current_fact:Name
    #	к пользовательскому атрибуту факта ?current_fact:<имя атрибута>
    # добавляет новый факт. первый аргумент - переменная в которую пападет новый факт
    add(<variable fact>, <fact name>, <attribute name 1>:<value>, <attribute name 2>:<value>, ...);
    # удаляет факт из памяти по идентификатору
    delete(<fact id>);
    # печатает факт по идентификатору
    print(<fact id>);
    # вызов применения правила
    call(<rule name>);
    # предикат равенства значения атрибута и константы
    ?<variable name>:<attribute name> == <const>
    # предикат равенства значения атрибута и значения атрибута
    ?<variable name 1>:<attribute name> == ?<variable name 2>:<attribute name>
    # if <predicate> { <list of operators> };
    };
};
<symbols, == defrule> <symbols, >  <{, > - определение правила
<symbols, == condition> <{, > - определение секции условий
<symbols, == consquence> <{, > - определение секции следствий
<symbols, == match> <(, > - определение оператора match
<symbols, == add> <(, > - определение оператора add
<symbols, == delete> <(, > - определение оператора delete
<symbols, == print> <(, > - определение оператора print
<symbols, == disable> <(, > - определение оператора disable - выключить правило
<symbols, == enable> <(, > - определение оператора enable - включить правило
<symbols, == delete> <(, > - определение оператора delete - удалить факт
<symbols, == call> <(, > - определение оператора call - вызвать правило по имени
<symbols, == quit> <(, > - определение оператора quit - закончить выполнение
<symbols, == if> <(, >  <{, > - если
<symbols, > <symbols, == ==>  <symbols, "строка"> - условие 1 if
<symbols, > <symbols, == ==>  <symbols, > - условие 2 if


<symbols, [0] == "?"> <symbols, == => > <symbols,> - переприсваивание значения переменой другой переменной
<string, > <symbols, == => > <symbols,> - константу строку в переменную
<(, > <symbols, == => > <symbols,> - константу список в переменную
<[, > <symbols, == => > <symbols,> - константу массив в переменную
<{, > <symbols, == => > <symbols,> - константу словарь в переменную
<symbols, == Факт> <(, >  <symbols, == => > <symbols,> - определение факта
<symbols, == Шаблон> <(, >  <symbols, == => > <symbols,> - определение шаблона
<symbols, == Лисп> <(, > - вставка на чистом Лиспе
<symbols, > <(, > - вызов функции без возврата значения
<symbols, > <symbols, == => > <symbols,>  - вызов функции с возвратом значения
среднеуровневые элементы
операторы
<symbols, == Если> <(, >  <{, > - если
<symbols, == Если> <(, >  <{, > <symbols, == Иначе> <{, > - если иначе
<symbols, == Цикл> <symbols, == по> <symbols, [0] == "?">  <symbols, == => > <symbols,> <{, > - Цикл по
<symbols, == Вернуть> <symbols, > - вернуть
список в обпределении тринара или шаблона
<symbols, > <symbols, > <symbols, > - список тринаров

*/

type RuleParserStackItem struct {
	ConditionOps []*rs.Operator
	ExecOps      []*rs.Operator
}

type RuleParser struct {
	CurOp        rs.Operator
	Conditions   []*rs.Operator // список условий срабатывания правила
	Consequences []*rs.Operator // список действий

	Stack    []*RuleParserStackItem
	StackPos int

	Env *rs.Environment
}

func ParseArg(val string) (*rs.Attribute, error) {
	if len(val) > 0 {
		if val[0] == '?' {
			// переменная
			vval := val[1:]
			lvval := strings.Split(vval, ":")
			if len(lvval) > 0 {
				if len(lvval) == 2 {
					var_name := lvval[0]
					value := lvval[1]
					a := rs.NewAttribute("variable_value", var_name, value)
					return a, nil
				} else {
					if len(lvval) > 2 {
						// error!
						return nil, errors.New("Many symbols :")
					}
					a := rs.NewAttribute("variable", vval, "")
					return a, nil
				}
			} else {
				return nil, errors.New("Too small symbols")
			}
		} else {
			// константа
			lvval := strings.Split(val, ":")
			if len(lvval) > 0 {
				if len(lvval) == 2 {
					const_name := lvval[0]
					value := lvval[1]
					a := rs.NewAttribute("const_value", const_name, value)
					return a, nil
				} else {
					if len(lvval) > 2 {
						// error!
						return nil, errors.New("Many symbols :")
					}
					a := rs.NewAttribute("const", val, "")
					return a, nil
				}
			} else {
				return nil, errors.New("Too small symbols")
			}
		}
	}
	return nil, nil
}

func f_defrule(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		rule_name := pi.Items[1].Data
		env.CE.StringVars["rule_name"] = rule_name
		// список операторов
		env.CE.Pi_cnt = 0
		env.CE.Next_state = 1
		env.CE.State = 100

		rp := env.Struct.(RuleParser)
		//CurOp rs.Operator
		rp.Conditions = []*rs.Operator{}
		rp.Consequences = []*rs.Operator{}
		rp.Stack = []*RuleParserStackItem{}
		rp.StackPos = -1
		env.Struct = rp

	case 1:
		body := env.CE.Result_generate
		rule_name := env.CE.StringVars["rule_name"]
		result = fmt.Sprintf("(defrule %v %v)", rule_name, body)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		r := rs.Rule{}
		uv := rs.Unique_Value(10)
		r.ID = uv
		r.Name = rule_name
		r.Conditions = rp.Conditions
		r.Consequences = rp.Consequences
		rp.Env.AddRule(r)
		env.Struct = rp
	}
	return result, nil
}

func f_condition(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// список операторов
		env.CE.Pi_cnt = 0
		env.CE.Next_state = 1
		env.CE.State = 100
	case 1:
		body := env.CE.Result_generate
		result = fmt.Sprintf("(condition %v)", body)
		env.CE.State = 1000
	}
	return result, nil
}

func f_consequence(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// список операторов
		env.CE.Pi_cnt = 0
		env.CE.Next_state = 1
		env.CE.State = 100
	case 1:
		body := env.CE.Result_generate
		result = fmt.Sprintf("(consequence %v)", body)
		env.CE.State = 1000
	}
	return result, nil
}

func f_match(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// список аргументов
		env.CE.Pi_cnt = 0
		env.CE.Next_state = 1
		env.CE.State = 200
	case 1:
		body := env.CE.Result_generate
		result = fmt.Sprintf("(match %v)", body)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "match"

		b := strings.Trim(body, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := strings.Trim(args[i], " ")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}
				op.Attributes = append(op.Attributes, a)
			}
		}
		rp.Conditions = append(rp.Conditions, &op)
		env.Struct = rp
	}
	return result, nil
}

func f_add(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// список аргументов
		env.CE.Pi_cnt = 0
		env.CE.Next_state = 1
		env.CE.State = 200
	case 1:
		body := env.CE.Result_generate
		result = fmt.Sprintf("(add %v)", body)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "add"

		b := strings.Trim(body, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := strings.Trim(args[i], " ")
			if len(arg) > 0 {
				a, err := ParseArg(arg)
				if err != nil {
					return "", err
				}

				op.Attributes = append(op.Attributes, a)
			}
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func f_delete(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		fact_name := pi.Items[1].Data
		result = fmt.Sprintf("(delete %v)", fact_name)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "delete"

		b := strings.Trim(fact_name, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func f_print(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// значение - константа или значение переменной
		value := pi.Items[1].Data
		result = fmt.Sprintf("(print %v)", value)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "print"

		b := strings.Trim(value, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func f_quit(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// значение - константа или значение переменной
		value := pi.Items[1].Data
		result = fmt.Sprintf("(quit %v)", value)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "quit"

		b := strings.Trim(value, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func f_enable(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		rule_name := pi.Items[1].Data
		result = fmt.Sprintf("(enable %v)", rule_name)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "enable"

		b := strings.Trim(rule_name, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func f_disable(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	fmt.Printf("disable\r\n")
	switch env.CE.State {
	case 0:
		// имя правила
		rule_name := pi.Items[1].Data
		result = fmt.Sprintf("(disable %v)", rule_name)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "disable"

		b := strings.Trim(rule_name, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func f_call(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		rule_name := pi.Items[1].Data
		result = fmt.Sprintf("(call %v)", rule_name)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "call"

		b := strings.Trim(rule_name, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ExecOps = append(rp.Stack[rp.StackPos].ExecOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}

		env.Struct = rp
	}
	return result, nil
}

func f_empty(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// имя правила
		rule_name := pi.Items[1].Data
		result = fmt.Sprintf("(empty %v)", rule_name)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "empty"

		b := strings.Trim(rule_name, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}
		env.Struct = rp
	}
	return result, nil
}

func f_if(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// список аргументов
		env.CE.Pi_cnt = 0
		env.CE.Next_state = 1
		env.CE.State = 100

		rp := env.Struct.(RuleParser)
		si := RuleParserStackItem{}

		rp.Stack = append(rp.Stack, &si)
		rp.StackPos = rp.StackPos + 1

		env.Struct = rp
	case 1:
		env.CE.StringVars["condition"] = env.CE.Result_generate
		env.CE.Pi_cnt = 1
		env.CE.Next_state = 2
		env.CE.State = 100
	case 2:
		body := env.CE.Result_generate
		cond := env.CE.StringVars["condition"]
		result = fmt.Sprintf("(if %v %v)", cond, body)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		// добавляем оператор условия
		cops := rp.Stack[rp.StackPos].ConditionOps

		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, cops...)
		} else {
			rp.Consequences = append(rp.Consequences, cops...)
		}
		// добавляем переход
		eops := rp.Stack[rp.StackPos].ExecOps
		l := len(eops)

		op := rs.Operator{}
		op.Name = "branch_if_false"
		a := rs.NewAttribute("special", "", fmt.Sprintf("%v", l))
		op.Attributes = append(op.Attributes, a)

		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}
		// добавляем выполняемые в случае исполнения условия команды
		if rp.StackPos > 0 {
			rp.Stack[rp.StackPos-1].ExecOps = append(rp.Stack[rp.StackPos-1].ExecOps, eops...)
		} else {
			rp.Consequences = append(rp.Consequences, eops...)
		}

		rp.StackPos = rp.StackPos - 1
		if len(rp.Stack) > 0 {
			rp.Stack = rp.Stack[:len(rp.Stack)-1]
		} else {
			rp.Stack = []*RuleParserStackItem{}
		}
		env.Struct = rp
	}
	return result, nil
}

func f_condition1(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// левая часть условия
		l := pi.Items[0].Data
		r := pi.Items[2].Data
		result = fmt.Sprintf("(eq %v %v)", l, r)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "eq"

		rule_name := l + " " + r
		b := strings.Trim(rule_name, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ConditionOps = append(rp.Stack[rp.StackPos].ConditionOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}
		env.Struct = rp
	}
	return result, nil
}

func f_condition2(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// левая часть условия
		l := pi.Items[0].Data
		r := pi.Items[2].Data
		result = fmt.Sprintf("(eq %v %v)", l, r)
		env.CE.State = 1000

		rp := env.Struct.(RuleParser)
		op := rs.Operator{}
		op.Name = "eq"

		rule_name := l + " " + r
		b := strings.Trim(rule_name, " ")
		args := strings.Split(b, " ")
		for i, _ := range args {
			arg := args[i]
			a, err := ParseArg(arg)
			if err != nil {
				return "", err
			}
			op.Attributes = append(op.Attributes, a)
		}

		if rp.StackPos >= 0 {
			rp.Stack[rp.StackPos].ConditionOps = append(rp.Stack[rp.StackPos].ConditionOps, &op)
		} else {
			rp.Consequences = append(rp.Consequences, &op)
		}

		env.Struct = rp
	}
	return result, nil
}

func f_symbol(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// символ
		s := pi.Items[0].Data
		result = fmt.Sprintf(" %v", s)
		env.CE.State = 1000
	}
	return result, nil
}

func f_string(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// строка
		s := pi.Items[0].Data
		result = fmt.Sprintf(" %v", s)
		env.CE.State = 1000
	}
	return result, nil
}

func f_variable(pi ParseItem, env *Env, level int) (string, error) {
	result := ""
	switch env.CE.State {
	case 0:
		// переменная
		s := pi.Items[1].Data
		result = fmt.Sprintf("var %v", s)
		env.CE.State = 1000
	}
	return result, nil
}

func Make_Rules(env *Env) {
	if true {
		defer func() {
			r := recover()
			if r != nil {
				fmt.Printf("%v\r\n", r)
				return
			}
		}()
	}
	// <symbols, == defrule> <symbols, >  <{, > - определение правила
	gr := MakeRule("правило", env)
	gr.AddItemToRule("symbols", "", 1, "defrule", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", []string{"условие", "следствие"}, env)
	gr.AddRuleHandler(f_defrule, env)

	// <symbols, == condition> <{, > - определение секции условий
	gr = MakeRule("условие", env)
	gr.AddItemToRule("symbols", "", 1, "condition", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", []string{"найти факт", "пусто"}, env)
	gr.AddRuleHandler(f_condition, env)

	// <symbols, == consquence> <{, > - определение секции следствий
	gr = MakeRule("следствие", env)
	gr.AddItemToRule("symbols", "", 1, "consequence", "", []string{}, env)
	gr.AddItemToRule("{", "", 0, "", ";", []string{"добавить факт",
		"удалить факт", "печатать факт", "вызвать правило", "если",
		"выключить правило", "включить правило", "завершить выполнение"}, env)
	gr.AddRuleHandler(f_consequence, env)

	// <symbols, == match> <(, > - определение оператора match
	gr = MakeRule("найти факт", env)
	gr.AddItemToRule("symbols", "", 1, "match", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"список_аргументов"}, env)
	gr.AddRuleHandler(f_match, env)

	// <symbols, == add> <(, > - определение оператора add
	gr = MakeRule("добавить факт", env)
	gr.AddItemToRule("symbols", "", 1, "add", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"список_аргументов"}, env)
	gr.AddRuleHandler(f_add, env)

	// <symbols, == delete> <(, > - определение оператора delete
	gr = MakeRule("удалить факт", env)
	gr.AddItemToRule("symbols", "", 1, "delete", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"символ"}, env) // "символ"
	gr.AddRuleHandler(f_delete, env)

	// <symbols, == disable> <(, > - определение оператора disable - выключить правило
	gr = MakeRule("выключить правило", env)
	gr.AddItemToRule("symbols", "", 1, "disable", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"символ"}, env) // "символ"
	gr.AddRuleHandler(f_disable, env)

	// <symbols, == enable> <(, > - определение оператора enable - включить правило
	gr = MakeRule("включить правило", env)
	gr.AddItemToRule("symbols", "", 1, "enable", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"символ"}, env) // "символ"
	gr.AddRuleHandler(f_enable, env)

	// <symbols, == print> <(, > - определение оператора print
	gr = MakeRule("печатать факт", env)
	gr.AddItemToRule("symbols", "", 1, "print", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"список_аргументов"}, env)   // "символ"
	gr.AddRuleHandler(f_print, env)

	// <symbols, == quit> <(, > - определение оператора quit - закончить выполнение
	gr = MakeRule("завершить выполнение", env)
	gr.AddItemToRule("symbols", "", 1, "quit", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"символ"}, env) // "символ"
	gr.AddRuleHandler(f_quit, env)

	// <symbols, == call> <(, > - определение оператора call
	gr = MakeRule("вызвать правило", env)
	gr.AddItemToRule("symbols", "", 1, "call", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"символ"}, env)
	gr.AddRuleHandler(f_call, env)

	// <symbols, == empty> <(, > - определение оператора empty
	gr = MakeRule("пусто", env)
	gr.AddItemToRule("symbols", "", 1, "empty", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"символ"}, env)
	gr.AddRuleHandler(f_empty, env)

	// <symbols, == if> <(, >  <{, > - если
	gr = MakeRule("если", env)
	gr.AddItemToRule("symbols", "", 1, "if", "", []string{}, env)
	gr.AddItemToRule("(", "", 0, "", "", []string{"условие1", "условие2"}, env)
	gr.AddItemToRule("{", "", 0, "", ";", []string{"добавить факт",
		"удалить факт", "печатать факт", "вызвать правило", "если"}, env)
	gr.AddRuleHandler(f_if, env)

	//<symbols, > <symbols, == ==>  <строка> - условие 1 if
	gr = MakeRule("условие1", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "==", "", []string{}, env)
	gr.AddItemToRule("string", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(f_condition1, env)

	//<symbols, > <symbols, == ==>  <symbols, > - условие 2 if
	gr = MakeRule("условие2", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "==", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(f_condition2, env)

	// среднеуровневые элементы
	// список в определении тринара или шаблона
	// <symbols, > - просто символ
	gr = MakeRule("символ", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(f_symbol, env)

	// <string, > - просто строка
	gr = MakeRule("строка", env)
	gr.AddItemToRule("string", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(f_string, env)

	// ?<variable name>
	// <symbols, == ?> - переменная
	gr = MakeRule("переменная", env)
	gr.AddItemToRule("symbols", "", 0, "?", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddRuleHandler(f_variable, env)

	// ?<variable name>:<attribute name>
	// <symbols, == ?> <symbols, > <symbols, == :> <symbols, >- атрибут переменной
	gr = MakeRule("атрибут переменной", env)
	gr.AddItemToRule("symbols", "", 0, "?", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, ":", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)

	// ?<variable name>:<attribute name>
	// <symbols, > <symbols, == :> <symbols, >- атрибут
	gr = MakeRule("атрибут", env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, ":", "", []string{}, env)
	gr.AddItemToRule("symbols", "", 0, "", "", []string{}, env)

	high_level_array := []string{"правило"}

	expr_array := []string{"атрибут переменной", "строка", "символ"}

	env.SetHLAEnv(high_level_array)
	env.SetEAEnv(expr_array)
	env.SetBGRAEnv()
}

func main() {
	var file_in string
	flag.StringVar(&file_in, "file_in", "", "input rls file")

	flag.Parse()

	rs.Init_Unique_Value()

	env := NewEnv()

	if true {
	}
	Make_Rules(env)

	r_env := rs.NewEnvironment()
	rp := RuleParser{}
	rp.Env = r_env
	env.Struct = rp

	if false {
		//		mgr, gr, hla, ea := CreateRule()
		//		env.SetEnv(mgr, gr, hla, ea)
	}
	if false {
		err := SaveGrammaticRule(env, "rules_n.json")
		if err != nil {
			fmt.Printf("%v\r\n", err)
			return
		}
	}

	//in_file_name := "test1.txt"
	out_file_name := "test2.txt"

	res, err := env.ParseFile(file_in, out_file_name)
	//fmt.Printf("res %v\r\n", res)
	if err != nil {
		fmt.Printf("%v %v\r\n", res, err)
		return
	}
	if false {
		for i, _ := range r_env.Rules {
			r := r_env.Rules[i]
			s := rs.PrintRule(*r)
			fmt.Printf("%v\r\n", s)
		}
	}
	if true {
		for {
			if len(r_env.Rules) == 0 {
				break
			}
			flag := false
			n := len(r_env.Rules)
			for i, _ := range r_env.Rules {
				r := r_env.Rules[i]
				if false {
				        s := rs.PrintRule(*r)
				        fmt.Printf("%v\r\n", s)
				}
				res, err := rs.ExecuteRule(r, r_env)
				if err != nil {
					flag = true
					break
				}
				if res {
					n = n - 1
				} else {
					flag = true
					break
				}
			}
			if flag {
				break
			}
			if n == 0 {
				break
			}
		}
	}
}
