// JWT authentication layer.
//
// HTTP is stateless: after POST /login, the next request knows nothing about the
// previous one. Two ways to fix that:
//
//   - Server-side sessions: store a random session ID in a DB/Redis and look it
//     up on every request. Revocable, but needs storage and a lookup per request.
//   - JWT (what we use): hand the client a signed, self-describing claim. The
//     server stores nothing and only verifies a signature, so there is no shared
//     session store and the app scales horizontally for free.
//
// The cost of that choice: we cannot revoke a token. Logout is client-side only
// (the client drops the token); the token itself stays valid until it expires.
// The "exp" claim is therefore our only damage-control lever - see GenerateToken.
package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// secretKey is the HMAC key. HS256 is symmetric: the same key signs AND
// verifies, which makes this the master password of the whole auth system.
// Anyone holding it can forge a token claiming any userID. So it must never be
// hardcoded - source goes to git, git goes to GitHub, and a leaked secret cannot
// be rotated out of history.
//
// []byte and not string because the HMAC signer type-asserts the key to []byte;
// a string compiles fine and then fails at runtime with ErrInvalidKeyType.
var secretKey []byte

// InitJWT loads the signing key from the environment. Called once from main()
// before the server starts, which log.Fatal's on error.
//
// Why an explicit init step instead of reading the env var inside GenerateToken?
// To fail fast and loudly:
//
//   - A missing secret kills the process at boot, when you will notice, instead
//     of surfacing as a 500 on the first login attempt at 2am.
//   - More importantly, the empty check closes a silent auth bypass. HMAC with an
//     empty key works perfectly well - it signs, it verifies, everything looks
//     healthy - but the key is trivially guessable, so anyone can forge tokens.
//     Without this guard, a missing JWT_SECRET is total auth bypass with no
//     visible symptom.
//
// Tradeoff we accept: nothing enforces that this runs before the functions
// below. With a single main() that is fine; a second entrypoint that forgot the
// call would leave secretKey nil.
func InitJWT() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return errors.New("JWT_SECRET environment variable is not set")
	}
	secretKey = []byte(secret)
	return nil
}

// GenerateToken mints a token after a successful login. Runs once per login -
// compare VerifyToken, which runs on every protected request.
//
// HS256 (symmetric, one shared secret) rather than RS256/ES256 because a single
// service both issues and verifies here. Asymmetric keys are for when some OTHER
// service must verify tokens without being trusted to mint them. RSA is not
// "stronger" for this shape of problem, just more moving parts.
//
// Note the payload is base64-encoded, NOT encrypted: anyone holding the token
// can read these claims in a devtools tab. The signature proves the claims were
// not ALTERED, it does not hide them. Hence: identifiers only, never passwords
// or anything sensitive.
func GenerateToken(email string, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Convenience for the frontend (show "logged in as ..." without an extra
		// API call). Also a liability worth knowing: it is plaintext-readable,
		// and it goes stale - change your email and old tokens carry the old one
		// for up to 24h. Rule of thumb: identifiers in a JWT, not mutable
		// profile data.
		"email": email,
		// The claim that actually matters. It answers "who is making this
		// request?" for every handler downstream.
		"userID": userID,
		// The security-critical claim. "exp" is a REGISTERED claim name from
		// RFC 7519, so the library recognises and enforces it automatically.
		// Rename it to "expiry" and it becomes meaningless custom data that
		// nothing checks. .Unix() because the spec mandates seconds-since-epoch,
		// not an RFC3339 string.
		//
		// Why 24h? Since we cannot revoke, exp is the entire blast-radius
		// control, and it pulls two ways: 15min = small breach window but
		// constant re-logins; 30 days = nice UX but a stolen token is a
		// month-long backdoor. 24h is a defensible middle for this app. The
		// production answer is the refresh-token pattern: a ~15min access token
		// like this one plus a long-lived refresh token that IS stored
		// server-side and so CAN be revoked - you give up statelessness only on
		// the rare refresh call, not on every request.
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	// Returns an error rather than panicking: signing can fail (wrong key type,
	// or a future change of algorithm), and the caller should turn that into a
	// 500 rather than take the server down.
	return token.SignedString(secretKey)
}

// VerifyToken checks a token and returns the user it belongs to. This is the
// function that actually enforces security, and it runs on every protected
// request - one HMAC computation, no DB, no network. That hot path is the whole
// payoff of choosing JWT.
func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token,
		// The keyfunc. It is a callback rather than a plain key argument because
		// with asymmetric algorithms or key rotation you would inspect the token
		// header (its "kid") to decide WHICH key to verify against. We have one
		// key, so we ignore the argument.
		//
		// The classic bug here is trusting token.Header["alg"] to pick the key -
		// that is how the "alg: none" attack works (strip the signature, claim
		// no algorithm is needed) and how algorithm confusion works (send an
		// HS256 token signed with the server's PUBLIC RSA key, which the
		// attacker has). We never look at the header; the option below is what
		// defeats both.
		func(token *jwt.Token) (any, error) {
			return secretKey, nil
		},
		// The most important line in this file. Reject, before verification,
		// anything whose header algorithm is not exactly HS256. "none" is out,
		// an RS256 token is out. The token gets no say in how it is validated -
		// we decided that here, in code. Using .Alg() instead of the literal
		// "HS256" removes typo risk.
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		// v5 validates "exp" only IF PRESENT - by the spec, a token with no exp
		// is simply a non-expiring token, and it would pass. This makes the
		// claim mandatory: defense in depth, so that if someone later drops exp
		// from GenerateToken, verification fails loudly instead of silently
		// minting immortal tokens.
		jwt.WithExpirationRequired(),
	)

	// We deliberately swallow the library's err and return something vague:
	// "signature invalid" vs "expired" vs "malformed" would tell an attacker
	// which part of their forgery to fix. (Worth adding: log the real err
	// server-side, otherwise the actual cause is lost for debugging.)
	if err != nil {
		return 0, errors.New("cloud not parse token.")
	}

	// Redundant in v5 - a non-nil error and Valid == false always come together -
	// but harmless belt-and-braces, and it documents the intent.
	if !parsedToken.Valid {
		return 0, errors.New("Invalid token!")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("Invalid token claims.")
	}

	// float64, not int64, and this catches everyone out: the payload is JSON, and
	// MapClaims is map[string]any, so encoding/json decodes EVERY number into a
	// float64 when the target is any. Our int64 went out as the JSON number 7 and
	// comes back as float64(7). Asserting .(int64) would fail every single time.
	//
	// And note the comma-ok form. A bare assertion PANICS on failure, and this
	// payload is attacker-influenced, so a bare assertion is a remote crash
	// vector. The habit to keep: never bare-assert data that came from outside
	// the process. (The commented-out email line below is exactly the bare form.)
	//
	// Is float64 lossy? Not at this scale - it holds integers exactly up to 2^53.
	// If IDs ever become 64-bit snowflakes, store the ID as a string claim and
	// strconv.ParseInt it instead.
	v, ok := claims["userID"].(float64)
	if !ok {
		return 0, errors.New("missing user id")
	}

	// email := claims["email"].(string)
	userID := int64(v)

	// Returning just the int64 rather than the claims map keeps the surface
	// narrow: it is all middleware.Auth needs, and handlers cannot reach into raw
	// untrusted token data. If roles are needed later, return a small typed
	// struct rather than leaking jwt.MapClaims.
	return userID, nil
}
