## Examples

### Running Zork:
```sh
$ go run main.go ro:zork1.cpm

a>zork1
```

### Running with the CP/M utilities:
```sh
$ go run main.go ro:cpma.cpm

a>dir
A: DUMP     COM : SDIR     COM : SUBMIT   COM : ED       COM
A: STAT     COM : BYE      COM : RMAC     COM : CREF80   COM
A: DDTZ     COM : L80      COM : M80      COM : SID      COM
A: WM       COM : WM       HLP : ZSID     COM : MAC      COM
A: TRACE    UTL : HIST     UTL : LIB80    COM : M        SUB
A: DDT      COM : CLS      COM : MOVCPM   COM : ASM      COM
A: LOAD     COM : XSUB     COM : HELLO    ASM : PIP      COM
A: SYSGEN   COM : REW      ASM : HELLO    COM : REW      COM
```

### Creating a new disk image for use with the emulator:
```sh
# Create a blank image of the correct size
$ dd if=/dev/zero of=newdisk.cpm bs=128 count=2002

# Format it for use with CP/M
$ mkfs.cpm newdisk.cpm

# Start the emulator with the new disk. cpma.cpm will be A:, newdisk.cpm will be B:
$ go run main.go ro:cpma.cpm newdisk.cpm
```

Disk images and CP/M binary obtained from https://github.com/jscrane/cpm80
