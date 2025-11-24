from pathlib import Path

import streamlit as st
import plotly.graph_objects as go
import plotly.express as px
import pandas as pd

PATH_CURR_DIR = Path(__file__).resolve().parent
PATH_DOCS = PATH_CURR_DIR / "docs"


# Importando datasets
summary_GO_1 = pd.read_csv(PATH_DOCS / "go_api_metrics" / "summaryCARGA_1_GO.csv")
results_GO_1 = pd.read_csv(PATH_DOCS / "go_api_metrics" / "resultsCARGA_1_GO.csv")

summary_GO_2 = pd.read_csv(PATH_DOCS / "go_api_metrics" / "summaryCARGA_2_GO.csv")
results_GO_2 = pd.read_csv(PATH_DOCS / "go_api_metrics" / "resultsCARGA_2_GO.csv")

summary_GO_3 = pd.read_csv(PATH_DOCS / "go_api_metrics" / "summaryCARGA_3_GO.csv")
results_GO_3 = pd.read_csv(PATH_DOCS / "go_api_metrics" / "resultsCARGA_3_GO.csv")


summary_NODE_1 = pd.read_csv(PATH_DOCS / "node_api_metrics" / "summaryCARGA_1_NODE.csv")
results_NODE_1 = pd.read_csv(PATH_DOCS / "node_api_metrics" / "resultsCARGA_1_NODE.csv")

summary_NODE_2 = pd.read_csv(PATH_DOCS / "node_api_metrics" / "summaryCARGA_2_NODE.csv")
results_NODE_2 = pd.read_csv(PATH_DOCS / "node_api_metrics" / "resultsCARGA_2_NODE.csv")

summary_NODE_3 = pd.read_csv(PATH_DOCS / "node_api_metrics" / "summaryCARGA_3_NODE.csv")
results_NODE_3 = pd.read_csv(PATH_DOCS / "node_api_metrics" / "resultsCARGA_3_NODE.csv")


summary_Python_1 = pd.read_csv(
    PATH_DOCS / "python_api_metrics" / "summaryCARGA_1_Python.csv"
)
results_Python_1 = pd.read_csv(
    PATH_DOCS / "python_api_metrics" / "resultsCARGA_1_Python.csv"
)

summary_Python_2 = pd.read_csv(
    PATH_DOCS / "python_api_metrics" / "summaryCARGA_2_Python.csv"
)
results_Python_2 = pd.read_csv(
    PATH_DOCS / "python_api_metrics" / "resultsCARGA_2_Python.csv"
)

summary_Python_3 = pd.read_csv(
    PATH_DOCS / "python_api_metrics" / "summaryCARGA_3_Python.csv"
)
results_Python_3 = pd.read_csv(
    PATH_DOCS / "python_api_metrics" / "resultsCARGA_3_Python.csv"
)

# Título da aplicação
st.title("Dashboard de Métricas de Desempenho das APIs")

st.markdown(
    """
### Testes Realizados:
- **Teste 1** → 100 threads, 2 repetições cada → **200 requisições**
- **Teste 2** → 500 threads, 5 repetições cada → **2.500 requisições**
- **Teste 3** → 1.000 threads, 10 repetições cada → **10.000 requisições**
"""
)


for df in [
    summary_GO_1,
    summary_GO_2,
    summary_GO_3,
    summary_NODE_1,
    summary_NODE_2,
    summary_NODE_3,
    summary_Python_1,
    summary_Python_2,
    summary_Python_3,
]:
    df.drop(df[df["Label"] == "TOTAL"].index, inplace=True)

# Gráfico 1: Tempo médio por operação
fig_avg = go.Figure()
fig_avg.add_trace(
    go.Bar(x=summary_GO_1["Label"], y=summary_GO_1["Average"], name="GO - Rodada 1")
)
fig_avg.add_trace(
    go.Bar(x=summary_GO_2["Label"], y=summary_GO_2["Average"], name="GO - Rodada 2")
)
fig_avg.add_trace(
    go.Bar(x=summary_GO_3["Label"], y=summary_GO_3["Average"], name="GO - Rodada 3")
)
fig_avg.add_trace(
    go.Bar(
        x=summary_NODE_1["Label"], y=summary_NODE_1["Average"], name="Node - Rodada 1"
    )
)
fig_avg.add_trace(
    go.Bar(
        x=summary_NODE_2["Label"], y=summary_NODE_2["Average"], name="Node - Rodada 2"
    )
)
fig_avg.add_trace(
    go.Bar(
        x=summary_NODE_3["Label"], y=summary_NODE_3["Average"], name="Node - Rodada 3"
    )
)
fig_avg.add_trace(
    go.Bar(
        x=summary_Python_1["Label"],
        y=summary_Python_1["Average"],
        name="Python - Rodada 1",
    )
)
fig_avg.add_trace(
    go.Bar(
        x=summary_Python_2["Label"],
        y=summary_Python_2["Average"],
        name="Python - Rodada 2",
    )
)
fig_avg.add_trace(
    go.Bar(
        x=summary_Python_3["Label"],
        y=summary_Python_3["Average"],
        name="Python - Rodada 3",
    )
)

