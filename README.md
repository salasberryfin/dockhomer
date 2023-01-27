# Dockhomer

Dockhomer is a personal project that I start with the goal of understanding how Linux namespaces are used to create containers.

This playground for container creation and configuration allows me to get my hands dirty with the low-level features that Docker uses to make containers simple and effective.
It starts with creating processes within the Linux system isolated from the 'host' machine using kernel namespaces:
- **UTS**: hostname and NIS domain name
- **PID**: process IDs
- **User**: user and group IDs
- **Mount**: mount points

The following are the features I'm focusing on -for now-:
[x] Running an interactive terminal - `/bin/bash`
[ ] Passing a given command as argument - `echo "hello, world"`
[ ] Run a program - `go run main.go`
[ ] Use custom images - `dockhomer run nginx`

