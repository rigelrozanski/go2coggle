# go2coggle

_generate basic coggle diagram from go repo_

---

### Installation

1. Make sure you [have Go installed][1] and [put $GOPATH/bin in your $PATH][2]
2. [Install Cobra][3]
3. run `go get github.com/rigelrozanski/go2coggle`
4. run `go install go2coggle`

[1]: https://golang.org/doc/install
[2]: https://github.com/tendermint/tendermint/wiki/Setting-GOPATH 
[3]: https://github.com/spf13/cobra#installing

###  Usage

once installed run the command `go2coggle` from any golang repo directory to create a coggle file loadable into https://coggle.it
(drag and drop the generated .txt file into a new coggle diagram)

### Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

### License

go2coggle is released under the Apache 2.0 license.
