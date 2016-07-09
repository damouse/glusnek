# gosnake

Call from python to go and back. Goroutine/GIL safe bindings with automatic type conversion. 

While either language can call the other, go has to start the show: in other words this library embeds python within go. To embed go within python check out (gopy)[https://github.com/go-python/gopy].

Tested with go1.6, Ubuntu 15.04, python 2.7. No support for python 3 yet. 

## Getting Started

"I dont care why it works, just run the demo"

```
go get https://github.com/damouse/gosnake
cd $GOPATH/src/github.com/damouse/gosnake
export PYTHONPATH=demo
go run demo/main.go
```

If you see a warning about `Python.h` then you need the python dev tools.

```
sudo apt-get install python-dev
```

### Imports and the PATH

`gosnake` imports and calls python packages the same way python does: by checking for the imported package in the `PYTHONPATH`. Lets see an example.

foo.py:

```
def hello():
    return 'Why do humans instinctively fear snakes?'
```

main.py:

```
import foo
print foo.hello()
```

Directory structure 1:

```
demo/
    main.py
    foo/
        hello.py
```

This directory structure works. `main.py` resolves the package `foo` by checking the current directory. 

```
$ cd demo
$ python main.py 
Why do humans instinctively fear snakes?
```

Directory structure 2:

```
demo/
    main.py
foo/
    hello.py
```

This one doesn't, since the package we're looking for is not present in the current directory. If `foo` is an installed package, (i.e. you can run `pip install foo`) then it works again.

An easier way to make the packages visible to python is to temporarily add the target directory to the `PYTHONPATH`.

```
$ cd demo
# export PYTHONPATH="${PYTHONPATH}:.."
$ python main.py 
Why do humans instinctively fear snakes?
```

# TODO

- Tests green on existing work
- go-py type conversion
- py-go type conversion
- Model objects for imports

#### Advanced Features:

* "Permanent", goroutine-safe imports
* Performance Cleanup
    * Dynamically create threads as needed to handle requests
    * Pre-create goroutines for outbound

* Getting object references from python
