package gerador_codigo

type Gerador struct {
	emissor *Emissor
}

func NovoGerador(saida string) *Gerador {
	return &Gerador{
		emissor: NovoEmissor(saida),
	}
}

func (gerador *Gerador) ProximaLinha() int {
	return gerador.emissor.ProximaLinha()
}

func (gerador *Gerador) EmitirComMarcador(prefixo string) int {
	return gerador.emissor.EmitirComMarcador(prefixo)
}

func (gerador *Gerador) Preencher(indice int, linhaAlvo int, prefixo string) {
	gerador.emissor.Preencher(indice, linhaAlvo, prefixo)
}

func (gerador *Gerador) InicioPrograma() {
	gerador.emissor.Emitir("INPP")
}

func (gerador *Gerador) FimPrograma() {
	gerador.emissor.Emitir("PARA")
	gerador.emissor.Fechar()
}
