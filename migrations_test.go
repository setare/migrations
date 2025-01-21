//go:generate go run go.uber.org/mock/mockgen -package migrations -destination migration_mock_test.go github.com/jamillosantos/migrations/v2 Source,Target,Migration,RunnerReporter

package migrations
