package gerador_codigo

import "strconv"

func (gerador *Gerador) GerarPusher(endereco int) {
	gerador.emissor.Emitir("PUSHER " + strconv.Itoa(endereco))
}

func (gerador *Gerador) GerarParametro(endereco int) {
	gerador.emissor.Emitir("PARAM " + strconv.Itoa(endereco))
}

func (gerador *Gerador) GerarChamada(endereco int) {
	gerador.emissor.Emitir("CHPR " + strconv.Itoa(endereco))
}

func (gerador *Gerador) GerarRetorno(quantidade int) {
	gerador.GerarDesalocacao(quantidade)
	gerador.emissor.Emitir("RTPR")
}
