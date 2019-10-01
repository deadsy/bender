#!/usr/bin/python3

def main():
  f = open('code.txt')
  x = f.readlines()
  f.close

  base = 0x200
  mem = []

  for l in x:
    l = l.strip()
    l = l.split()
    adr = int(l[0], 16)
    if adr != base:
      print("%04x != %04x" % (adr, base))
      return
    buf = [int(l[i + 1], 16) for i in range(len(l) - 1)]
    mem.extend(buf)
    base += 0x10

  for i in range(len(mem)):
    if i % 8 == 0 and i != 0:
      print()
    print("0x%02x, " % mem[i], end = '')
  print()

main()

