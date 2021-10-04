# Golang RPN calc

**DEPRECATED** 

See https://github.com/leucos/grpn instead.

[![CircleCI](https://circleci.com/gh/leucos/gorpn/tree/master.svg?style=svg)](https://circleci.com/gh/leucos/gorpn/tree/master)

Simple yet efficient terminal RPN calcultator.

[![asciicast](https://asciinema.org/a/207322.png)](https://asciinema.org/a/207322)

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
  - `pow`, `^` (e.g. `2⏎3⏎pow⏎` yields `8`), `sqrt`
  - `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, 
  - `abs`, `ceil`, `floor`, `round`, `trunc` (e.g. `3.14159⏎2⏎trunc⏎`)
  - `rad`, `deg` for angle modes
  - `dup` (a.k.a empty input and ⏎) duplicates last stack item
  - `swap` exchanges last 2 items in the stack
  - `drop` removes last item in the stack
  - `pi`, `phi` constants
  - `precision`, `#` (e.g. `2⏎precision⏎` or ``2⏎#⏎`) limits number of displayed digits
  - `cu1_cu2` returns exchange rate between `cu1` and `cu2` (e.g. `isk_eur` for xchange rate between ISK and EUR)
  - `quit` or `<ESC>` exists `gorpn`
  - `<UP>/<DOWN>` key walks input history up or down

Mode is shown in the bottom line. If an error occurs, a red `E` will
show at the bottom right corner.

## TODO

- [ ] undo
- [ ] non-intercative mode support (stdin / stdout)
- [ ] sentence recall (vs token recall) with PgUp/PgDown
- [ ] stack save on exit / load on start
- [ ] completion

## Licence

DWTFYWPL

## Authors

@leucos

Inspired by https://medium.com/@jhh3/anonymous-functions-and-reflection-in-go-71274dd9e83a
