go mod init example_program
go get github.com/BalazsNyiro/jyp

# https://stackoverflow.com/a/63730888


# elegant update method - use the wanted commit hash:
go get github.com/BalazsNyiro/jyp@a7e4925a3ef3f10c91ed9b3dedd37d3c3b8c733c
go mod vendor

# Manual update method:
# if you want to update the necessary commit:
 - update the commit hash in go.mod
 - the 'go build will display the correct timestamp that you have to use
 - go mod tidy
 - go mod vendor
