package sintatico

import (
	"fmt"

	"ufmt.br/luiz-crema/compilador/lexer"
	"ufmt.br/luiz-crema/compilador/semantico"
)

// def funcao(a, b):
//
//	bloco_funcao
func (parser *Parser) Funcao() error {
	if err := parser.Confirmar(lexer.TOKEN_DEF); err != nil {
		return err
	}
	nomeFuncao := parser.proximoToken.Lexeme
	if err := parser.Confirmar(lexer.TOKEN_VARIAVEL); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	nomesParams, err := parser.ParametrosNomes()
	if err != nil {
		return err
	}
	numParams := len(nomesParams)
	if err := parser.tabelaDeSimbolos.DeclararFuncao(nomeFuncao, numParams); err != nil {
		return err
	}
	primeiraFuncao := len(parser.indicesDSVI) == 0
	parser.indiceDSVI = parser.gerador.EmitirComMarcador("DSVI")
	linhaEntrada := parser.gerador.ProximaLinha()
	if primeiraFuncao {
		linhaEntrada++
	}
	if err := parser.tabelaDeSimbolos.AtualizarEnderecoFuncao(nomeFuncao, linhaEntrada); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES); err != nil {
		return err
	}
	parser.tabelaDeSimbolos.EntrarFuncao()

	proximoEnderecoAntes := parser.tabelaDeSimbolos.GetProximoEndereco()

	base := parser.enderecoBaseFuncoes
	for i, nome := range nomesParams {
		parser.tabelaDeSimbolos.DeclararParametro(nome, base+i)
	}
	parser.tabelaDeSimbolos.SetProximoEndereco(base + numParams)
	parser.numParamsAtual = numParams
	parser.numLocaisAtual = 0
	parser.nomeFuncaoAtual = nomeFuncao
	if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
		return err
	}
	parser.dentroDeFuncao = true
	erro := parser.Bloco()
	parser.dentroDeFuncao = false
	parser.exigeIndentacao = false
	parser.tabelaDeSimbolos.SairFuncao()
	parser.gerador.GerarRetorno(parser.numParamsAtual + parser.numLocaisAtual)
	parser.tabelaDeSimbolos.SetProximoEndereco(proximoEnderecoAntes)

	parser.indicesDSVI = append(parser.indicesDSVI, parser.indiceDSVI)
	return erro
}

func (parser *Parser) ParametrosNomes() ([]string, error) {
	if parser.proximoToken.Tag != lexer.TOKEN_VARIAVEL {
		return nil, nil
	}
	var nomes []string
	nome := parser.proximoToken.Lexeme
	nomes = append(nomes, nome)
	if err := parser.Confirmar(lexer.TOKEN_VARIAVEL); err != nil {
		return nil, err
	}
	for parser.proximoToken.Tag == lexer.TOKEN_VIRGULA {
		parser.Avancar()
		nome = parser.proximoToken.Lexeme
		nomes = append(nomes, nome)
		if err := parser.Confirmar(lexer.TOKEN_VARIAVEL); err != nil {
			return nil, err
		}
	}
	return nomes, nil
}

// ChamadaFuncaoComGerador: PUSHER <linha_retorno>, PARAM por argumento, CHPR <endereço>.
func (parser *Parser) ChamadaFuncaoComGerador(nome string) error {
	simbolo, encontrado := parser.tabelaDeSimbolos.Buscar(nome)
	if !encontrado || simbolo.Tipo != semantico.SIMBOLO_FUNCAO {
		return fmt.Errorf("erro semântico: '%s' não é uma função declarada", nome)
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	var enderecos []int
	if parser.proximoToken.Tag != lexer.TOKEN_FECHA_PARENTESES {
		endereco, err := parser.ArgumentoEndereco()
		if err != nil {
			return err
		}
		enderecos = append(enderecos, endereco)
		for parser.proximoToken.Tag == lexer.TOKEN_VIRGULA {
			parser.Avancar()
			endereco, err := parser.ArgumentoEndereco()
			if err != nil {
				return err
			}
			enderecos = append(enderecos, endereco)
		}
	}
	if err := parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES); err != nil {
		return err
	}
	if len(enderecos) != simbolo.QtdeParametros {
		return fmt.Errorf("erro semântico: função '%s' espera %d parâmetro(s), recebeu %d", nome, simbolo.QtdeParametros, len(enderecos))
	}
	linhaRetorno := parser.gerador.ProximaLinha() + len(enderecos) + 2
	parser.gerador.GerarPusher(linhaRetorno)
	for _, endereco := range enderecos {
		parser.gerador.GerarParametro(endereco)
	}
	parser.gerador.GerarChamada(simbolo.Endereco)
	return nil
}

func (parser *Parser) ArgumentoEndereco() (int, error) {
	if parser.proximoToken.Tag != lexer.TOKEN_VARIAVEL {
		return 0, fmt.Errorf("erro sintático: esperado identificador como argumento, veio %v", parser.proximoToken.Tag)
	}
	nome := parser.proximoToken.Lexeme
	sim, ok := parser.tabelaDeSimbolos.Buscar(nome)
	if !ok {
		return 0, fmt.Errorf("erro semântico: variável '%s' não declarada", nome)
	}
	parser.Avancar()
	return sim.Endereco, nil
}
