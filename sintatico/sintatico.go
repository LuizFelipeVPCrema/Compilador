package sintatico

import (
	"fmt"

	"ufmt.br/luiz-crema/compilador/gerador_codigo"
	"ufmt.br/luiz-crema/compilador/lexer"
	"ufmt.br/luiz-crema/compilador/semantico"
)

type Parser struct {
	lexer              *lexer.Lexer
	proximoToken       lexer.Token
	tabelaDeSimbolos   *semantico.TabelaDeSimbolos
	gerador            *gerador_codigo.Gerador
	dentroDeFuncao     bool
	exigeIndentacao    bool
	indentacaoBloco    string
	indentacaoPendente string

	numParamsAtual  int
	numLocaisAtual  int
	indiceDSVI      int
	nomeFuncaoAtual string
	indicesDSVI     []int

	enderecoBaseFuncoes int
}

func NovoParser(lexer *lexer.Lexer, tabela *semantico.TabelaDeSimbolos, gerador *gerador_codigo.Gerador) *Parser {
	parser := &Parser{
		lexer:               lexer,
		tabelaDeSimbolos:    tabela,
		gerador:             gerador,
		enderecoBaseFuncoes: -1, // -1 = ainda não definido (será definido na primeira função)
	}
	parser.proximoToken = lexer.ProximoToken()
	return parser
}

func (parser *Parser) Avancar() {
	parser.proximoToken = parser.lexer.ProximoToken()
}

func (parser *Parser) Confirmar(tag lexer.Tag) error {
	if parser.proximoToken.Tag != tag {
		return fmt.Errorf("erro sintático: esperado %v, veio %v (%q)", tag, parser.proximoToken.Tag, parser.proximoToken.Lexeme)
	}
	parser.Avancar()
	return nil
}

// Programa: início, corpo (declarações + comandos) e fim.
func (parser *Parser) Programa() error {
	parser.gerador.InicioPrograma()
	if err := parser.Corpo(); err != nil {
		return err
	}
	parser.gerador.FimPrograma()
	return nil
}

// Corpo: <corpo> -> <dc> <comandos>
func (parser *Parser) Corpo() error {
	if err := parser.DC(); err != nil {
		return err
	}
	linhaInicioMain := parser.gerador.LinhaAtual()
	for _, idx := range parser.indicesDSVI {
		parser.gerador.Preencher(idx, linhaInicioMain, "DSVI")
	}
	parser.indicesDSVI = nil
	return parser.Comandos()
}

// DC: <dc> -> <dc_v> <mais_dc> | <dc_f> | λ
func (parser *Parser) DC() error {
	for parser.proximoToken.Tag == lexer.TOKEN_FIM_LINHA {
		parser.Avancar()
	}
	if parser.proximoToken.Tag == lexer.TOKEN_INICIO_LINHA_SEM_INDENT {
		parser.Avancar()
	}
	switch parser.proximoToken.Tag {
	case lexer.TOKEN_VARIAVEL:
		_, encontrou := parser.tabelaDeSimbolos.Buscar(parser.proximoToken.Lexeme)
		if encontrou {
			return nil
		}
		return parser.dcVar()
	case lexer.TOKEN_DEF:
		return parser.dcFunc()
	default:
		return nil
	}
}

