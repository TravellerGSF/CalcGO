package calculation

import "errors"

var (
	ErrBrackets       = errors.New("Неверно. Количество скобок не совпадает")
	ErrValues         = errors.New("Неверно. Недостаточно значений")
	ErrDivisionByZero = errors.New("Неверно. Деление на ноль")
	ErrAllowed        = errors.New("Недопустимо. Допускаются только числа и ( ) + - * /")
)
