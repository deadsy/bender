//-----------------------------------------------------------------------------

#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <ctype.h>
#include <6502.h>
#include <cbm.h>

//-----------------------------------------------------------------------------

static void outs(const char *s) {
	while (*s != 0) {
		putchar(*s++);
	}
}

static void *xmemcpy(void *dst, const void *src, size_t n) {
	char *d = dst;
	const char *s = src;
	int i;
	for (i = 0; i < n; i++) {
		d[i] = s[i];
	}
	return dst;
}

//-----------------------------------------------------------------------------

static char nybble(uint8_t val) {
	val &= 0xf;
	if (val >= 0xa && val <= 0xf) {
		return val - 0xa + 'a';
	}
	return val + '0';
}

static char *hex8(char *s, uint8_t val) {
	s[0] = nybble(val >> 4);
	s[1] = nybble(val);
	s[2] = 0;
	return s;
}

static char *hex16(char *s, uint16_t val) {
	hex8(s, val >> 8);
	hex8(&s[2], val);
	return s;
}

//-----------------------------------------------------------------------------

#define BYTES_PER_LINE 8

static void mem_display8(uint16_t addr, void *ptr, size_t n) {
	char ascii[BYTES_PER_LINE + 1];
	char tmp[4 + 1];
	size_t ofs = 0;

	n = (n + BYTES_PER_LINE - 1) & ~(BYTES_PER_LINE - 1);
	ascii[BYTES_PER_LINE] = 0;

	while (ofs < n) {
		int i;
		outs(hex16(tmp, addr + ofs));
		outs(" ");
		for (i = 0; i < BYTES_PER_LINE; i++) {
			uint8_t c = ((uint8_t *) ptr)[ofs];
			outs(hex8(tmp, c));
			outs(" ");
			// ascii string
			ascii[i] = c;
			c &= 0xe0;
			if ((c == 0) || (c == 0x80)) {
				ascii[i] = '.';
			}
			ofs++;
		}
		outs(ascii);
		outs("\n");
	}
}

//-----------------------------------------------------------------------------

#define REG_CLR(ptr, bits) (*(ptr) &= ~(bits))
#define REG_SET(ptr, bits) (*(ptr) |= (bits))

//-----------------------------------------------------------------------------

#define DATA_IO (uint8_t *)1

int main(void) {
	uint16_t addr = 0xd000;
	uint8_t buf[8];

	while (addr < 0xd000 + (4 << 10)) {
		SEI();
		REG_CLR(DATA_IO, 1 << 2);
		xmemcpy(buf, (void *)addr, sizeof(buf));
		REG_SET(DATA_IO, 1 << 2);
		CLI();
		mem_display8(addr, buf, sizeof(buf));
		addr += sizeof(buf);
	}

	return 0;
}
