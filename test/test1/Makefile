
BASE = /opt/cc65

XCC = $(BASE)/bin/cc65
XAS = $(BASE)/bin/ca65
XLD = $(BASE)/bin/ld65
SIM = $(BASE)/bin/sim65

#TGT = c64
TGT = sim6502

OUT = test1

all:
	$(XCC) -g -Oi -t $(TGT) main.c
	$(XAS) -g main.s
	$(XLD) --mapfile $(OUT).map -o $(OUT) -t $(TGT) main.o $(TGT).lib

test:
	$(SIM) -c $(OUT)

clean:
	-rm *.s
	-rm *.o
	-rm $(OUT)
	-rm $(OUT).map
