# 📊 Log Cluster Counter com go-drain3

Este projeto em Go processa um arquivo de log, agrupa mensagens de log semelhantes em **clusters de templates** usando o [go-drain3](https://github.com/Jaeyo/go-drain3), e contabiliza a quantidade de ocorrências por hora para cada cluster.

---

## 🚀 Funcionalidades

- Analisa um arquivo de log (`exemplo.log`).
- Agrupa mensagens semelhantes com o algoritmo Drain.
- Conta quantas vezes cada **template de log (cluster)** ocorre por hora.
- Mostra um resumo final total por template.

---
O script espera que cada linha do log siga esse padrão:
[Mon May 12 14:01:48 2025] [error] mod_jk child workerEnv in error state 7

O timestamp entre colchetes é extraído para agrupar por hora.

---

## 🧱 Dependências

- Go 1.18+
- Biblioteca `go-drain3` (instalada via Go modules)