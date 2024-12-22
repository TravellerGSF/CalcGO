package calculation

import (
	"strconv"
	"unicode"
)

// Приоритет операторов для алгоритма преобразования в обратную польскую нотацию (RPN)
var priority = map[rune]int{
	'+': 1, // Операторы сложения и вычитания имеют приоритет 1
	'-': 1,
	'*': 2, // Операторы умножения и деления имеют приоритет 2
	'/': 2,
	'(': 0, // Открывающая скобка имеет наименьший приоритет (чтобы она не мешала вычислениям)
}

// Основная функция для вычисления математического выражения
func Calc(expression string) (float64, error) {
	// Преобразуем входную строку в обратную польскую нотацию (RPN)
	rpn, err := convertToRPN(expression)

	// Если при преобразовании возникла ошибка, возвращаем ошибку
	if err != nil {
		return 0, err
	}

	// Вычисляем результат из полученной RPN
	return calculateRPN(rpn)
}

// Функция преобразования математического выражения в обратную польскую нотацию
// RPN (Reverse Polish Notation) - используется для упрощения вычислений без скобок
func convertToRPN(expression string) ([]string, error) {
	var rpn []string     // Массив для хранения итогового выражения в RPN
	var operators []rune // Стек для хранения операторов

	// Функция для добавления оператора в стек с учётом приоритета
	pushOperator := func(op rune) {
		// Пока в стеке есть операторы с приоритетом больше или равным текущему,
		// переносим их в RPN
		for len(operators) > 0 && priority[operators[len(operators)-1]] >= priority[op] {
			rpn = append(rpn, string(operators[len(operators)-1]))
			operators = operators[:len(operators)-1]
		}
		// Добавляем текущий оператор в стек
		operators = append(operators, op)
	}

	i := 0
	for i < len(expression) {
		char := rune(expression[i]) // Преобразуем символ в rune

		// Если символ - это цифра или точка (для чисел с плавающей запятой)
		if unicode.IsDigit(char) || char == '.' {
			j := i
			// Ищем все символы, составляющие число (цифры и точка)
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || rune(expression[i]) == '.') {
				i++
			}
			// Добавляем найденное число в RPN
			rpn = append(rpn, expression[j:i])
			continue
		}

		// Обработка операторов
		switch char {
		case '+', '-', '/', '*': // Для операторов сложения, вычитания, деления и умножения
			pushOperator(char)
		case '(': // Открывающая скобка просто добавляется в стек
			operators = append(operators, char)
		case ')': // Закрывающая скобка
			// Пока не встретится открывающая скобка, переносим операторы в RPN
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				rpn = append(rpn, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			// Если открывающая скобка не найдена, возвращаем ошибку
			if len(operators) == 0 {
				return nil, ErrBrackets
			}
			// Убираем открывающую скобку из стека
			operators = operators[:len(operators)-1]
		default:
			// Если символ не является оператором, числом или пробелом - ошибка
			if !unicode.IsSpace(char) {
				return nil, ErrAllowed
			}
		}
		i++ // Переходим к следующему символу
	}

	// Переносим оставшиеся операторы в RPN
	for len(operators) > 0 {
		// Если осталась открывающая скобка, то это ошибка
		if operators[len(operators)-1] == '(' {
			return nil, ErrBrackets
		}
		// Переносим операторы в RPN
		rpn = append(rpn, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	// Возвращаем полученную обратную польскую нотацию
	return rpn, nil
}

// Функция для вычисления значения выражения, записанного в обратной польской нотации
func calculateRPN(rpn []string) (float64, error) {
	var stack []float64 // Стек для промежуточных результатов

	// Проходим по каждому элементу из выражения в RPN
	for _, elem := range rpn {
		// Если это оператор, выполняем операцию
		switch elem {
		case "+", "-", "*", "/":
			// Проверяем, что в стеке есть хотя бы два числа
			if len(stack) < 2 {
				return 0, ErrValues
			}
			// Берём два числа из стека
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2] // Убираем эти два числа из стека

			var result float64
			// Выполняем операцию в зависимости от оператора
			switch elem {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				// Проверка на деление на ноль
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				result = a / b
			}
			// Результат операции добавляем обратно в стек
			stack = append(stack, result)
		default:
			// Если элемент - число, парсим его и кладём в стек
			value, err := strconv.ParseFloat(elem, 64)
			if err != nil {
				return 0, ErrAllowed
			}
			stack = append(stack, value)
		}
	}

	// В стеке должно остаться только одно значение - результат
	if len(stack) != 1 {
		return 0, ErrValues
	}

	// Возвращаем результат
	return stack[0], nil
}
