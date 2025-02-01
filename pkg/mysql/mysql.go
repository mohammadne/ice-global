package mysql

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	*sqlx.DB
}

const (
	driver      = "mysql"
	pingTimeout = time.Second * 20
)

func Open(cfg *Config) (*Mysql, error) {
	connString := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	database, err := sqlx.Open(driver, connString)
	if err != nil {
		return nil, fmt.Errorf("error while opening connection to mysql: %v", err)
	}
	database.SetMaxIdleConns(0)

	ctx, cf := context.WithTimeout(context.Background(), pingTimeout)
	defer cf()
	if err = database.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error while pinging database: %v", err)
	}

	r := &Mysql{DB: database}

	return r, nil
}

type MigrateDirection string

const (
	MigrateDirectionUp   MigrateDirection = "UP"
	MigrateDirectionDown MigrateDirection = "DOWN"

	migrateDirectionUpFileExtension   string = ".up.sql"
	migrateDirectionDowbFileExtension string = ".down.sql"
)

type migration struct {
	content string
	sort    int
}

func (m *Mysql) Migrate(directory string, files *embed.FS, direction MigrateDirection) error {
	fileEntries, err := files.ReadDir(directory)
	if err != nil {
		slog.Error(`error reading migration files`, `Err`, err)
		os.Exit(1)
	}
	if len(fileEntries) == 0 {
		return fmt.Errorf("no migration files has been given")
	}

	migrations := make([]migration, 0, len(fileEntries)/2)

	extension := migrateDirectionUpFileExtension
	if direction == MigrateDirectionDown {
		extension = migrateDirectionDowbFileExtension
	}

	for _, file := range fileEntries {
		sort, _ := strconv.Atoi(strings.Split(file.Name(), "_")[0])
		path := directory + "/" + file.Name()

		if path[len(path)-len(extension):] == extension {
			content, err := files.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading migration file:%v", err)
			}

			migrations = append(migrations, migration{
				content: string(content),
				sort:    sort,
			})
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].sort < migrations[j].sort
	})

	for fileIndex := 0; fileIndex < len(migrations); fileIndex++ {
		queries, err := splitQueries(migrations[fileIndex].content)
		if err != nil {
			return fmt.Errorf("error splitting queries of SQL file:%v", err)
		}

		for queryIndex := 0; queryIndex < len(queries); queryIndex++ {
			fmt.Println("------------------------")
			fmt.Println(queries[queryIndex])
			fmt.Println("------------------------")

			_, err := m.DB.Exec(queries[queryIndex])
			if err != nil {
				return fmt.Errorf("error migrating file:%v", err)
			}
		}
	}

	return nil
}

func splitQueries(content string) ([]string, error) {
	var queries []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	var query strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Add the current line to the query builder
		query.WriteString(line + "\n")

		// If we encounter a semicolon, it means the query is complete
		if strings.HasSuffix(line, ";") {
			queries = append(queries, query.String())
			query.Reset() // Reset the builder for the next query
		}
	}

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading SQL file:%v", err)
	}

	return queries, nil
}
