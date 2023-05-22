# Running The Functional Tests

```
./cmd/emu/emu -f 6502_functional_test.bin 
6502_functional_test.bin code 0000-ffff
emu> go 400
PC is stuck at 3469, 96241376 cpu cycles
```

Per the code:

```
                        ; S U C C E S S ************************************************
                        ; -------------       
                                success         ;if you get here everything went well
3469 : 4c6934          >        jmp *           ;test passed, no errors


```

