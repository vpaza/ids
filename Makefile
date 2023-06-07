RUN=./scripts/run.sh
MAKE_CONTAINER=$(RUN) make --no-print-directory -e -f Makefile.core.mk
FRONTEND_MAKE_CONTAINER=$(RUN) make --no-print-directory -C frontend -e -f Makefile

%:
	@$(MAKE_CONTAINER) $@

default:
	@$(MAKE_CONTAINER)

shell:
	@$(RUN) /bin/bash

local-run:
	@$(MAKE_CONTAINER) dev

backend-run:
	MISC_OPTIONS="-p 3000:3000" $(MAKE_CONTAINER) dev

frontend-run:
	MISC_OPTIONS="-p 5174:5174" $(FRONTEND_MAKE_CONTAINER) dev

.PHONY: default shell