package lexer

import (
	"bufio"
	"fmt"
	"os"
)


type Lexer struct {
	reader *bufio.Reader
	caracter rune
	noInicioLinha bool
}



func NovoLexer(nomeArquivo string) (*Lexer, error) {
	file, err := os.Open(nomeArquivo)
	if err != nil {
		return nil, err
	}

	lexer := &Lexer{
		reader: bufio.NewReader(file),
		noInicioLinha: true,
	}

	lexer.lerCaracter()
	return lexer, nil
}

func (lexer *Lexer) lerCaracter() {
	caracter, _, err := lexer.reader.ReadRune()
	if err != nil {
		lexer.caracter = 0
		return
	}
	lexer.caracter = caracter
}

func (lexer *Lexer) lerNumero() Token {
	var lexeme string

	for ehNumero(lexer.caracter) || lexer.caracter == '.' {
		lexeme += string(lexer.caracter)
		lexer.lerCaracter()
	}

	return Token{Tag: TOKEN_NUMERO, Lexeme: lexeme}
}

func (lexer *Lexer) pularEspacos() {
	for lexer.caracter == ' ' || lexer.caracter == '\t' {
		lexer.lerCaracter()
	}
}

func (lexer *Lexer) indentificarIndentacao() Token {
	var indentacao string
	for lexer.caracter == ' ' || lexer.caracter == '\t' {
		indentacao += string(lexer.caracter)
		lexer.lerCaracter()
	}
	if len(indentacao) > 0 {
		return Token{Tag: TOKEN_IDENTACAO, Lexeme: indentacao}
	}
	return Token{Tag: TOKEN_ERROR, Lexeme: ""}
}

func (lexer *Lexer) pularComentarios() {
	if lexer.caracter == '"' {
		caracter2, err := lexer.reader.Peek(1)
		if err == nil && len(caracter2) > 0 && caracter2[0] == '"' {
			caracter3, err := lexer.reader.Peek(2)
			if err == nil && len(caracter3) > 1 && caracter3[1] == '"' {
				lexer.lerCaracter()
				lexer.lerCaracter()
				lexer.lerCaracter()
				
				for {
					if lexer.caracter == 0 {
						return 
					}
					
					if lexer.caracter == '"' {
						caracter2, err := lexer.reader.Peek(1)
						if err == nil && len(caracter2) > 0 && caracter2[0] == '"' {
							caracter3, err := lexer.reader.Peek(2)
							if err == nil && len(caracter3) > 1 && caracter3[1] == '"' {
								lexer.lerCaracter()
								lexer.lerCaracter()
								lexer.lerCaracter()
								return
							}
						}
					}
					
					lexer.lerCaracter()
				}
			}
		}
	}
}

func (lexer *Lexer) identificarPalavraReservada() Token {
	var lexeme string

	for ehLetra(lexer.caracter) || ehNumero(lexer.caracter) {
		lexeme += string(lexer.caracter)
		lexer.lerCaracter()
	}

	switch lexeme {
	case "def":
		return Token{Tag: TOKEN_DEF, Lexeme: lexeme}
	case "if":
		return Token{Tag: TOKEN_IF, Lexeme: lexeme}
	case "else":
		return Token{Tag: TOKEN_ELSE, Lexeme: lexeme}
	case "while":
		return Token{Tag: TOKEN_WHILE, Lexeme: lexeme}
	case "print":
		return Token{Tag: TOKEN_PRINT, Lexeme: lexeme}
	case "input":
		return Token{Tag: TOKEN_INPUT, Lexeme: lexeme}
	case "read":
		return Token{Tag: TOKEN_READ, Lexeme: lexeme} 
	default:
		return Token{Tag: TOKEN_VARIAVEL, Lexeme: lexeme}
	}

}



