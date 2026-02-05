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

func (gerador *Gerador) LinhaAtual() int {
	return gerador.emissor.LinhaAtual()
}

func (gerador *Gerador) EmitirComMarcador(prefixo string) int {
	return gerador.emissor.EmitirComMarcador(prefixo)
}

func (gerador *Gerador) Preencher(indice int, linhaAlvo int, prefixo string) {
	gerador.emissor.Preencher(indice, linhaAlvo, prefixo)
}

func (gerador *Gerador) StartBuffer() {
	gerador.emissor.StartBuffer()
}

func (gerador *Gerador) GetBufferAndClear() []string {
	return gerador.emissor.GetBufferAndClear()
}

func (gerador *Gerador) EmitBuffer(instrucoes []string) {
	gerador.emissor.EmitBuffer(instrucoes)
}

func (gerador *Gerador) RemoverUltima() string {
	return gerador.emissor.RemoverUltima()
}

func (gerador *Gerador) ReemitirInstrucao(instrucao string) {
	gerador.emissor.Emitir(instrucao)
}

func (gerador *Gerador) InicioPrograma() {
	gerador.emissor.Emitir("INPP")
}

func (gerador *Gerador) FimPrograma() {
	gerador.emissor.Emitir("PARA")
	gerador.emissor.Fechar()
}
