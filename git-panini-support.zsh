compdef _git-pannini git-panini
function _git-panini() {
  local context curcontext="$curcontext" state line
  typeset -A opt_args

  _arguments -C \
    "1: :->subcmd" \
    "*::arg:->args"

  case $state in
    (subcmd)
      _values 'subcommand' world apply noop fetch find-nonpanini status
      ;;
    (args)
      case $line[1] in
        (world)
          _arguments \
            '--verbose[more verbose]' \
            '--local[show only local]'
          ;;
        (apply|noop)
          _arguments \
            '--verbose[more verbose]' \
            '--force[force apply]'
          ;;
        (fetch|find-nonpanini|status)
          _message 'no more arguments'
          ;;
      esac
      ;;
  esac
}

compdef _cdp cdp
function cdp() {
  cd `git panini path panini:$1`
}

function _cdp() {
  local -U panini_names
  panini_names=(`git panini world --local`)
  _values 'panini repos' ${panini_names:s/panini://}
}
