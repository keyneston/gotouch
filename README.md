# gotouch

`gotouch` is a simple command line tool to automate the creation of new go
files and their matching test file.

It will make any required directories, then create the go file with the
(hopefully) correct `package <name>`.

```shell
$ gotouch place/thing.go
$ ls place
thing.go      thing_test.go
$ cat place/*
package place
package place
```
