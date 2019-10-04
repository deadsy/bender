
DIRS = da \
       emu \
       gen \

all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./cmd/$$dir $@; \
	done

clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./cmd/$$dir $@; \
	done
