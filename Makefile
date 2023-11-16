## gitadr Makefile
##  
## This Makefile is used to build, test, and install gitadr.
##  
#

.PHONY: help
## Show the help
help:
	@awk '/^(## |##)/ \
			{ if (c) {print c}; c=substr($$0, 4); next } \
				c && /(^[[:alpha:]][[:alnum:]_-]+:)/ \
			{print $$1, "\t", c; c=0} \
				END { print c }' $(MAKEFILE_LIST)

# Keep these lines at the end of the file to keep the help tidy
## 