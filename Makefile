build:
	@echo -e "\e[96mBuilding for \e[95mARM (version 7)\e[39m"
	@echo -e "\e[96mBuilding \e[93mBinary\e[39m"
	@export GHW_DISABLE_WARNINGS=1
	@go build -o bin/robot_controller cmd/main.go 1>/dev/null
	@echo -e "\e[92mBuild Complete\e[39m"
run:
	@bin/robot_controller
install-deps:
	@go get -u github.com/d2r2/go-i2c
	@go get -u github.com/stianeikeland/go-rpio