fig_avg.update_layout(
    barmode="group",
    title="Comparativo: Tempo Médio por Operação (ms)",
    xaxis_title="Operação",
    yaxis_title="Tempo (ms)",
    template="presentation",
    height=600,
)

st.plotly_chart(fig_avg, use_container_width=True)
st.set_page_config(layout="wide")

# Gráfico 2: Throughput por operação
fig_throughput = go.Figure()
fig_throughput.add_trace(
    go.Scatter(
        x=summary_GO_1["Label"],
        y=summary_GO_1["Throughput"],
        mode="lines+markers",
        name="GO - Rodada 1",
    )
)
fig_throughput.add_trace(
    go.Scatter(
        x=summary_GO_2["Label"],
        y=summary_GO_2["Throughput"],
        mode="lines+markers",
        name="GO - Rodada 2",
    )
)
fig_throughput.add_trace(
    go.Scatter(
        x=summary_GO_3["Label"],
        y=summary_GO_3["Throughput"],
        mode="lines+markers",
        name="GO - Rodada 3",
    )
)
fig_throughput.add_trace(
    go.Scatter(
        x=summary_NODE_1["Label"],
        y=summary_NODE_1["Throughput"],
        mode="lines+markers",
        name="Node - Rodada 1",
    )
)
fig_throughput.add_trace(
    go.Scatter(
        x=summary_NODE_2["Label"],
        y=summary_NODE_2["Throughput"],
        mode="lines+markers",
        name="Node - Rodada 2",
    )
)
fig_throughput.add_trace(
    go.Scatter(
        x=summary_NODE_3["Label"],
        y=summary_NODE_3["Throughput"],
        mode="lines+markers",
        name="Node - Rodada 3",
    )
)
fig_throughput.add_trace(
    go.Scatter(
        x=summary_Python_1["Label"],
        y=summary_Python_1["Throughput"],
        mode="lines+markers",
        name="Python - Rodada 1",
    )
)
fig_throughput.add_trace(
    go.Scatter(
        x=summary_Python_2["Label"],
        y=summary_Python_2["Throughput"],
        mode="lines+markers",
        name="Python - Rodada 2",
    )
)
fig_throughput.add_trace(
    go.Scatter(
        x=summary_Python_3["Label"],
        y=summary_Python_3["Throughput"],
        mode="lines+markers",
        name="Python - Rodada 3",
    )
)

fig_throughput.update_layout(
    title="Comparativo: Throughput por Operação",
    xaxis_title="Operação",
    yaxis_title="Req/s",
    template="presentation",
    height=600,
)

st.plotly_chart(fig_throughput, use_container_width=True)

# Gráfico 3: Percentual de Erros por Operação
fig_error = go.Figure()
for df, api, rodada in [
    (summary_GO_1, "GO", 1),
    (summary_GO_2, "GO", 2),
    (summary_GO_3, "GO", 3),
    (summary_NODE_1, "Node", 1),
    (summary_NODE_2, "Node", 2),
    (summary_NODE_3, "Node", 3),
    (summary_Python_1, "Python", 1),
    (summary_Python_2, "Python", 2),
    (summary_Python_3, "Python", 3),
]:
    fig_error.add_trace(
        go.Bar(
            x=df["Label"],
            y=[float(x.strip("%")) for x in df["Error %"]],
            name=f"{api} - Rodada {rodada}",
        )
    )

fig_error.update_layout(
    barmode="group",
    title="Comparativo: Percentual de Erros por Operação",
    xaxis_title="Operação",
    yaxis_title="Erro (%)",
    template="presentation",
    height=600,
)

st.plotly_chart(fig_error, use_container_width=True)


# Gráfico 4: Latência Média por API e Rodada
def calcular_latencia_media(df_latencia):
    return df_latencia["Latency"].mean()


arquivos = {
    "GO": [results_GO_1, results_GO_2, results_GO_3],
    "Node": [results_NODE_1, results_NODE_2, results_NODE_3],
    "Python": [results_Python_1, results_Python_2, results_Python_3],
}

dados = []
for api, lista_arquivos in arquivos.items():
    for i, arquivo in enumerate(lista_arquivos, start=1):
        latencia = calcular_latencia_media(arquivo)
        dados.append(
            {"API": api, "Rodada": f"Rodada {i}", "Latência Média (ms)": latencia}
        )

df_latencias = pd.DataFrame(dados)

fig_latencia = px.bar(
    df_latencias,
    x="Rodada",
    y="Latência Média (ms)",
    color="API",
    barmode="group",
    title="Latência Média por API",
)
st.plotly_chart(fig_latencia, use_container_width=True)
