package cobra


// CreateAppFunc initializes dependencies.
type CreateAppFunc func(*Command) (*App, error)

// App contains dependencies to execute a generator.
type App struct {
	Generator Generator
	UI        cli.UI
}
