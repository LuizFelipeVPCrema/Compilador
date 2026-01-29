package gerador_codigo

import "strconv"

func (gerador *Gerador) GerarIf(
	gerarCondicao func(),
	blocoVerdadeiro func(),
	blocoFalso func(),
) {
	condicaoSenao := gerador.NovoRotulo()
	condicaoFim := gerador.NovoRotulo()

	gerarCondicao()
	gerador.emissor.Emitir("DSVF " + strconv.Itoa(condicaoSenao))

	blocoVerdadeiro()
	gerador.emissor.Emitir("DSVI " + strconv.Itoa(condicaoFim))

	gerador.emissor.Emitir(strconv.Itoa(condicaoSenao))

	if blocoFalso != nil {
		blocoFalso()
	}

	gerador.emissor.Emitir(strconv.Itoa(condicaoFim))
}

func (gerador *Gerador) GerarWhile(
	gerarCondicao func(),
	bloco func(),
) {
	condicaoInicio := gerador.NovoRotulo()
	condicaoFim := gerador.NovoRotulo()

	gerador.emissor.Emitir(strconv.Itoa(condicaoInicio))
	gerarCondicao()
	gerador.emissor.Emitir("DSVF " + strconv.Itoa(condicaoFim))

	bloco()
	gerador.emissor.Emitir("DSVI " + strconv.Itoa(condicaoInicio))
	gerador.emissor.Emitir(strconv.Itoa(condicaoFim))
}
