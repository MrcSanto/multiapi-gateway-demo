# Dashboard de Métricas de Desempenho das APIs
Este dashboard interativo em Streamlit foi desenvolvido para visualizar e analisar os resultados dos testes de desempenho realizados nas três APIs (Go, Python e Node.js).

## Por quê?

#### Análise Visual Interativa
  - Exploração dinâmica: Permite visualizar e comparar métricas das três APIs de forma interativa
  - Identificação de padrões: Facilita a detecção de anomalias e tendências nos dados
  - Validação de resultados: Essencial para verificar os dados antes de incluí-los no artigo acadêmico

#### Complemento ao Artigo

  - Geração de gráficos de alta qualidade para o artigo LaTeX
  - Análise exploratória dos dados brutos do JMeter
  - Validação das métricas consolidadas apresentadas nas tabelas

## Como executar

- Instale as dependências
  ```bash
  uv sync
  ```

- Rode o projeto
  ```bash
  streamlit run ./main.py
  ```