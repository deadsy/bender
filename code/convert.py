#!/usr/bin/python3


def code_dump():
  f = open('code.txt')
  x = f.readlines()
  f.close

  base = 0x200
  mem = []

  for l in x:
    l = l.strip()
    l = l.split()
    adr = int(l[0], 16)
    assert adr == base
    buf = [int(l[i + 1], 16) for i in range(len(l) - 1)]
    mem.extend(buf)
    base += 0x10

  for i in range(len(mem)):
    if i % 8 == 0 and i != 0:
      print()
    print("0x%02x, " % mem[i], end = '')
  print()

def symbol_dump():
  f = open('symbols.txt')
  x = f.readlines()
  f.close

  table = []

  for l in x:
    l = l.strip()
    l = l.split()
    assert len(l) == 2
    adr = int(l[1], 16)
    name = l[0].strip().lower()
    table.append((adr, name))

  table = sorted(table, key=lambda x: x[0])
  for (adr, name) in table:
    print("0x%04x: \"%s\"," % (adr, name))

def main():
  code_dump()
  symbol_dump()

main()

