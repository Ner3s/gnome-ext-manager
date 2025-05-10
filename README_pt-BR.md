# Gerenciador de Extensões GNOME

Uma ferramenta de linha de comando para backup e restauração de extensões do GNOME Shell.

[EN](README.md)

## Funcionalidades

- **Backup de Extensões**: Crie um backup de todas as suas extensões GNOME ativadas junto com suas configurações
- **Restauração de Extensões**: Restaure extensões previamente salvas e suas configurações
- **Indicadores de Progresso**: Barras de progresso visuais para mostrar o status da operação
- **Interface Simples**: Interface de linha de comando fácil de usar

## Requisitos

- GNOME Shell
- Linguagem de programação Go
- Ferramenta de linha de comando `gnome-extensions`
- Ferramenta de linha de comando `dconf`

## Instalação

Clone este repositório e construa a aplicação:

```bash
curl -L https://github.com/ner3s/gnome-ext-manager/releases/latest/download/gnome-ext-manager -o /usr/local/bin/gnome-ext-manager
chmod +x /usr/local/bin/gnome-ext-manager
```

## Uso

Execute a aplicação:

```bash
gnome-ext-manager
```

### Backup de Suas Extensões

Selecione a opção 1 no menu principal e digite um nome de arquivo para seu backup (ou pressione Enter para o nome padrão).

O processo de backup irá:
1. Detectar sua versão do GNOME Shell
2. Listar todas as extensões ativadas
3. Salvar os dados da extensão em um arquivo JSON
4. Exportar configurações de extensão usando dconf
5. Criar um arquivo zip contendo todos os dados do backup

### Restauração de Suas Extensões

Selecione a opção 2 no menu principal e forneça o caminho para seu arquivo zip de backup.

O processo de restauração irá:
1. Extrair o arquivo zip
2. Instalar cada extensão a partir de extensions.gnome.org
3. Ativar todas as extensões previamente ativadas
4. Restaurar configurações de extensão usando dconf

## Como Funciona

A funcionalidade de backup:
- Usa `gnome-extensions list --enabled` para identificar extensões instaladas
- Busca metadados de extensão de extensions.gnome.org
- Usa `dconf dump` para exportar configurações de extensão
- Empacota tudo em um arquivo zip

A funcionalidade de restauração:
- Extrai o arquivo de backup
- Baixa extensões de extensions.gnome.org
- Instala-as com `gnome-extensions install`
- Ativa cada extensão com `gnome-extensions enable`
- Restaura configurações com `dconf load`

## Licença

Este projeto está licenciado sob os termos do arquivo LICENSE incluído neste repositório.

## Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para enviar um Pull Request.
