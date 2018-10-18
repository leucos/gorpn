# Golang RPN calc

Simple yet efficient terminal RPN calcultator.

## Installation

### Binary

Grab a [release](https://github.com/leucos/gorpn/releases).

### From sources

Using Go 1.11:

```
export GO111MODULES=on
git clone https://github.com/leucos/gorpn
cd gorpn
go install
```

## Supported ops

  - `+`, `-`, `*`, `/`, `%`
  - `pow` (e.g. `2⏎3⏎pow` yields `8`), `sqrt`
  - `sin`, `cos`, `tan`
  - `abs`, `ceil`, `floor`, `round`, `trunc` (e.g. `3.14159⏎2⏎trunc⏎`)
  - `rad`, `deg` for angle modes
  - `dup` (a.k.a empty input and ⏎) duplicates last stack item
  - `swap` exchanges last 2 items in the stack
  - `drop` removes last item in the stack
  - `pi`, `phi` constants
  - `quit` or `<ESC>` exists `gorpn`

Mode is shown in the bottom line. If an error occurs, a red `E` will
show at the bottom right corner.

## Licence

DWTFYWPL

## Authors

@leucos

Inspired by https://medium.com/@jhh3/anonymous-functions-and-reflection-in-go-71274dd9e83a

