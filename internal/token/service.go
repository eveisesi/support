package token

import (
	"context"
	"fmt"
	"time"

	"github.com/embersyndicate/support"
	"github.com/embersyndicate/support/internal/key"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Service interface {
	BuildAndSignUserKey(ctx context.Context, user *support.User) ([]byte, error)
	ParseAndVerifyToken(context.Context, string) (jwt.Token, error)
	GetUserIDFromToken(t jwt.Token) (string, error)
}

type service struct {
	key key.Service
}

func New(
	key key.Service,
) Service {
	return &service{
		key: key,
	}
}

func (s *service) BuildAndSignUserKey(ctx context.Context, user *support.User) ([]byte, error) {

	now := time.Now().In(time.UTC)
	t := jwt.New()
	var err error
	err = t.Set(jwt.JwtIDKey, uuid.New().String())
	if err != nil {
		return nil, fmt.Errorf("failed to set %s on token: %w", jwt.JwtIDKey, err)
	}

	err = t.Set(jwt.SubjectKey, fmt.Sprintf("SUPPORT::USER::%s", user.ID.Hex()))
	if err != nil {
		return nil, fmt.Errorf("failed to set %s on token: %w", jwt.SubjectKey, err)
	}

	err = t.Set(jwt.AudienceKey, "Ember Syndicate Support Portal Users")
	if err != nil {
		return nil, fmt.Errorf("failed to set %s on token: %w", jwt.AudienceKey, err)
	}

	err = t.Set(jwt.IssuerKey, "Ember Syndicate")
	if err != nil {
		return nil, fmt.Errorf("failed to set %s on token: %w", jwt.IssuerKey, err)
	}

	err = t.Set(jwt.IssuedAtKey, now.Unix())
	if err != nil {
		return nil, fmt.Errorf("failed to set %s on token: %w", jwt.IssuedAtKey, err)
	}

	err = t.Set(jwt.ExpirationKey, now.Add(time.Hour*8).Unix())
	if err != nil {
		return nil, fmt.Errorf("failed to set %s key  on token: %w", jwt.ExpirationKey, err)
	}

	err = t.Set(`username`, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to set %s on token: %w", "username", err)
	}

	err = t.Set(`id`, user.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to set %s on token: %w", "user id", err)
	}

	signed, err := jwt.Sign(t, jwa.RS256, s.key.GetPrivateJWK())
	if err != nil {
		return nil, err
	}

	return signed, nil

}

func (s *service) ParseAndVerifyToken(ctx context.Context, t string) (jwt.Token, error) {

	seg := newrelic.FromContext(ctx).StartSegment("parse and verify token")
	defer seg.End()

	set, err := s.getSet()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jwk set from key service: %w", err)
	}

	token, err := jwt.ParseString(t, jwt.WithKeySet(set))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token, nil

}

func (s *service) GetUserIDFromToken(t jwt.Token) (string, error) {

	id, ok := t.Get("id")
	if !ok {
		return "", fmt.Errorf("failed to retrieve id from token")
	}

	if _, ok := id.(string); !ok {
		return "", fmt.Errorf("id of type %T is invalid", id)
	}

	return id.(string), nil

}

// Returns a *jwk.Set that ParseToken uses to validate a JWT
func (s *service) getSet() (*jwk.Set, error) {

	// ctx := context.Background()

	// Leaving this commented out because one day in the possible not to distance future
	// we may have an SSO Server to reachout to and this code will help
	// Future Self, reference zrule for the rest of the code needed for this

	// result, err := s.redis.Get(ctx, zrule.CACHE_CCP_JWKS).Bytes()
	// if err != nil && err.Error() != "redis: nil" {
	// 	return nil, errors.Wrap(err, "unexpected error looking for jwk in redis")
	// }

	// if len(result) == 0 {
	// 	res, err := s.client.Get(s.jwksURL)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "unable to retrieve jwks from sso")
	// 	}

	// 	if res.StatusCode != http.StatusOK {
	// 		return nil, fmt.Errorf("unexpected status code recieved while fetching jwks. %d", res.StatusCode)
	// 	}

	// 	buf, err := ioutil.ReadAll(res.Body)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "faile dto read jwk response body")
	// 	}

	// 	_, err = s.redis.Set(ctx, zrule.CACHE_CCP_JWKS, buf, time.Hour*24).Result()
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "failed to cache jwks in redis")
	// 	}

	// 	result = buf
	// }

	set, err := s.key.GetPublicJWKSBytes()
	if err != nil {
		return nil, err
	}

	return jwk.ParseBytes(set)
}
