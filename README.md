## Install

1. Install ffmpeg: `brew install ffmpeg`
2. Move the .app file to your Applications folder

## Build

#### Install Go
```sh
brew install go
```


#### Create a Go module

```sh
go mod init video2gif
go mod tidy
```

#### Install Fyne CLI

```sh
go install fyne.io/tools/cmd/fyne@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

#### Build .app bundle

```sh
fyne package -os darwin -icon icon.png -name video2gif
```
