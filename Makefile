user-rpc-dev:
	@make -f deploy/mk/user-rpc.mk release-test

user-api-dev:
	@make -f deploy/mk/user-api.mk release-test

release-test:user-rpc-dev user-api-dev

install-server:
	cd ./deploy/scripts && chomd +x release-test.sh && ./release-test.sh