</$objtype/mkfile

TARG=storage

BIN=/$objtype/bin
</sys/src/cmd/mkmany

STORAGE=storage.$O fsys.$O match.$O rules.$O

$STORAGE:	$HFILES storage.h

$O.storage:	$STORAGE

syms:V:
	$CC -a storage.c		> syms
	$CC -aa fsys.c match.c rules.c	>>syms
