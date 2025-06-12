## Build

#### Create a Go module

```sh
go mod init video2gif
go get fyne.io/fyne/v2
go mod tidy
```

## Install Fyne CLI

```sh
go install fyne.io/tools/cmd/fyne@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Build .app bundle

```sh
fyne package -os darwin -icon icon.png -name video2gif
```