func (lexer *Lexer) identificarOperador() Token {
	caracterAtual := lexer.caracter
	lexer.lerCaracter()
	
	switch caracterAtual {
	case '+':
		return Token{Tag: TOKEN_MAIS, Lexeme: "+"}
	case '-':
		return Token{Tag: TOKEN_MENOS, Lexeme: "-"}
	case '*':
		return Token{Tag: TOKEN_VEZES, Lexeme: "*"}
	case '/':
		return Token{Tag: TOKEN_DIVIDIDO, Lexeme: "/"}
	case '=':
		if lexer.caracter == '=' {
			lexer.lerCaracter()
			return Token{Tag: TOKEN_IGUAL, Lexeme: "=="}
		}
		return Token{Tag: TOKEN_ATRIBUICAO, Lexeme: "="}
	case '!':
		if lexer.caracter == '=' {
			lexer.lerCaracter()
			return Token{Tag: TOKEN_DIFERENTE, Lexeme: "!="}
		}
		return Token{Tag: TOKEN_ERROR, Lexeme: "!"}
	case '>':
		if lexer.caracter == '=' {
			lexer.lerCaracter()
			return Token{Tag: TOKEN_MAIOR_IGUAL_QUE, Lexeme: ">="}
		}
		return Token{Tag: TOKEN_MAIOR_QUE, Lexeme: ">"}
	case '<':
		if lexer.caracter == '=' {
			lexer.lerCaracter()
			return Token{Tag: TOKEN_MENOR_IGUAL_QUE, Lexeme: "<="}
		}
		return Token{Tag: TOKEN_MENOR_QUE, Lexeme: "<"}
	default:
		return Token{Tag: TOKEN_ERROR, Lexeme: string(caracterAtual)}
	}
}

func (lexer *Lexer) indentificarDelimitadores() Token {
	caracterAtual := lexer.caracter
	lexer.lerCaracter()
	
	switch caracterAtual {
	case '(':
		return Token{Tag: TOKEN_ABRE_PARENTESES, Lexeme: "("}
	case ')':
		return Token{Tag: TOKEN_FECHA_PARENTESES, Lexeme: ")"}
	case ':':
		return Token{Tag: TOKEN_DOIS_PONTOS, Lexeme: ":"}
	default:
		return Token{Tag: TOKEN_ERROR, Lexeme: string(caracterAtual)}
	}
}


func (lexer *Lexer) ProximoToken() Token {
	for {
		if lexer.noInicioLinha {
			lexer.pularComentarios()
			if lexer.caracter == ' ' || lexer.caracter == '\t' {
				tokenIndentacao := lexer.indentificarIndentacao()
				if tokenIndentacao.Tag != TOKEN_ERROR {
					lexer.pularComentarios()
					return tokenIndentacao
				}
			}
			lexer.pularComentarios()
			lexer.noInicioLinha = false
		}

		lexer.pularEspacos()
		lexer.pularComentarios()
		lexer.pularEspacos()
		
		if lexer.caracter == 0 {
			return Token{Tag: TOKEN_FIM, Lexeme: "PARA"}
		}
		
		if lexer.caracter == '\n' {
			lexer.lerCaracter()
			lexer.noInicioLinha = true
			return Token{Tag: TOKEN_FIM_LINHA, Lexeme: "\n"}
		}
		
		if lexer.caracter == '\r' {
			lexer.lerCaracter()
			continue
		}
		
		if ehLetra(lexer.caracter) {
			return lexer.identificarPalavraReservada()
		}

		if ehNumero(lexer.caracter) {
			return lexer.lerNumero()
		}

		switch lexer.caracter {
		case '+', '-', '*', '/', '=', '!', '>', '<':
			return lexer.identificarOperador()
		case ',':
			lexer.lerCaracter()
			return Token{Tag: TOKEN_VIRGULA, Lexeme: ","} 
		case '(', ')', ':':
			return lexer.indentificarDelimitadores()
		case '"':
			caracter2, err := lexer.reader.Peek(1)
			if err == nil && len(caracter2) > 0 && caracter2[0] == '"' {
				caracter3, err := lexer.reader.Peek(2)
				if err == nil && len(caracter3) > 1 && caracter3[1] == '"' {
					lexer.pularComentarios()
					lexer.pularEspacos()
					continue
				}
			}
			return lexer.indentificarDelimitadores()
		default:
			caracterDesconhecido := lexer.caracter
			lexer.lerCaracter()
			return Token{Tag: TOKEN_ERROR, Lexeme: string(caracterDesconhecido)}
		}
	}
}

func GerarTabelaLexemas(nomeArquivo string) error {
	lexer, err := NovoLexer(nomeArquivo)
	if err != nil {
		return fmt.Errorf("erro ao criar lexer: %v", err)
	}

	var tokens []Token
	for {
		token := lexer.ProximoToken()
		tokens = append(tokens, token)
		if token.Tag == TOKEN_FIM {
			break
		}
	}
	arquivo, err := os.Create("lexemas.txt")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo lexemas.txt: %v", err)
	}
	defer arquivo.Close()

	for _, token := range tokens {
		arquivo.WriteString(fmt.Sprintf("<%s, \"%s\">\n", token.Tag.String(), token.Lexeme))
	}

	return nil
}