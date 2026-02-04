package cmd

import (
	"fmt"
	"os"
)

type CompletionCmd struct {
	Bash CompletionBashCmd `cmd:"" help:"Generate bash completions"`
	Zsh  CompletionZshCmd  `cmd:"" help:"Generate zsh completions"`
	Fish CompletionFishCmd `cmd:"" help:"Generate fish completions"`
}

type CompletionBashCmd struct{}

func (c *CompletionBashCmd) Run() error {
	script := `_delijn_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local commands="version auth config stops lines departures info completion"

    if [ $COMP_CWORD -eq 1 ]; then
        COMPREPLY=($(compgen -W "$commands" -- "$cur"))
    fi
}

complete -F _delijn_completions delijn
`
	fmt.Fprint(os.Stdout, script)

	return nil
}

type CompletionZshCmd struct{}

func (c *CompletionZshCmd) Run() error {
	script := `#compdef delijn

_delijn() {
    local -a commands
    commands=(
        'version:Print version'
        'auth:Manage API key'
        'config:Manage configuration'
        'stops:Search and view stops'
        'lines:Search and view lines'
        'departures:Show realtime departures'
        'info:Show CLI and API info'
        'completion:Generate shell completions'
    )

    _arguments \
        '1: :->command' \
        '*::arg:->args'

    case $state in
        command)
            _describe 'command' commands
            ;;
    esac
}

compdef _delijn delijn
`
	fmt.Fprint(os.Stdout, script)

	return nil
}

type CompletionFishCmd struct{}

func (c *CompletionFishCmd) Run() error {
	script := `complete -c delijn -f

complete -c delijn -n '__fish_use_subcommand' -a 'version' -d 'Print version'
complete -c delijn -n '__fish_use_subcommand' -a 'auth' -d 'Manage API key'
complete -c delijn -n '__fish_use_subcommand' -a 'config' -d 'Manage configuration'
complete -c delijn -n '__fish_use_subcommand' -a 'stops' -d 'Search and view stops'
complete -c delijn -n '__fish_use_subcommand' -a 'lines' -d 'Search and view lines'
complete -c delijn -n '__fish_use_subcommand' -a 'departures' -d 'Show realtime departures'
complete -c delijn -n '__fish_use_subcommand' -a 'info' -d 'Show CLI and API info'
complete -c delijn -n '__fish_use_subcommand' -a 'completion' -d 'Generate shell completions'
`
	fmt.Fprint(os.Stdout, script)

	return nil
}
