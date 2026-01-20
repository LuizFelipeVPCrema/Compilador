package lexer

func ehLetra(caracter rune) bool {
	return (caracter >= 'a' && caracter <= 'z') || (caracter >= 'A' && caracter <= 'Z')
}

func ehNumero(caracter rune) bool {
	return caracter >= '0' && caracter <= '9'
}