package user_test

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nixpig/syringe.sh/internal/user"
	"github.com/nixpig/syringe.sh/pkg/turso"
	"github.com/nixpig/syringe.sh/pkg/validation"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
)

type MockTursoClient struct {
	mock.Mock
}

func (m *MockTursoClient) CreateDatabase(name, group string) (*turso.TursoDatabaseResponse, error) {
	args := m.Called(name, group)

	return args.Get(0).(*turso.TursoDatabaseResponse), args.Error(1)
}

func (m *MockTursoClient) ListDatabases() (*turso.TursoDatabases, error) {
	args := m.Called()

	return args.Get(0).(*turso.TursoDatabases), args.Error(1)
}

func (m *MockTursoClient) CreateToken(name, expiration string) (*turso.TursoToken, error) {
	args := m.Called(name, expiration)

	return args.Get(0).(*turso.TursoToken), args.Error(1)
}

func (m *MockTursoClient) New(organization, apiToken string, httpClient http.Client) turso.TursoClient {
	m.Called(organization, apiToken, httpClient)

	return turso.TursoClient{}
}

var mockTursoClient = new(MockTursoClient)

func TestUserCmd(t *testing.T) {
	scenarios := map[string]func(
		t *testing.T,
		cmd *cobra.Command,
		service user.UserService,
		mock sqlmock.Sqlmock,
	){
		"test user register happy path": testUserRegisterHappyPath,
	}

	for scenario, fn := range scenarios {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: \n%s", err)
		}

		store := user.NewSqliteUserStore(db)

		service := user.NewUserServiceImpl(
			store,
			validation.New(),
			http.Client{},
			user.TursoAPISettings{URL: "", Token: ""},
			mockTursoClient,
		)

		cmd := user.NewCmdUser()

		t.Run(scenario, func(t *testing.T) {
			fn(
				t,
				cmd,
				service,
				mock,
			)
		})
	}
}

func testUserRegisterHappyPath(
	t *testing.T,
	cmd *cobra.Command,
	service user.UserService,
	mock sqlmock.Sqlmock,
) {

}
