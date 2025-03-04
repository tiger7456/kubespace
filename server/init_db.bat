@echo ##########Formatting code#########################################
go fmt ./ && go vet ./
@echo ##########Format the code successfully############################
SET CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
@echo ##########Compiling gva.exe#######################################
go build -o gva.exe cmd/main.go
@echo ##########Successfully compiled gva.exe###########################
@echo ##########Initializing data using gva.exe#########################
gva.exe initdb
@echo ##########Use gva.exe to initialize data successfully#############
@echo ##########Deleting gva.exe########################################
del gva.exe
@echo ##########Deleting gva.exe successfully###########################