# gitadr: cli to manage your adrs the gitops way

## Description

This is the cli that helps you manage your architecture decision records (adrs) the gitops way. 

## Install

### Package Managers

- [ ] Add to homebrew
- [ ] Add to scoop
- [ ] Add to chocolatey
- [ ] Add to apt
- [ ] Add to yum

### Manual

You can download the latest release from the [releases](https://github.com/git-adr/gitadr/releases) page.

## Build

To build the cli you will want to use the Makefile. For more information run `make help`.

## Development

### Viewer

The adr viewer server is best tested using a supported adr repository. You can either use the [git-adr](https://github.com/git-adr/adr)
or your own adr repository. The `adr` directory is a great place to clone the repository to since its ignored by git.

```bash
git clone https://github.com/git-adr/adr.git adr
```

To start the development server it's best to use air. Air will watch for changes and restart the server automatically.

```bash
go get -u github.com/cosmtrek/air
air
```

#### Tailwind

The viewer uses tailwind for styling. To make changes to the styling you will need to install the dependencies.

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
chmod +x tailwindcss-macos-arm64
mv tailwindcss-macos-arm64 /usr/local/bin/tailwindcss
```

To run the watcher for tailwind run the following command.

```bash
tailwindcss -i assets/input.css -o build/assets/output.css --watch
```
