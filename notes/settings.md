# Settings Goland (IntelliJ)

* GOPATH = The project root folder
* Create a Go Build in Edit configuration
* The name given will be the name of the compiled file. Suggested: oscar-helper
* Choose **Run Kind** as *package*
* Set **Package path** as main. (The executable package name should be main)
* go get github.com/jroimartin/gocui creates the pkg and src folders in root
* create the *main* package folder in the newly created *src* folder
* set output folder as /root/bin

## gitignore
```
bin
pkg
src/github*
go_build_*
```

## Set the following folders as *Excluded*
* bin
* pkg

## Add README.md to bin as the default file.