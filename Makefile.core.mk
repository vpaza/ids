FINDFILES=find . \( -path ./.git -o -path ./out -o -path ./.github -o -path ./vendor \) -prune -o -type f
XARGS=xargs -0 -r

lint-copyright:
	@${FINDFILES} \( -name '*.go' -o -name '*.sh' \) \( ! \( -name '*.gen.go' -o -name '*.pb.go' \) \) -print0 |\
		${XARGS} scripts/lint_copyright_license.sh

lint-go:
	@${FINDFILES} -name '*.go' \( ! \( -name '*.gen.go' -o -name '*.pb.go' \) \) -print0 | ${XARGS} scripts/lint_go.sh

lint-markdown:
	@${FINDFILES} -name '*.md' -print0 | ${XARGS} mdl --ignore-front-matter --style .mdl.rb

lint: lint-copyright lint-go lint-markdown
