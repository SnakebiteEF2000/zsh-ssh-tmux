# generate ssh config from ansible inventory *fast*

~~~
go build -o ssh-conf-gen main.go
mv ssh-conf-get ~/.local/bin/
~~~

~~~
# default setting
ssh-conf-gen -inventory <ansible inventory.yml> -user <default ssh user>

# use alternative user
ssh-conf-gen -inventory <ansible inventory.yml> -user <default ssh user> -altuser <alternative ssh user> -altuserregex <use regex>
~~~
# zsh-ssh-tmux

depends on tmux, zsh and fzf
