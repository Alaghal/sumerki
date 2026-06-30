package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"sumerki/backend/internal/db"
	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/repository"
	"sumerki/backend/internal/service"

	"golang.org/x/crypto/bcrypt"
)

const (
	defaultDatabaseURL = "postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable"
	devPassword        = "password123"
)

type seedAccount struct {
	Email   string
	Name    string
	Culture string
	Patron  *string
	Boost   bool
}

var seedAccounts = []seedAccount{
	{
		Email:   "northern@example.com",
		Name:    "Воронья Сечь",
		Culture: service.CultureNorthernPrincipality,
		Patron:  stringPtr("old_pact"),
		Boost:   true,
	},
	{
		Email:   "lizard@example.com",
		Name:    "Тёплый Камень",
		Culture: service.CultureLizardGrad,
		Patron:  stringPtr("independent"),
	},
	{
		Email:   "posad@example.com",
		Name:    "Серый Посад",
		Culture: service.CultureFreePosad,
		Patron:  stringPtr("empire_of_dusk"),
	},
	{
		Email:   "raider@example.com",
		Name:    "Чёрный Брод",
		Culture: service.CultureNorthernPrincipality,
		Patron:  nil,
	},
}

func main() {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}

	database, err := db.Open(databaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if err := db.Ping(ctx, database); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	seeder := newDevSeeder(database)
	for _, account := range seedAccounts {
		result, err := seeder.Seed(ctx, account)
		if err != nil {
			log.Fatalf("seed %s: %v", account.Email, err)
		}

		patron := "none"
		if account.Patron != nil {
			patron = *account.Patron
		}
		fmt.Printf("OK %s user=%s kingdom=%s patron=%s\n", account.Email, result.User.ID, result.Kingdom.ID, patron)
	}

	fmt.Println("Dev seed complete. Login with password: password123")
}

type devSeeder struct {
	db                    *sql.DB
	kingdoms              *repository.KingdomRepository
	rulers                *service.RulerService
	resources             *service.ResourcesService
	buildings             *service.BuildingService
	army                  *service.ArmyService
	patrons               *service.PatronService
	kingdomsWithTasks     *service.KingdomService
	newKingdomCreatedTime time.Time
}

type seedResult struct {
	User    domain.User
	Kingdom domain.Kingdom
}

func newDevSeeder(database *sql.DB) *devSeeder {
	kingdoms := repository.NewKingdomRepository(database)
	rulers := service.NewRulerService(kingdoms, repository.NewRulerRepository(database))
	resources := service.NewResourcesService(kingdoms, repository.NewResourcesRepository(database))
	buildings := service.NewBuildingService(kingdoms, repository.NewBuildingRepository(database), resources)
	army := service.NewArmyService(kingdoms, repository.NewArmyRepository(database), resources, buildings)
	patrons := service.NewPatronService(kingdoms, repository.NewPatronRepository(database))
	patronPressure := service.NewPatronPressureService(kingdoms, repository.NewPatronRepository(database), repository.NewPatronPressureRepository(database), resources)
	resources.SetProductionProvider(buildings)
	patrons.SetPressureLifecycle(patronPressure)

	return &devSeeder{
		db:                    database,
		kingdoms:              kingdoms,
		rulers:                rulers,
		resources:             resources,
		buildings:             buildings,
		army:                  army,
		patrons:               patrons,
		kingdomsWithTasks:     service.NewKingdomService(kingdoms, rulers, resources, buildings, army),
		newKingdomCreatedTime: time.Now().UTC().Add(-48 * time.Hour),
	}
}

func (s *devSeeder) Seed(ctx context.Context, account seedAccount) (seedResult, error) {
	user, err := s.upsertUser(ctx, account.Email)
	if err != nil {
		return seedResult{}, err
	}

	kingdom, err := s.ensureKingdom(ctx, user.ID, account.Name, account.Culture)
	if err != nil {
		return seedResult{}, err
	}

	if err := s.ensureSupportRows(ctx, user.ID, kingdom.ID); err != nil {
		return seedResult{}, err
	}
	if err := s.applyDevState(ctx, kingdom.ID, account.Boost); err != nil {
		return seedResult{}, err
	}
	if err := s.applyPatron(ctx, user.ID, kingdom.ID, account.Patron); err != nil {
		return seedResult{}, err
	}

	updated, err := s.kingdoms.FindByUserID(ctx, user.ID)
	if err != nil {
		return seedResult{}, err
	}

	return seedResult{User: user, Kingdom: updated}, nil
}

