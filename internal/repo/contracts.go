// Package repo implements application outer layer logic. Each logic group in own file.
package repo

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test
