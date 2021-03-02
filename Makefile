build:
	@echo -e "\e[96mBuilding for \e[95mARM (version 7)\e[39m"
	@echo -e "\e[96mBuilding \e[93mBinary\e[39m"
	@export GHW_DISABLE_WARNINGS=1
	@go build -o bin/robot_controller cmd/main.go 1>/dev/null
	@echo -e "\e[92mBuild Complete\e[39m"
build-ui:
	@if [[ ! -e webserver/www/index.html ]]; then\
    	echo "Robot Controller - Need to do \"make build\" or similar build (for other platforms) before using GUI!" > webserver/www/index.html;\
		exit 1
	fi
	@cd webserver && statik -src=www -f 1>/dev/null
run:
	@bin/robot_controller
install-dep:
	@go get -u github.com/rakyll/statik
	@mkdir -p bin
	@mkdir -p "webserver/www"
	@go get -u -v all