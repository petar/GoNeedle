# Copyright 2010 GoNeedle Authors. All rights reserved.
# Use of this source code is governed by a 
# license that can be found in the LICENSE file.

# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# After editing the DIRS= list or adding imports to any Go files
# in any of those directories, run:
#
#	./deps.bash
#
# to rebuild the dependency information in Make.deps.

nullstring :=
space := $(nullstring) # a space at the end
ifndef GOBIN
QUOTED_HOME=$(subst $(space),\ ,$(HOME))
GOBIN=$(QUOTED_HOME)/bin
endif
QUOTED_GOBIN=$(subst $(space),\ ,$(GOBIN))

all: install

DIRS=\
	needle/proto\
	needle\
	cmd/needle-daemon\
	cmd/needle-listen\
	cmd/needle-connect\

TEST=\
	$(filter-out $(NOTEST),$(DIRS))

BENCH=\
	$(filter-out $(NOBENCH),$(TEST))

clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
nuke.dirs: $(addsuffix .nuke, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))
bench.dirs: $(addsuffix .bench, $(BENCH))

%.clean:
	+cd $* && $(QUOTED_GOBIN)/gomake clean

%.install:
	+cd $* && $(QUOTED_GOBIN)/gomake install

%.nuke:
	+cd $* && $(QUOTED_GOBIN)/gomake nuke

%.test:
	+cd $* && $(QUOTED_GOBIN)/gomake test

%.bench:
	+cd $* && $(QUOTED_GOBIN)/gomake bench

clean: clean.dirs

install: install.dirs

test:	test.dirs

bench:	bench.dirs

nuke: nuke.dirs

deps:
	./deps.bash

-include Make.deps
