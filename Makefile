default: frontend_lint

frontend_prod: export NODE_ENV=production
frontend_prod: frontend

frontend: node_modules
	node ci/build.mjs

frontend_lint: node_modules
	./node_modules/.bin/eslint \
		--fix \
		src

node_modules:
	npm ci --include dev
