function zsh-ssh-tmux () {
  local selected_host=$(awk '/^Host / {
    host = $2
    desc = ""
    for (i=3; i<=NF; i++) desc = desc " " $i
    if (host != "*") print host "\t" desc
  }' ~/.ssh/config | fzf-tmux -p --reverse -1 -0 +m +s --query "$LBUFFER" --prompt="SSH Host > " | cut -f1)

  if [ -n "$selected_host" ]; then
    BUFFER="tmux new-window -n ssh-${selected_host} 'ssh ${selected_host}'"
    zle accept-line
  fi
  zle reset-prompt
}

zle -N zsh-ssh-tmux
bindkey '^n' zsh-ssh-tmux
