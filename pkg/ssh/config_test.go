package ssh_test

import (
	"os"
	"testing"

	"github.com/nixpig/syringe.sh/pkg/ssh"
	"github.com/stretchr/testify/require"
)

func TestSSHConfig(t *testing.T) {
	scenarios := map[string]func(t *testing.T){
		"test add identity to ssh config new host":      testAddIdentityToSSHConfigNewHost,
		"test add identity to ssh config existing host": testAddIdentityToSSHConfigExistingHost,
	}

	for scenario, fn := range scenarios {
		t.Run(scenario, fn)
	}
}

func testAddIdentityToSSHConfigNewHost(t *testing.T) {
	os.Setenv("APP_HOST", "localhost")
	f, err := os.CreateTemp("", "tmp_ssh_config")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	defer f.Close()

	id := "../../test/crypt_test_rsa"

	err = ssh.AddIdentityToSSHConfig(id, f)

	// read contents of file and check
	w, err := os.ReadFile(f.Name())
	require.NoError(t, err)

	require.Equal(
		t,
		"Host localhost\nAddKeysToAgent yes\nIgnoreUnknown UseKeychain\nUseKeychain yes\nIdentityFile ../../test/crypt_test_rsa\n",
		string(w),
	)
}

func testAddIdentityToSSHConfigExistingHost(t *testing.T) {
	os.Setenv("APP_HOST", "localhost")
	f, err := os.CreateTemp("", "tmp_ssh_config")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	defer f.Close()

	f.WriteString("Host localhost\nAddKeysToAgent yes\nIgnoreUnknown UseKeychain\nUseKeychain yes\n")

	id := "../../test/crypt_test_rsa"

	err = ssh.AddIdentityToSSHConfig(id, f)

	w, err := os.ReadFile(f.Name())
	require.NoError(t, err)

	require.Equal(
		t,
		"Host localhost\nAddKeysToAgent yes\nIgnoreUnknown UseKeychain\nUseKeychain yes\nIdentityFile ../../test/crypt_test_rsa\n",
		string(w),
	)
}
