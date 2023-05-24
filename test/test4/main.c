//-----------------------------------------------------------------------------

#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <ctype.h>
#include <6502.h>
#include <cbm.h>

//-----------------------------------------------------------------------------

#define REG_CLR(ptr, bits) (*(ptr) &= ~(bits))
#define REG_SET(ptr, bits) (*(ptr) |= (bits))

//-----------------------------------------------------------------------------

void *xmemcpy(void *dst, const void *src, size_t n) {
	if (n != 0) {
		char *d = dst;
		const char *s = src;
		do {
			*d++ = *s++;
		} while (--n != 0);
	}
	return dst;
}

//-----------------------------------------------------------------------------

#define BYTES_PER_LINE 8

static void mem_display(uint16_t addr, const void *ptr, size_t n) {
	char ascii[BYTES_PER_LINE + 1];
	size_t ofs = 0;

	n = (n + BYTES_PER_LINE - 1) & ~(BYTES_PER_LINE - 1);
	ascii[BYTES_PER_LINE] = 0;

	while (ofs < n) {
		int i;
		printf("%04x ", addr + ofs);
		for (i = 0; i < BYTES_PER_LINE; i++) {
			uint8_t c = ((uint8_t *) ptr)[ofs];
			printf("%02x ", c);
			// ascii string
			ascii[i] = c;
			c &= 0xe0;
			if ((c == 0) || (c == 0x80)) {
				ascii[i] = '.';
			}
			ofs++;
		}
		printf("%s\n", ascii);
	}
}

//-----------------------------------------------------------------------------

#define DATA_IO (uint8_t *)1

#define ROMSIZE (4 << 10)

int main(void) {
	uint16_t addr = 0xd000;
	uint8_t buf[8];

	while (addr < 0xd000 + ROMSIZE) {
		SEI();
		REG_CLR(DATA_IO, 1 << 2);
		memcpy(buf, (void *)addr, sizeof(buf));
		REG_SET(DATA_IO, 1 << 2);
		CLI();
		mem_display(addr, buf, sizeof(buf));
		addr += sizeof(buf);
	}

	return 0;
}

//-----------------------------------------------------------------------------
