package sintatico

import (
	"fmt"

	"ufmt.br/luiz-crema/compilador/lexer"
)

// Expressao: <expressao> -> <termo> <outros_termos> | input() | read()
func (parser *Parser) Expressao() error {
	if parser.proximoToken.Tag == lexer.TOKEN_INPUT {
		return parser.Input()
	}
	if parser.proximoToken.Tag == lexer.TOKEN_READ {
		return parser.Read()
	}
	parser.gerador.StartBuffer()
	if err := parser.Termo(); err != nil {
		return err
	}
	buf := parser.gerador.GetBufferAndClear()
	primeiraOp := true
	for parser.proximoToken.Tag == lexer.TOKEN_MAIS || parser.proximoToken.Tag == lexer.TOKEN_MENOS {
		op := parser.proximoToken.Tag
		parser.Avancar()

		linhaAntes := parser.gerador.LinhaAtual()
		if err := parser.Termo(); err != nil {
			return err
		}
		linhaDepois := parser.gerador.LinhaAtual()
		instrucoesGeradas := linhaDepois - linhaAntes

		if primeiraOp && len(buf) > 0 {
			if instrucoesGeradas == 1 {
				ultimaInstrucao := parser.gerador.RemoverUltima()
				parser.gerador.EmitBuffer(buf)
				parser.gerador.ReemitirInstrucao(ultimaInstrucao)
			} else {
				parser.gerador.EmitBuffer(buf)
			}
			buf = nil
			primeiraOp = false
		}

		if op == lexer.TOKEN_MAIS {
			parser.gerador.GerarSoma()
		} else {
			parser.gerador.GerarSubtracao()
		}
	}
	if len(buf) > 0 {
		parser.gerador.EmitBuffer(buf)
	}
	return nil
}

func (parser *Parser) Termo() error {
	if err := parser.Fator(); err != nil {
		return err
	}
	for parser.proximoToken.Tag == lexer.TOKEN_VEZES || parser.proximoToken.Tag == lexer.TOKEN_DIVIDIDO {
		op := parser.proximoToken.Tag
		parser.Avancar()
		if err := parser.Fator(); err != nil {
			return err
		}
		if op == lexer.TOKEN_VEZES {
			parser.gerador.GerarMultiplicacao()
		} else {
			parser.gerador.GerarDivisao()
		}
	}
	return nil
}

func (parser *Parser) Fator() error {
	switch parser.proximoToken.Tag {
	case lexer.TOKEN_NUMERO:
		valor := parser.proximoToken.Lexeme
		if err := parser.Confirmar(lexer.TOKEN_NUMERO); err != nil {
			return err
		}
		parser.gerador.GeradorNumero(valor)
		return nil
	case lexer.TOKEN_VARIAVEL:
		nome := parser.proximoToken.Lexeme
		sim, _ := parser.tabelaDeSimbolos.Buscar(nome)
		if err := parser.Confirmar(lexer.TOKEN_VARIAVEL); err != nil {
			return err
		}
		parser.gerador.GerarVariavel(sim.Endereco)
		return nil
	case lexer.TOKEN_ABRE_PARENTESES:
		parser.Avancar()
		if err := parser.Expressao(); err != nil {
			return err
		}
		return parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES)
	default:
		return fmt.Errorf("erro sintático: fator inválido %v (%q)", parser.proximoToken.Tag, parser.proximoToken.Lexeme)
	}
}
