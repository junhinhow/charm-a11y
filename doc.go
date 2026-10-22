// Package a11y detecta as capacidades de acessibilidade do terminal e adapta
// a saída de aplicações Charm a perfis degradados.
//
// O ponto de entrada é [Detect] (ou [Default], que detecta uma única vez por
// processo). O [Profile] resultante pode ser usado para remover cores, trocar
// glifos Unicode por equivalentes ASCII, desligar animações, simplificar
// bordas e produzir texto linear para leitores de tela.
//
// A detecção é conservadora e baseada apenas em variáveis de ambiente, sem
// consultar o terminal. Isso mantém o comportamento previsível e testável.
package a11y
