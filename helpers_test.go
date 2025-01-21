package migrations

func RepositoryBuilder() repositoryBuilder {
	return repositoryBuilder{}
}

type repositoryBuilder struct {
	migration []Migration
}

func (b repositoryBuilder) WithMigration(m ...Migration) repositoryBuilder {
	if b.migration == nil {
		b.migration = make([]Migration, 0)
	}
	b.migration = append(b.migration, m...)
	return b
}

func (b repositoryBuilder) Build() Repository {
	r := Repository{}
	for _, m := range b.migration {
		_ = r.Add(m)
	}
	return r
}
