package server

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/nixpig/syringe.sh/internal/serrors"
	"github.com/nixpig/syringe.sh/internal/stores"
)

func NewIdentityMiddleware(s *stores.SystemStore) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			publicKeyHash := fmt.Sprintf("%x", sha1.Sum(sess.PublicKey().Marshal()))
			sess.Context().SetValue("publicKeyHash", publicKeyHash)

			sessionID := sess.Context().SessionID()
			username := sess.Context().User()
			email := time.Now().GoString()

			cmd := sess.Command()
			if len(cmd) > 0 && cmd[0] == "register" {
				log.Debug("register user", "session", sessionID, "username", username, "email", email, "publicKeyHash", publicKeyHash)
				next(sess)
			}

			user, err := s.GetUser(username)
			if err != nil || user == nil {
				sess.Write([]byte(fmt.Sprintf("User '%s' with provided key not found.\n", username)))

				log.Debug(
					"creating user",
					"session", sessionID,
					"username", username,
					"email", email,
					"publicKeyHash", publicKeyHash,
					"reason", err,
				)

				// Currently only supporting a user with a single key, since to add a
				// new key we'd need to authenticate them with a previously saved key
				// first, which seems like a pain this early on - so let's do it later

				user = &stores.User{
					Username:      username,
					PublicKeySHA1: publicKeyHash,
					Email:         email,
				}
				userID, err := s.CreateUser(user)
				if err != nil {
					log.Error("failed to create user", "session", sessionID, "err", err)

					sess.Stderr().Write([]byte(serrors.New("user", fmt.Sprintf("failed to create user '%s'", username), sessionID).Error()))

					sess.Exit(1)
					return
				}

				sess.Write([]byte(
					fmt.Sprintf("Created user '%s' with provided key.\n", username),
				))

				user.ID = userID
			}

			if user.PublicKeySHA1 != publicKeyHash {
				sess.Stderr().Write([]byte("Public key provided does not match stored"))
				sess.Exit(1)
			}

			// TODO: add email verification

			next(sess)
		}
	}
}
