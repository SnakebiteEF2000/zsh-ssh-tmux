# zsh-ssh-tmux

depends on tmux, zsh and fzf

source in .zshrc

~~~
go build -o ssh-conf-gen main.go
~~~

~~~
./ssh-conf-gen -inventory /path/to/inventory.yml -user <user used for ssh>
~~~