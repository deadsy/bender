//-----------------------------------------------------------------------------
/*

C&N Tetris

*/
//-----------------------------------------------------------------------------

//#include <stdio.h>
#include <conio.h>
//#include <time.h>

#include <c64.h>

typedef unsigned int u16;
typedef unsigned char u8;
typedef signed char s8;
typedef signed int s16;

//-----------------------------------------------------------------------------

// return the current screen ram location
u8 *get_screenram(void) {
	return (u8 *) ((VIC.addr >> 4) << 10);
}

#define SCREEN_RAM 0x0400
#define SCREEN_COLOR 0xD800

//-----------------------------------------------------------------------------

typedef s8 xy[2];
typedef xy layout[3];

#define BG_COLOR COLOR_BLACK

#define P_ON 0xa0
#define P_OFF 0x20

static const layout p_layout[] = {
	// I
	{{1, 0}, {2, 0}, {-1, 0}},
	{{0, 40}, {0, 80}, {0, -40}},
	{{-1, 0}, {-2, 0}, {1, 0}},
	{{0, -40}, {0, -80}, {0, 40}},
	// J
	{{-1, 0}, {-1, -40}, {1, 0}},
	{{0, -40}, {1, -40}, {0, 40}},
	{{1, 0}, {1, 40}, {-1, 0}},
	{{0, 40}, {-1, 40}, {0, -40}},
	// L
	{{1, 0}, {1, -40}, {-1, 0}},
	{{0, 40}, {1, 40}, {0, -40}},
	{{-1, 0}, {-1, 40}, {1, 0}},
	{{0, -40}, {-1, -40}, {0, 40}},
	// S
	{{0, -40}, {1, -40}, {-1, 0}},
	{{1, 0}, {1, 40}, {0, -40}},
	{{0, 40}, {-1, 40}, {1, 0}},
	{{-1, 0}, {-1, -40}, {0, 40}},
	// Z
	{{0, -40}, {-1, -40}, {1, 0}},
	{{1, 0}, {1, -40}, {0, 40}},
	{{0, 40}, {1, 40}, {-1, 0}},
	{{-1, 0}, {-1, 40}, {0, -40}},
	// T
	{{0, -40}, {-1, 0}, {1, 0}},
	{{1, 0}, {0, -40}, {0, 40}},
	{{0, 40}, {1, 0}, {-1, 0}},
	{{-1, 0}, {0, 40}, {0, -40}},
	// O
	{{0, -40}, {1, -40}, {1, 0}},
	{{0, -40}, {1, -40}, {1, 0}},
	{{0, -40}, {1, -40}, {1, 0}},
	{{0, -40}, {1, -40}, {1, 0}},
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

static void set(u16 addr, u8 c, u8 color) {
	u8 *ofs = (u8 *) addr;
	ofs[SCREEN_RAM] = c;
	ofs[SCREEN_COLOR] = color;
}

static void piece(u8 x, u8 y, u8 c, u8 color, const layout * p) {
	s16 base = (40 * y) + x;
	set(base, c, color);
	set(base + (*p)[0][0] + (*p)[0][1], c, color);
	set(base + (*p)[1][0] + (*p)[1][1], c, color);
	set(base + (*p)[2][0] + (*p)[2][1], c, color);
}

static void piece_on(u8 id, u8 x, u8 y, u8 dirn) {
	piece(x, y, P_ON, p_color[id], &p_layout[(id << 2) + (dirn & 3)]);
}

static void piece_off(u8 id, u8 x, u8 y, u8 dirn) {
	piece(x, y, P_OFF, BG_COLOR, &p_layout[(id << 2) + (dirn & 3)]);
}

static void delay(void) {
	volatile int i;
	for (i = 0; i < 6000; i++) {
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
