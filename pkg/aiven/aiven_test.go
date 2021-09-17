package aiven

import (
	aiven_nais_io_v1 "github.com/nais/liberator/pkg/apis/aiven.nais.io/v1"
	"github.com/nais/nais-cli/pkg/client"
	"github.com/nais/nais-cli/pkg/common"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
	"time"
)

var scheme = runtime.NewScheme()

func TestAivenGenerateApplicationCreated(t *testing.T) {

	client.InitScheme(scheme)

	username := "user"
	team := "team"
	pool := "pool"
	secretName := "secret-name"
	expiry := 1

	namespace := v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: team,
		},
	}

	fakeClient := fake.NewFakeClientWithScheme(scheme, &namespace)
	aiven := SetupAiven(fakeClient, username, team, pool, secretName, expiry)
	currentAivenApp, err := aiven.GenerateApplication()
	assert.NoError(t, err)

	assert.Equal(t, username, currentAivenApp.Name, "Name has the same value")
	assert.Equal(t, team, currentAivenApp.Namespace, "Namespace has the same value")
	assert.Equal(t, secretName, currentAivenApp.Spec.SecretName, "SecretName has the same value")
	assert.Equal(t, pool, currentAivenApp.Spec.Kafka.Pool, "Pool has the same value")

	parsedDate, err := time.Parse(time.RFC3339, currentAivenApp.Spec.ExpiresAt)
	assert.NoError(t, err)
	assert.True(t, parsedDate.After(time.Now()), "Parsed date is still valid")
}

func TestAivenGenerateApplicationUpdated(t *testing.T) {

	client.InitScheme(scheme)

	username := "user"
	team := "team"
	pool := "pool"
	secretName := "secret-name"
	expiry := 1

	aivenApp := aiven_nais_io_v1.AivenApplication{
		ObjectMeta: metav1.ObjectMeta{
			Name:      username,
			Namespace: team,
		},
	}

	namespace := v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: team,
		},
	}

	fakeClient := fake.NewFakeClientWithScheme(scheme, &namespace, &aivenApp)
	aiven := SetupAiven(fakeClient, username, team, pool, secretName, expiry)
	currentAivenApp, err := aiven.GenerateApplication()
	assert.NoError(t, err)

	assert.Equal(t, username, currentAivenApp.Name, "Name has the same value")
	assert.Equal(t, team, currentAivenApp.Namespace, "Namespace has the same value")
	assert.Equal(t, secretName, currentAivenApp.Spec.SecretName, "SecretName has the same value")
	assert.Equal(t, pool, currentAivenApp.Spec.Kafka.Pool, "Pool has the same value")

	parsedDate, err := time.Parse(time.RFC3339, currentAivenApp.Spec.ExpiresAt)
	assert.NoError(t, err)
	assert.True(t, parsedDate.After(time.Now()), "Parsed date is still valid")
}

func TestAiven_SetSecretName(t *testing.T) {
	username := "user"
	team := "team"
	secretName := "secret-name"

	s, err := common.SetSecretName(secretName, username, team)
	assert.NoError(t, err)
	assert.Equal(t, secretName, s, "SecretName has the same value as input")

	s, err = common.SetSecretName("", username, team)
	assert.NoError(t, err)
	assert.Equal(t, "user-team-3d735979", s, "SecretName is generated")
}
