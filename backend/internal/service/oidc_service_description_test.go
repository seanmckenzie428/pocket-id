package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pocket-id/pocket-id/backend/internal/dto"
	"github.com/pocket-id/pocket-id/backend/internal/model"
	testutils "github.com/pocket-id/pocket-id/backend/internal/utils/testing"
)

func TestOidcService_CreateClient_WithDescription(t *testing.T) {
	db := testutils.NewDatabaseForTest(t)

	s := &OidcService{
		db: db,
	}

	description := "This is a test description for the OIDC client"
	input := dto.OidcClientCreateDto{
		OidcClientUpdateDto: dto.OidcClientUpdateDto{
			Name:        "Test Client",
			Description: &description,
			CallbackURLs: []string{"https://example.com/callback"},
		},
		ID: "test-client-id",
	}

	client, err := s.CreateClient(t.Context(), input, "test-user-id")
	require.NoError(t, err)

	// Verify description was saved
	require.NotNil(t, client.Description)
	assert.Equal(t, description, *client.Description)

	// Verify by fetching from DB
	var fetchedClient model.OidcClient
	err = db.First(&fetchedClient, "id = ?", client.ID).Error
	require.NoError(t, err)
	require.NotNil(t, fetchedClient.Description)
	assert.Equal(t, description, *fetchedClient.Description)
}

func TestOidcService_CreateClient_WithoutDescription(t *testing.T) {
	db := testutils.NewDatabaseForTest(t)

	s := &OidcService{
		db: db,
	}

	input := dto.OidcClientCreateDto{
		OidcClientUpdateDto: dto.OidcClientUpdateDto{
			Name:         "Test Client No Description",
			CallbackURLs: []string{"https://example.com/callback"},
			// Description is nil
		},
		ID: "test-client-no-desc",
	}

	client, err := s.CreateClient(t.Context(), input, "test-user-id")
	require.NoError(t, err)

	// Verify description is nil
	assert.Nil(t, client.Description)

	// Verify by fetching from DB
	var fetchedClient model.OidcClient
	err = db.First(&fetchedClient, "id = ?", client.ID).Error
	require.NoError(t, err)
	assert.Nil(t, fetchedClient.Description)
}

func TestOidcService_UpdateClient_Description(t *testing.T) {
	db := testutils.NewDatabaseForTest(t)

	s := &OidcService{
		db: db,
	}

	// Create initial client without description
	initialClient := model.OidcClient{
		Base: model.Base{
			ID: "test-update-client",
		},
		Name:         "Initial Name",
		CallbackURLs: model.UrlList{"https://example.com/callback"},
	}
	err := db.Create(&initialClient).Error
	require.NoError(t, err)

	// Update with description
	description := "Updated description"
	updateInput := dto.OidcClientUpdateDto{
		Name:         "Updated Name",
		Description:  &description,
		CallbackURLs: []string{"https://example.com/callback"},
	}

	updatedClient, err := s.UpdateClient(t.Context(), initialClient.ID, updateInput)
	require.NoError(t, err)

	// Verify description was updated
	require.NotNil(t, updatedClient.Description)
	assert.Equal(t, description, *updatedClient.Description)
	assert.Equal(t, "Updated Name", updatedClient.Name)

	// Verify in DB
	var fetchedClient model.OidcClient
	err = db.First(&fetchedClient, "id = ?", initialClient.ID).Error
	require.NoError(t, err)
	require.NotNil(t, fetchedClient.Description)
	assert.Equal(t, description, *fetchedClient.Description)
}

func TestOidcService_UpdateClient_ClearDescription(t *testing.T) {
	db := testutils.NewDatabaseForTest(t)

	s := &OidcService{
		db: db,
	}

	// Create initial client with description
	description := "Initial description"
	initialClient := model.OidcClient{
		Base: model.Base{
			ID: "test-clear-desc",
		},
		Name:         "Test Client",
		Description:  &description,
		CallbackURLs: model.UrlList{"https://example.com/callback"},
	}
	err := db.Create(&initialClient).Error
	require.NoError(t, err)

	// Update with nil description (clear it)
	updateInput := dto.OidcClientUpdateDto{
		Name:         "Test Client",
		Description:  nil, // Clear description
		CallbackURLs: []string{"https://example.com/callback"},
	}

	updatedClient, err := s.UpdateClient(t.Context(), initialClient.ID, updateInput)
	require.NoError(t, err)

	// Verify description was cleared
	assert.Nil(t, updatedClient.Description)

	// Verify in DB
	var fetchedClient model.OidcClient
	err = db.First(&fetchedClient, "id = ?", initialClient.ID).Error
	require.NoError(t, err)
	assert.Nil(t, fetchedClient.Description)
}