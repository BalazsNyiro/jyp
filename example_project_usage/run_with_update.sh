cd ..
gitId=$(git show HEAD --pretty=format:"%H" --no-patch)
cd -
go get github.com/BalazsNyiro/jyp@${gitId}
go mod vendor

go run jyp_example_usage.go