func (s *devSeeder) upsertUser(ctx context.Context, email string) (domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(devPassword), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	var id string
	err = s.db.QueryRowContext(ctx, `SELECT id::text FROM users WHERE lower(email) = lower($1)`, email).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return repository.NewUserRepository(s.db).Create(ctx, email, string(hash))
	}
	if err != nil {
		return domain.User{}, err
	}

	const query = `
		UPDATE users
		SET email = $2,
			password_hash = $3,
			updated_at = now()
		WHERE id = $1
		RETURNING id::text, email, password_hash, created_at, updated_at
	`
	var user domain.User
	err = s.db.QueryRowContext(ctx, query, id, email, string(hash)).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *devSeeder) ensureKingdom(ctx context.Context, userID string, name string, culture string) (domain.Kingdom, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		created, err := s.kingdomsWithTasks.Create(ctx, userID, name, culture)
		if err != nil {
			return domain.Kingdom{}, err
		}
		return s.updateKingdomForSeed(ctx, created.ID, name, culture)
	}
	if err != nil {
		return domain.Kingdom{}, err
	}
	return s.updateKingdomForSeed(ctx, kingdom.ID, name, culture)
}

func (s *devSeeder) updateKingdomForSeed(ctx context.Context, kingdomID string, name string, culture string) (domain.Kingdom, error) {
	const query = `
		UPDATE kingdoms
		SET name = $2,
			culture = $3,
			created_at = LEAST(created_at, $4),
			updated_at = now()
		WHERE id = $1
		RETURNING id::text, user_id::text, name, culture, patron, dread, honor, created_at, updated_at
	`
	return scanKingdom(s.db.QueryRowContext(ctx, query, kingdomID, name, culture, s.newKingdomCreatedTime))
}

func (s *devSeeder) ensureSupportRows(ctx context.Context, userID string, kingdomID string) error {
	if _, err := s.rulers.Current(ctx, userID); err != nil {
		return err
	}
	if _, err := s.resources.CreateForKingdom(ctx, kingdomID); err != nil {
		return err
	}
	if err := s.buildings.EnsureForKingdom(ctx, kingdomID); err != nil {
		return err
	}
	return s.army.EnsureForKingdom(ctx, kingdomID)
}

func (s *devSeeder) applyDevState(ctx context.Context, kingdomID string, boost bool) error {
	if _, err := s.db.ExecContext(ctx, `
		UPDATE kingdom_resources
		SET gold = 2000,
			food = 1500,
			wood = 1500,
			stone = 1200,
			population = 250,
			last_calculated_at = now(),
			updated_at = now()
		WHERE kingdom_id = $1
	`, kingdomID); err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, `
		UPDATE kingdom_units AS units
		SET amount = seed.amount,
			updated_at = now()
		FROM (VALUES
			('militia', 40::bigint),
			('scouts', 10::bigint),
			('spearmen', 15::bigint),
			('archers', 15::bigint),
			('cavalry', 5::bigint)
		) AS seed(unit_type, amount)
		WHERE units.kingdom_id = $1
		  AND units.unit_type = seed.unit_type
	`, kingdomID); err != nil {
		return err
	}

	if boost {
		_, err := s.db.ExecContext(ctx, `
			UPDATE kingdom_buildings AS buildings
			SET level = seed.level,
				upgrade_started_at = NULL,
				upgrade_finishes_at = NULL,
				updated_at = now()
			FROM (VALUES
				('farm', 2),
				('lumberyard', 2),
				('quarry', 2),
				('market', 2),
				('barracks', 2),
				('walls', 1)
			) AS seed(type, level)
			WHERE buildings.kingdom_id = $1
			  AND buildings.type = seed.type
		`, kingdomID)
		return err
	}

	return nil
}

func (s *devSeeder) applyPatron(ctx context.Context, userID string, kingdomID string, patron *string) error {
	if patron != nil {
		if _, err := s.patrons.Join(ctx, userID, *patron); err != nil {
			return err
		}
		return nil
	}

	if _, err := s.patrons.Break(ctx, userID); err != nil {
		return err
	}
	if _, err := s.kingdoms.UpdatePatronByID(ctx, kingdomID, nil); err != nil {
		return err
	}
	return nil
}

func scanKingdom(row scanner) (domain.Kingdom, error) {
	var kingdom domain.Kingdom
	var patron sql.NullString
	err := row.Scan(
		&kingdom.ID,
		&kingdom.UserID,
		&kingdom.Name,
		&kingdom.Culture,
		&patron,
		&kingdom.Dread,
		&kingdom.Honor,
		&kingdom.CreatedAt,
		&kingdom.UpdatedAt,
	)
	if err != nil {
		return domain.Kingdom{}, err
	}
	if patron.Valid {
		kingdom.Patron = &patron.String
	}
	return kingdom, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func stringPtr(value string) *string {
	return &value
}
