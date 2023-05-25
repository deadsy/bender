//-----------------------------------------------------------------------------

//#include <string.h>
//#include <ctype.h>
//#include <6502.h>
//#include <cbm.h>

#include <stdint.h>
#include <stdio.h>
#include <conio.h>
#include <time.h>

#include <c64.h>

typedef unsigned char u8;
typedef char s8;

//-----------------------------------------------------------------------------

// return the current screen ram location
uint8_t *get_screenram(void) {
	return (uint8_t *) ((VIC.addr >> 4) << 10);
}

#define SCREEN_RAM 0x0400
#define SCREEN_COLOR 0xD800

static void set_xy(u8 x, u8 y, u8 c, u8 color) {
	uint8_t *ofs = (uint8_t *) ((40 * y) + x);
	ofs[SCREEN_RAM] = c;
	ofs[SCREEN_COLOR] = color;
}

//-----------------------------------------------------------------------------

typedef s8 xy[2];
typedef xy layout[4];

#define BG_COLOR COLOR_BLACK

#define P_CHAR 0xa0

static const layout p_layout[] = {
	// I
	{{0, 0}, {1, 0}, {2, 0}, {-1, 0}},
	{{0, 0}, {0, 1}, {0, 2}, {0, -1}},
	{{0, 0}, {-1, 0}, {-2, 0}, {1, 0}},
	{{0, 0}, {0, -1}, {0, -2}, {0, 1}},
	// J
	{{0, 0}, {-1, 0}, {-1, -1}, {1, 0}},
	{{0, 0}, {0, -1}, {1, -1}, {0, 1}},
	{{0, 0}, {1, 0}, {1, 1}, {-1, 0}},
	{{0, 0}, {0, 1}, {-1, 1}, {0, -1}},
	// L
	{{0, 0}, {1, 0}, {1, -1}, {-1, 0}},
	{{0, 0}, {0, 1}, {1, 1}, {0, -1}},
	{{0, 0}, {-1, 0}, {-1, 1}, {1, 0}},
	{{0, 0}, {0, -1}, {-1, -1}, {0, 1}},
	// S
	{{0, 0}, {0, -1}, {1, -1}, {-1, 0}},
	{{0, 0}, {1, 0}, {1, 1}, {0, -1}},
	{{0, 0}, {0, 1}, {-1, 1}, {1, 0}},
	{{0, 0}, {-1, 0}, {-1, -1}, {0, 1}},
	// Z
	{{0, 0}, {0, -1}, {-1, -1}, {1, 0}},
	{{0, 0}, {1, 0}, {1, -1}, {0, 1}},
	{{0, 0}, {0, 1}, {1, 1}, {-1, 0}},
	{{0, 0}, {-1, 0}, {-1, 1}, {0, -1}},
	// T
	{{0, 0}, {0, -1}, {-1, 0}, {1, 0}},
	{{0, 0}, {1, 0}, {0, -1}, {0, 1}},
	{{0, 0}, {0, 1}, {1, 0}, {-1, 0}},
	{{0, 0}, {-1, 0}, {0, 1}, {0, -1}},
	// O
	{{0, 0}, {0, -1}, {1, -1}, {1, 0}},
	{{0, 0}, {0, -1}, {1, -1}, {1, 0}},
	{{0, 0}, {0, -1}, {1, -1}, {1, 0}},
	{{0, 0}, {0, -1}, {1, -1}, {1, 0}},
};

static const u8 p_color[] = {
	COLOR_LIGHTBLUE,	// I
	COLOR_BLUE,		// J
	COLOR_ORANGE,		// L
	COLOR_GREEN,		// S
	COLOR_RED,		// Z
	COLOR_PURPLE,		// T
	COLOR_YELLOW,		// O
};

static void piece(u8 x, u8 y, u8 color, const layout * p) {
	int i;
	for (i = 0; i < 4; i++) {
		set_xy(x + (*p)[i][0], y + (*p)[i][1], P_CHAR, color);
	}
}

static void piece_on(u8 id, u8 x, u8 y, u8 dirn) {
	piece(x, y, p_color[id], &p_layout[(id << 2) + (dirn & 3)]);
}

static void piece_off(u8 id, u8 x, u8 y, u8 dirn) {
	piece(x, y, BG_COLOR, &p_layout[(id << 2) + (dirn & 3)]);
}

static void delay(void) {
	volatile int i;
	for (i = 0; i < 3000; i++) {
	}
}

void game(void) {
	u8 dirn = 0;
	u8 y;

	clrscr();
	bgcolor(BG_COLOR);

	y = 3;

	while (1) {
		piece_on(0, 2, y, dirn);
		piece_on(1, 7, y, dirn);
		piece_on(2, 11, y, dirn);
		piece_on(3, 15, y, dirn);
		piece_on(4, 19, y, dirn);
		piece_on(5, 23, y, dirn);
		piece_on(6, 27, y, dirn);
		delay();
		piece_off(0, 2, y, dirn);
		piece_off(1, 7, y, dirn);
		piece_off(2, 11, y, dirn);
		piece_off(3, 15, y, dirn);
		piece_off(4, 19, y, dirn);
		piece_off(5, 23, y, dirn);
		piece_off(6, 27, y, dirn);
		dirn += 1;
		y += 1;
		if (y == 20) {
			y = 3;
		}
	}

}

//-----------------------------------------------------------------------------

int main(void) {
	game();
	while (1) ;
	return 0;
}

//-----------------------------------------------------------------------------
