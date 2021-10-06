//go:generate go run github.com/golang/mock/mockgen -package migrations -destination migration_mock_test.go github.com/jamillosantos/migrations Source,Target,Migration,RunnerReporter

package migrations
