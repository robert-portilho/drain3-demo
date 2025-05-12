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

## Saída esperada

```bash
Hora: 2025-05-11 04:00
Cluster (2) [error] mod jk child workerEnv in error state <*>: 26
Cluster (3) [notice] jk2 init() Found child <*> in scoreboard slot <*>: 32
Cluster (1) [notice] workerEnv.init() ok /etc/httpd/conf/workers2.properties: 27
Hora: 2025-05-11 05:00
Cluster (3) [notice] jk2 init() Found child <*> in scoreboard slot <*>: 19
Cluster (1) [notice] workerEnv.init() ok /etc/httpd/conf/workers2.properties: 15
Cluster (2) [error] mod jk child workerEnv in error state <*>: 15
Cluster (4) [error] [client <*> Directory index forbidden by rule: /var/www/html/: 1
.
.
.
----------------------- TOTAL ------------------------
Cluster (6) [error] mod jk child init 1 -2: 12
Cluster (3) [notice] jk2 init() Found child <*> in scoreboard slot <*>: 836
Cluster (1) [notice] workerEnv.init() ok /etc/httpd/conf/workers2.properties: 569
Cluster (4) [error] [client <*> Directory index forbidden by rule: /var/www/html/: 32
Cluster (5) [error] jk2 init() Can't find child <*> in scoreboard: 12
Cluster (2) [error] mod jk child workerEnv in error state <*>: 539
