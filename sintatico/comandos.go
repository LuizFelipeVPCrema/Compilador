package sintatico

import (
	"fmt"

	"ufmt.br/luiz-crema/compilador/lexer"
)

// Comandos: <comandos> -> <comando> <mais_comandos>
func (parser *Parser) Comandos() error {
	if parser.indentacaoPendente != "" && !parser.exigeIndentacao {
		return nil
	}

	for parser.proximoToken.Tag == lexer.TOKEN_FIM_LINHA {
		parser.Avancar()
	}
	if parser.proximoToken.Tag == lexer.TOKEN_FIM {
		return nil
	}
	if parser.proximoToken.Tag == lexer.TOKEN_INICIO_LINHA_SEM_INDENT {
		parser.Avancar()
	}

	if err := parser.Comando(); err != nil {
		return err
	}
	return parser.MaisComandos()
}

// MaisComandos: <mais_comandos> -> <comandos> | λ
func (parser *Parser) MaisComandos() error {
	for parser.proximoToken.Tag == lexer.TOKEN_FIM_LINHA {
		parser.Avancar()
	}

	var indent string
	if parser.indentacaoPendente != "" {
		indent = parser.indentacaoPendente
		parser.indentacaoPendente = ""
	} else if parser.proximoToken.Tag == lexer.TOKEN_IDENTACAO {

		for parser.proximoToken.Tag == lexer.TOKEN_IDENTACAO {
			indent += parser.proximoToken.Lexeme
			parser.Avancar()
		}

		if parser.exigeIndentacao && len(indent) < len(parser.indentacaoBloco) {
			parser.indentacaoPendente = indent
			return nil
		}
	} else if parser.exigeIndentacao {
		return nil
	}

	if parser.proximoToken.Tag == lexer.TOKEN_INICIO_LINHA_SEM_INDENT {
		parser.Avancar()
	}

	if indent != "" {
		if parser.proximoToken.Tag == lexer.TOKEN_INICIO_LINHA_SEM_INDENT {
			parser.Avancar()
			if parser.proximoToken.Tag == lexer.TOKEN_ELSE {
				return nil
			}
		}
		return parser.Comandos()
	}

	if parser.exigeIndentacao {
		return nil
	}
	switch parser.proximoToken.Tag {
	case lexer.TOKEN_VARIAVEL, lexer.TOKEN_PRINT, lexer.TOKEN_IF, lexer.TOKEN_WHILE:
		return parser.Comandos()
	default:
		return nil
	}
}

func (parser *Parser) Comando() error {
	switch parser.proximoToken.Tag {
	case lexer.TOKEN_VARIAVEL:
		return parser.AtribuicaoOuChamada()
	case lexer.TOKEN_PRINT:
		return parser.Imprimir()
	case lexer.TOKEN_IF:
		return parser.Condicional()
	case lexer.TOKEN_WHILE:
		return parser.While()
	case lexer.TOKEN_ELSE:
		return nil
	default:
		return fmt.Errorf("erro sintático: comando inesperado %v (%q)", parser.proximoToken.Tag, parser.proximoToken.Lexeme)
	}
}

func (parser *Parser) AtribuicaoOuChamada() error {
	nome := parser.proximoToken.Lexeme

	parser.Avancar()

	if parser.proximoToken.Tag == lexer.TOKEN_ABRE_PARENTESES {
		return parser.ChamadaFuncaoComGerador(nome)
	}

	if parser.dentroDeFuncao && !parser.tabelaDeSimbolos.ExisteNaFuncaoAtual(nome) && parser.proximoToken.Tag == lexer.TOKEN_ATRIBUICAO {
		parser.tabelaDeSimbolos.DeclararVariavel(nome)
		parser.gerador.GerarAlocacao(1)
		parser.numLocaisAtual++
	}

	if err := parser.tabelaDeSimbolos.VerificarVariavelDeclarada(nome); err != nil {
		return err
	}

	if parser.proximoToken.Tag == lexer.TOKEN_ATRIBUICAO {
		parser.Avancar()

		if parser.proximoToken.Tag == lexer.TOKEN_NUMERO && (parser.proximoToken.Lexeme == "0" || parser.proximoToken.Lexeme == "0.0") {
			parser.Avancar()
			return nil
		}
		if err := parser.Expressao(); err != nil {
			return err
		}
		sim, _ := parser.tabelaDeSimbolos.Buscar(nome)
		parser.gerador.GerarAtribuicao(sim.Endereco)
		return nil
	}

	return fmt.Errorf("erro sintático: esperado '=' ou '(' após uma variável")
}
