


build:
	cd web && pnpm run build
	rm -rf server/core/genCode/dist && mv web/dist server/core/genCode/dist