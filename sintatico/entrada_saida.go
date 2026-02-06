package sintatico

import (
	"ufmt.br/luiz-crema/compilador/lexer"
)

// print(expressao).
func (parser *Parser) Imprimir() error {
	if err := parser.Confirmar(lexer.TOKEN_PRINT); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	parser.gerador.GerarImprimirPilha()
	return parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES)
}

// read() emite LEIT.
func (parser *Parser) Read() error {
	if err := parser.Confirmar(lexer.TOKEN_READ); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES); err != nil {
		return err
	}
	parser.gerador.GerarLeitura()
	return nil
}

// input() emite LEIT.
func (parser *Parser) Input() error {
	if err := parser.Confirmar(lexer.TOKEN_INPUT); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES); err != nil {
		return err
	}
	parser.gerador.GerarLeitura()
	return nil
}
