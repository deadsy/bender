
DIRS = cmd \
       test \

.PHONY: all
all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done

.PHONY: clean
clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done

# clean + remove the 3rd party libraries
.PHONY: clobber
clobber: clean
	$(MAKE) -C ext clean
	-rm .stamp_ext

# build 3rd party tools/libraries
.stamp_ext:
	$(MAKE) -C ext
	touch $@

