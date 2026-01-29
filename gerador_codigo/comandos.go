package gerador_codigo

import "strconv"

func (gerador *Gerador) GerarAtribuicao(endereco int) {
	gerador.emissor.Emitir("ARMZ " + strconv.Itoa(endereco))
}

func (gerador *Gerador) GerarPrint(endereco int) {
	gerador.emissor.Emitir("CRVL " + strconv.Itoa(endereco))
	gerador.emissor.Emitir("IMPR")
}
