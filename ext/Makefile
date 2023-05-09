TOP = ..

CURRENT_DIR = $(shell pwd)
DL_DIR = $(TOP)/dl
USR_DIR = $(CURRENT_DIR)/usr

# cc65 tool
CC65_VER = 79018fd
CC65_URL = https://github.com/cc65/cc65/tarball/$(CC65_VER)
CC65_FILE = cc65-$(CC65_VER).tar.gz
CC65_TGZ = $(DL_DIR)/$(CC65_FILE)
CC65_DIR = $(TOP)/ext/cc65

.PHONY: all
all: .stamp_cc65

.PHONY: clean
clean:
	-rm -rf $(USR_DIR)
	-rm -rf $(CC65_DIR)
	-rm .stamp_*

$(CC65_TGZ):
	mkdir -p $(DL_DIR)
	wget $(CC65_URL) -O $(CC65_TGZ)

.stamp_cc65: $(CC65_TGZ) 
	mkdir -p $(CC65_DIR)
	tar -C $(CC65_DIR) -zxf $(CC65_TGZ) --strip-components 1
	PREFIX=$(USR_DIR) $(MAKE) -C $(CC65_DIR) all
	mkdir -p $(USR_DIR)
	PREFIX=$(USR_DIR) $(MAKE) -C $(CC65_DIR) install
	touch $@
