package semantico

import "fmt"

type TipoSimbolo int

const (
	SIMBOLO_VARIAVEL TipoSimbolo = iota
	SIMBOLO_FUNCAO
)

type Simbolo struct {
	Nome           string
	Tipo           TipoSimbolo
	QtdeParametros int
	Endereco       int
}

type TabelaDeSimbolos struct {
	simbolos             map[string]Simbolo
	proximoEndereco      int
	adicionadosNaFuncao  []string
	enderecoBaseGlobal   int
	simbolosSobrescritos map[string]Simbolo
}

func NovaTabelaDeSimbolos() *TabelaDeSimbolos {
	return &TabelaDeSimbolos{
		simbolos:            make(map[string]Simbolo),
		proximoEndereco:     0,
		adicionadosNaFuncao: nil,
	}
}

func (tabelaDeSimbolos *TabelaDeSimbolos) DeclararVariavel(nome string) {
	endereco := tabelaDeSimbolos.proximoEndereco
	simboloAntigo, jaExistia := tabelaDeSimbolos.simbolos[nome]

	if jaExistia && tabelaDeSimbolos.simbolosSobrescritos != nil &&
		simboloAntigo.Endereco < tabelaDeSimbolos.enderecoBaseGlobal {
		tabelaDeSimbolos.simbolosSobrescritos[nome] = simboloAntigo
	}

	tabelaDeSimbolos.simbolos[nome] = Simbolo{
		Nome:     nome,
		Tipo:     SIMBOLO_VARIAVEL,
		Endereco: endereco,
	}
	tabelaDeSimbolos.proximoEndereco++

	if tabelaDeSimbolos.adicionadosNaFuncao != nil && !jaExistia {
		tabelaDeSimbolos.adicionadosNaFuncao = append(tabelaDeSimbolos.adicionadosNaFuncao, nome)
	}
}

func (tabelaDeSimbolos *TabelaDeSimbolos) DeclararFuncao(nome string, numParametros int) error {
	if _, existe := tabelaDeSimbolos.simbolos[nome]; existe {
		return fmt.Errorf("erro semântico: função '%s' já declarada", nome)
	}
	tabelaDeSimbolos.simbolos[nome] = Simbolo{
		Nome:           nome,
		Tipo:           SIMBOLO_FUNCAO,
		QtdeParametros: numParametros,
		Endereco:       0,
	}
	return nil
}

func (tabelaDeSimbolos *TabelaDeSimbolos) AtualizarEnderecoFuncao(nome string, endereco int) error {
	simbolo, encontrado := tabelaDeSimbolos.simbolos[nome]
	if !encontrado || simbolo.Tipo != SIMBOLO_FUNCAO {
		return fmt.Errorf("erro semântico: função '%s' não encontrada para atualizar endereço", nome)
	}
	simbolo.Endereco = endereco
	tabelaDeSimbolos.simbolos[nome] = simbolo
	return nil
}

func (tabelaDeSimbolos *TabelaDeSimbolos) ProximoEndereco() int {
	return tabelaDeSimbolos.proximoEndereco
}

func (tabelaDeSimbolos *TabelaDeSimbolos) SetProximoEndereco(endereco int) {
	tabelaDeSimbolos.proximoEndereco = endereco
}

func (tabelaDeSimbolos *TabelaDeSimbolos) GetProximoEndereco() int {
	return tabelaDeSimbolos.proximoEndereco
}

func (tabelaDeSimbolos *TabelaDeSimbolos) DeclararParametro(nome string, endereco int) {
	simboloAntigo, jaExistia := tabelaDeSimbolos.simbolos[nome]

	if jaExistia && tabelaDeSimbolos.simbolosSobrescritos != nil {
		if simboloAntigo.Endereco < tabelaDeSimbolos.enderecoBaseGlobal {
			tabelaDeSimbolos.simbolosSobrescritos[nome] = simboloAntigo
		}
	}

	tabelaDeSimbolos.simbolos[nome] = Simbolo{
		Nome:     nome,
		Tipo:     SIMBOLO_VARIAVEL,
		Endereco: endereco,
	}

	if tabelaDeSimbolos.adicionadosNaFuncao != nil && !jaExistia {
		tabelaDeSimbolos.adicionadosNaFuncao = append(tabelaDeSimbolos.adicionadosNaFuncao, nome)
	}
}

func (tabelaDeSimbolos *TabelaDeSimbolos) EntrarFuncao() {
	tabelaDeSimbolos.enderecoBaseGlobal = tabelaDeSimbolos.proximoEndereco
	tabelaDeSimbolos.adicionadosNaFuncao = make([]string, 0)
	tabelaDeSimbolos.simbolosSobrescritos = make(map[string]Simbolo)
}

func (tabelaDeSimbolos *TabelaDeSimbolos) SairFuncao() {
	for _, nome := range tabelaDeSimbolos.adicionadosNaFuncao {
		if simbolo, ok := tabelaDeSimbolos.simbolos[nome]; ok {
			if simbolo.Endereco >= tabelaDeSimbolos.enderecoBaseGlobal || simbolo.Endereco < 0 {
				delete(tabelaDeSimbolos.simbolos, nome)
			}
		}
	}

	for nome, simboloAntigo := range tabelaDeSimbolos.simbolosSobrescritos {
		tabelaDeSimbolos.simbolos[nome] = simboloAntigo
	}

	tabelaDeSimbolos.adicionadosNaFuncao = nil
	tabelaDeSimbolos.simbolosSobrescritos = nil
}

func (tabelaDeSimbolos *TabelaDeSimbolos) Existe(nome string) bool {
	_, ok := tabelaDeSimbolos.simbolos[nome]
	return ok
}

func (tabelaDeSimbolos *TabelaDeSimbolos) ExisteNaFuncaoAtual(nome string) bool {
	for _, n := range tabelaDeSimbolos.adicionadosNaFuncao {
		if n == nome {
			return true
		}
	}
	if tabelaDeSimbolos.adicionadosNaFuncao != nil {
		if simbolo, existe := tabelaDeSimbolos.simbolos[nome]; existe {
			if simbolo.Endereco < 0 || simbolo.Endereco >= tabelaDeSimbolos.enderecoBaseGlobal {
				return true
			}
		}
	}

	return false
}

func (tabelaDeSimbolos *TabelaDeSimbolos) Adicionar(nome string, tipo TipoSimbolo) error {
	if _, existe := tabelaDeSimbolos.simbolos[nome]; existe {
		return fmt.Errorf("erro semântico: '%s' já declarado", nome)
	}
	tabelaDeSimbolos.simbolos[nome] = Simbolo{Nome: nome, Tipo: tipo}
	return nil
}

func (tabelaDeSimbolos *TabelaDeSimbolos) Listar() {
	for nome, simbolo := range tabelaDeSimbolos.simbolos {
		fmt.Printf("Nome: %s, Tipo: %v\n", nome, simbolo.Tipo)
	}
}

func (tabelaDeSimbolos *TabelaDeSimbolos) Buscar(nome string) (Simbolo, bool) {
	simbolo, encontrado := tabelaDeSimbolos.simbolos[nome]
	return simbolo, encontrado
}

func (tabelaDeSimbolos *TabelaDeSimbolos) VerificarVariavelDeclarada(nome string) error {
	if _, existe := tabelaDeSimbolos.Buscar(nome); !existe {
		return fmt.Errorf("erro semântico: variável '%s' não declarada", nome)
	}
	return nil
}
