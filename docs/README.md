# G8Emu

An emulator reads the original machine code instructions that were assembled for the target machine, interprets them, and then replicates the functionality of the target machine on the host machine. The ROM files contain the instructions, the emulator reads those instructions, and then does work to mimic the original machine.

### Notes

- Opening project in Zed editor for webassembly work: `GOOS=js GOARCH=wasm zed .`
