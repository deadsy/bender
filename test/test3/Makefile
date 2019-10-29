
BASE = /opt/cc65

XCC = $(BASE)/bin/cc65
XAS = $(BASE)/bin/ca65
XLD = $(BASE)/bin/ld65
SIM = $(BASE)/bin/sim65

#TGT = c64
TGT = sim6502

OUT = bcdtest

CFLAGS := -g
CFLAGS += -O -Oi -Or -Os

all:
	$(XCC) $(CFLAGS) -t $(TGT) main.c
	$(XAS) -g main.s
	$(XAS) -g bcd.S
	$(XLD) --mapfile $(OUT).map -o $(OUT) -t $(TGT) main.o bcd.o $(TGT).lib

test:
	$(SIM) -c $(OUT)

clean:
	-rm *.s
	-rm *.o
	-rm $(OUT)
	-rm $(OUT).map
