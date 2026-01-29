package gerador_codigo

import "strconv"

func (gerador *Gerador) GeradorNumero(valor string) {
	gerador.emissor.Emitir("CRCT " + valor)
}

func (gerador *Gerador) GerarVariavel(endereco int) {
	gerador.emissor.Emitir("CRVL " + strconv.Itoa(endereco))
}

func (gerador *Gerador) GerarSoma() {
	gerador.emissor.Emitir("SOMA")
}

func (gerador *Gerador) GerarSubtracao() {
	gerador.emissor.Emitir("SUBT")
}

func (gerador *Gerador) GerarMultiplicacao() {
	gerador.emissor.Emitir("MULT")
}

func (gerador *Gerador) GerarDivisao() {
	gerador.emissor.Emitir("DIVI")
}
