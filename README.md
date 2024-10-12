# TreeTags

Generate [ctags](https://ctags.io) compatible tag file for your code. Powered by [tree-sitter](https://tree-sitter.github.io/tree-sitter/).

### Installation

Clone repo and build the project and copy executable to a folder in your path.

```bash
go build # build the project
cp ./tree-tags ~/bin # copy executable to a folder in your path

## USAGE
cd ~/your-go-project
tree-tags # will output a tags file in vim compatible format
```