// dcVar: <dc_v> -> ident = expressao
func (parser *Parser) dcVar() error {
	nome := parser.proximoToken.Lexeme
	parser.Avancar()
	if parser.proximoToken.Tag != lexer.TOKEN_ATRIBUICAO {
		return fmt.Errorf("erro sintático: esperado '=' em declaração de variável, veio %v (%q)", parser.proximoToken.Tag, parser.proximoToken.Lexeme)
	}
	parser.tabelaDeSimbolos.DeclararVariavel(nome)
	parser.gerador.GerarAlocacao(1)
	parser.Confirmar(lexer.TOKEN_ATRIBUICAO)
	if parser.proximoToken.Tag == lexer.TOKEN_NUMERO && (parser.proximoToken.Lexeme == "0" || parser.proximoToken.Lexeme == "0.0") {
		parser.Avancar()
		return parser.MaisDC()
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	return parser.MaisDC()
}

// dcFunc: <dc_f> -> def ident parametros : corpo_f
func (parser *Parser) dcFunc() error {

	if parser.enderecoBaseFuncoes < 0 {
		parser.enderecoBaseFuncoes = parser.tabelaDeSimbolos.ProximoEndereco()
	}
	if err := parser.Funcao(); err != nil {
		return err
	}
	return parser.MaisDC()
}

// MaisDC: <mais_dc> -> <dc> | λ
func (parser *Parser) MaisDC() error {
	for parser.proximoToken.Tag == lexer.TOKEN_FIM_LINHA {
		parser.Avancar()
	}
	return parser.DC()
}

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

// <relacao> -> == | != | >= | <= | > | <
func (parser *Parser) Relacao() error {
	switch parser.proximoToken.Tag {
	case lexer.TOKEN_IGUAL:
		return parser.Confirmar(lexer.TOKEN_IGUAL)
	case lexer.TOKEN_DIFERENTE:
		return parser.Confirmar(lexer.TOKEN_DIFERENTE)
	case lexer.TOKEN_MAIOR_IGUAL_QUE:
		return parser.Confirmar(lexer.TOKEN_MAIOR_IGUAL_QUE)
	case lexer.TOKEN_MENOR_IGUAL_QUE:
		return parser.Confirmar(lexer.TOKEN_MENOR_IGUAL_QUE)
	case lexer.TOKEN_MAIOR_QUE:
		return parser.Confirmar(lexer.TOKEN_MAIOR_QUE)
	case lexer.TOKEN_MENOR_QUE:
		return parser.Confirmar(lexer.TOKEN_MENOR_QUE)
	default:
		return fmt.Errorf("erro sintático: operador de comparação inválido %v", parser.proximoToken.Tag)
	}
}

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

// Condicao: <condicao> -> <expressao> <relacao> <expressao>.
func (parser *Parser) Condicao() error {
	if err := parser.Expressao(); err != nil {
		return err
	}
	if err := parser.Relacao(); err != nil {
		return err
	}
	return parser.Expressao()
}

// emitirCondicao parse a condição e emite o código (expr1, expr2, comparação).
func (parser *Parser) emitirCondicao() error {
	if err := parser.Expressao(); err != nil {
		return err
	}
	relTag := parser.proximoToken.Tag
	if err := parser.Relacao(); err != nil {
		return err
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	parser.emitirComparacao(relTag)
	return nil
}

func (parser *Parser) emitirComparacao(relTag lexer.Tag) {
	switch relTag {
	case lexer.TOKEN_IGUAL:
		parser.gerador.GerarIgual()
	case lexer.TOKEN_DIFERENTE:
		parser.gerador.GerarDiferente()
	case lexer.TOKEN_MAIOR_IGUAL_QUE:
		parser.gerador.GerarMaiorIgual()
	case lexer.TOKEN_MENOR_IGUAL_QUE:
		parser.gerador.GerarMenorIgual()
	case lexer.TOKEN_MAIOR_QUE:
		parser.gerador.GerarMaior()
	case lexer.TOKEN_MENOR_QUE:
		parser.gerador.GerarMenor()
	}
}

// Condicional: condição (emitida), DSVF, bloco então, DSVI, [bloco senão].
func (parser *Parser) Condicional() error {
	if err := parser.Confirmar(lexer.TOKEN_IF); err != nil {
		return err
	}
	return parser.gerador.GerarIf(
		parser.emitirCondicao,
		func() error {
			if parser.proximoToken.Tag == lexer.TOKEN_FECHA_PARENTESES {
				parser.Avancar()
			}
			if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
				return err
			}
			return parser.Bloco()
		},
		func() error {
			if parser.proximoToken.Tag != lexer.TOKEN_ELSE {
				return nil
			}
			parser.Avancar()
			if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
				return err
			}
			return parser.Bloco()
		},
	)
}

// Bloco analisa: <bloco> -> tabulacao <comandos>
func (parser *Parser) Bloco() error {
	antigaIndentacaoBloco := parser.indentacaoBloco
	defer func() {
		parser.indentacaoBloco = antigaIndentacaoBloco
	}()

	parser.exigeIndentacao = true

	parser.indentacaoPendente = ""

	for parser.proximoToken.Tag == lexer.TOKEN_FIM_LINHA {
		parser.Avancar()
	}

	if parser.proximoToken.Tag != lexer.TOKEN_IDENTACAO {
		return nil
	}
	parser.indentacaoBloco = ""
	for parser.proximoToken.Tag == lexer.TOKEN_IDENTACAO {
		parser.indentacaoBloco += parser.proximoToken.Lexeme
		parser.Avancar()
	}
	return parser.Comandos()
}

// While: condição (emitida), DSVF, bloco, DSVI (início).
func (parser *Parser) While() error {
	if err := parser.Confirmar(lexer.TOKEN_WHILE); err != nil {
		return err
	}
	return parser.gerador.GerarWhile(
		parser.emitirCondicao,
		func() error {
			if parser.proximoToken.Tag == lexer.TOKEN_FECHA_PARENTESES {
				parser.Avancar()
			}
			if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
				return err
			}
			return parser.Bloco()
		},
	)
}

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
	parser.indiceDSVI = parser.gerador.EmitirComMarcador("DSVI")
	linhaEntrada := parser.gerador.ProximaLinha()
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
	linhaRetorno := parser.gerador.ProximaLinha() + len(enderecos) + 1
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
