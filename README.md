# SnapTUI 📸

Uma aplicação TUI (Terminal User Interface) profissional para backup de bancos de dados PostgreSQL, desenvolvida em Go com Bubble Tea.

## 🚀 Funcionalidades

- **Interface Terminal Moderna**: TUI intuitiva e responsiva
- **Conexão PostgreSQL**: Configuração fácil de conexão com banco
- **Backup Múltiplo**: Seleção individual ou de todos os bancos
- **Progresso Visual**: Spinner animado durante operações
- **Relatório Completo**: Resumo detalhado com sucessos e erros
- **Arquitetura Profissional**: Código organizado em packages

## 🛠️ Instalação

### Pré-requisitos

```bash
# Ubuntu/Debian
sudo apt install postgresql-client

# CentOS/RHEL/Fedora
sudo yum install postgresql
# ou
sudo dnf install postgresql
```

### Compilação

```bash
git clone https://github.com/Luiz-F3lipe/snapTUI.git
cd snapTUI
go build -o snapTUI ./cmd/main.go
```

### Execução

```bash
./snapTUI
```

## 📁 Estrutura do Projeto

```
snapTUI/
├── cmd/
│   └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── backup/              # Serviços de backup
│   │   └── backup.go
│   ├── config/              # Configurações e estilos
│   │   └── config.go
│   ├── database/            # Serviços de banco de dados
│   │   └── database.go
│   ├── types/               # Tipos e estruturas
│   │   └── types.go
│   └── ui/                  # Interface do usuário
│       ├── model.go         # Lógica principal da TUI
│       └── views/           # Views/telas da aplicação
│           └── views.go
├── go.mod
├── go.sum
└── README.md
```

## 🎯 Como Usar

### 1. Configuração da Conexão
- Configure host, porta, usuário, senha e banco de dados
- Use **Tab** ou **↑/↓** para navegar entre campos
- **Espaço** limpa o campo atual
- **Enter** para conectar

### 2. Menu Principal
- **Fazer Backup**: Acessa lista de bancos para backup
- **Restaurar Backup**: (Em desenvolvimento)
- **Configurar Conexão**: Volta para tela de configuração
- **Sair**: Encerra a aplicação

### 3. Seleção de Bancos
- **Espaço** para selecionar/desselecionar bancos
- **All Databases** seleciona todos de uma vez
- **Enter** inicia o backup dos bancos selecionados

### 4. Progresso e Resultados
- Spinner animado durante o processo
- Relatório final com sucessos e erros
- Lista dos arquivos de backup criados

## ⌨️ Atalhos de Teclado

| Tecla | Ação |
|-------|------|
| `↑/↓` ou `k/j` | Navegar |
| `Enter` | Selecionar/Confirmar |
| `Espaço` | Marcar/Desmarcar |
| `Esc` | Voltar |
| `Tab` | Próximo campo (conexão) |
| `Q` ou `Ctrl+C` | Sair |

## 🏗️ Arquitetura

### Packages

- **`cmd/`**: Ponto de entrada da aplicação
- **`internal/backup/`**: Lógica de backup com pg_dump
- **`internal/config/`**: Configurações, cores e estilos
- **`internal/database/`**: Operações de banco de dados
- **`internal/types/`**: Definições de tipos e estruturas
- **`internal/ui/`**: Interface e lógica da TUI

### Padrões Utilizados

- **Clean Architecture**: Separação clara de responsabilidades
- **Service Pattern**: Serviços especializados para cada domínio
- **MVC Pattern**: Model-View-Controller para a TUI
- **Dependency Injection**: Injeção de dependências entre services

## 🔧 Tecnologias

- **[Go](https://golang.org/)**: Linguagem principal
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: Framework TUI
- **[Bubbles](https://github.com/charmbracelet/bubbles)**: Componentes UI
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)**: Estilização
- **[pq](https://github.com/lib/pq)**: Driver PostgreSQL

## 📝 Formato dos Backups

Os backups são salvos no diretório do executável com o formato:
```
<nome_do_banco>_YYYYMMDD_HHMMSS.backup
```

Exemplo: `meu_banco_20231030_143022.backup`

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para detalhes.

## ✨ Desenvolvido por

**Luiz Felipe** - [GitHub](https://github.com/Luiz-F3lipe)
