# go-i8080: an Intel 8080 emulator

### Usage
[![GoDoc](https://godoc.org/github.com/GinjaNinja32/go-i8080?status.svg)](https://godoc.org/github.com/GinjaNinja32/go-i8080)  
See the [`example/`](https://github.com/GinjaNinja32/go-i8080/tree/master/example) folder for an example main file for the emulator which loads CP/M 2.2 and user-selected disk image files.

### Details
- BIOS-level emulation of CP/M
- Disk support (somewhat hardcoded, see `diskio.go`)
- Can run CP/M 2.2 (see `example/` folder)

### TODOs and limitations
- Disk format is currently hardcoded
- Instruction dispatch could probably be optimised
- Split Run() into Run() and Step(), to allow per-instruction stepping
