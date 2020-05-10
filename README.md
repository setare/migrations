# migrations

Migrations is an abstraction for a migration system. It can migrate anything.

## How to use

TODO

## Extending

To get it working the system as divided into 4 main components which 2 are interfaces and 2 are concrete implementations:

_**PS**: The examples here will be database related. But, this library is able to migrate anything that migrations would
apply._

### 1. Source

In some migration systems (like [golang-migrate/migrate](https://github.com/golang-migrate/migrate)) are stored as `.sql`
files and are stored into a directory as `<timestamp>_<description>.(up|down).sql`. In other systems, the migration is
a `func` that need to do some complex work and should connect many components before the actual database migration.

So, in the first example, the `Source` is a directory containing a bunch of .sql files with specific names. In the second
example, the `Source` are function that should be organized chronologically.

Hence, `Source` is the media that persists the migrations themselves. In practice, it is just an `interface{}` with a
bunch of methods that will list all available migrations ([check the code]()).

**TODO**: Link the Source interface.

### 2. Target

A `Target` are what the migrations are transforming. If you are dealing with relational databases, like postgres, you would
use a `TargetSQL` implementation (we provide one, check the our examples folder|**TODO**).

In practice, a `Target` is just an `interface{}` with a bunch of methods that will list executed migrations, mark and
unmark migrations as executed ([check the code]()).

**TODO**: Link the Target interface.

### 3. Executer

An Executer integrations `Source` and `Target` and is responsible for step actions, like `Do` and `Undo`. Each call will
step forward or backward one migration at a time.

### 4. Runner

Runners are, also, concrete. They capture the developer intentions and call the `Executer`.

Let's say that you want to migrate your system. By that, you mean to run all pending migrations. So the runner will
use the `Executer.Do` calling it multiple times to get all migrations executed.
