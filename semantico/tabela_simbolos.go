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
	simbolos            map[string]Simbolo
	proximoEndereco     int
	adicionadosNaFuncao []string
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
	tabelaDeSimbolos.simbolos[nome] = Simbolo{
		Nome:     nome,
		Tipo:     SIMBOLO_VARIAVEL,
		Endereco: endereco,
	}
	tabelaDeSimbolos.proximoEndereco++
	if tabelaDeSimbolos.adicionadosNaFuncao != nil {
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

func (tabelaDeSimbolos *TabelaDeSimbolos) DeclararParametro(nome string, endereco int) {
	tabelaDeSimbolos.simbolos[nome] = Simbolo{
		Nome:     nome,
		Tipo:     SIMBOLO_VARIAVEL,
		Endereco: endereco,
	}
	if tabelaDeSimbolos.adicionadosNaFuncao != nil {
		tabelaDeSimbolos.adicionadosNaFuncao = append(tabelaDeSimbolos.adicionadosNaFuncao, nome)
	}
}

func (tabelaDeSimbolos *TabelaDeSimbolos) EntrarFuncao() {
	tabelaDeSimbolos.adicionadosNaFuncao = make([]string, 0)
}

func (tabelaDeSimbolos *TabelaDeSimbolos) SairFuncao() {
	for _, nome := range tabelaDeSimbolos.adicionadosNaFuncao {
		delete(tabelaDeSimbolos.simbolos, nome)
	}
	tabelaDeSimbolos.adicionadosNaFuncao = nil
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
