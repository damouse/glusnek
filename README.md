# gosnake

Bi-directional, high level language bindings with automatic type conversion between go and python. Features goroutine-safe python invocation!

Tested on Ubuntu 15.04 with python 2.7. No support for python 3 yet. 

## Running

Execute the demo script:

```
export PYTHONPATH=demo
go run demo/main.go
```

If you see a warning about `Python.h` then you don't have the python dev tools.

```
sudo apt-get install python-dev
```
