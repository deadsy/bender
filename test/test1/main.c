
#include <stdio.h>
#include <conio.h>
#include <limits.h>
#include <stdint.h>

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

static char *hex32(char *s, uint32_t val) {
	hex16(s, val >> 16);
	hex16(&s[4], val);
	return s;
}

static char *itoa(char *s, int val) {
	unsigned int uval;
	int i = 0;
	int j = 0;

	// make it positive
	if (val < 0) {
		uval = -val;
	} else {
		uval = val;
	}

	// digits in reverse order
	do {
		s[i++] = (uval % 10) + '0';
		uval /= 10;
	} while (uval != 0);

	// add negative sign
	if (val < 0) {
		s[i++] = '-';
	}
	// null terminate
	s[i] = 0;
	i--;

	// reverse the string
	while (j < i) {
		char tmp = s[j];
		s[j++] = s[i];
		s[i--] = tmp;
	}

	return s;
}

void main(void) {
	char tmp[32];

	printf("sizeof(char) %d bytes\n", sizeof(char));
	printf("sizeof(int) %d bytes\n", sizeof(int));
	printf("sizeof(long) %d bytes\n", sizeof(long));

	puts(itoa(tmp, 0));
	puts(itoa(tmp, 1234));
	puts(itoa(tmp, -1234));
	puts(itoa(tmp, INT_MAX));
	puts(itoa(tmp, INT_MIN));
	puts(hex8(tmp, 0xAB));
	puts(hex16(tmp, 0xABCD));
	puts(hex32(tmp, 0xDEADBEEF));
}
