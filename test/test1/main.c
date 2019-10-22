
#include <stdio.h>
#include <conio.h>

void foo(void) {
  char tmp[2];
  int i;
  tmp[1] = 0;
  for (i = 0; i < 10; i ++) {
    tmp[0] = '0' + i;
    puts(tmp);
  }
}

void main(void) {
  puts("Hello World!");
  foo();
}
