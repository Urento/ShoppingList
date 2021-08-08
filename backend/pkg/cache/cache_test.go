package cache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	. "github.com/stretchr/testify/assert"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestCacheJWTToken(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	exists, err := EmailExists(email)

	Equal(t, nil, err)
	Equal(t, true, exists)
}

func TestGetTokenByEmail(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	val, err := GetJWTByEmail(email)
	if err != nil {
		t.Errorf("Error while getting Token by Email %s", err)
	}

	Equal(t, nil, err)
	Equal(t, token, val)
}

func TestDoesTokenExpireAfter1Day(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	ttl, err := GetTTLByEmail(email)
	if err != nil {
		t.Errorf("Error getting the ttl from the key by email %s", err)
	}

	if ttl < 86200 {
		t.Errorf("ttl is too low")
	}

	Equal(t, nil, err)
}

func TestGetEmailByJWT(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT Token %s", err)
	}

	val, err := GetEmailByJWT(token)
	if err != nil {
		t.Errorf("Error while getting Email by Token: %s", err)
	}

	Equal(t, email, val)
}

func TestDeleteToken(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT Token %s", err)
	}

	ok, err := DeleteTokenByEmail(email, token)
	if err != nil || !ok {
		t.Errorf("Error while deleting Token by Email: %s", err)
	}

	_, err = GetJWTByEmail(email)
	if err == nil {
		t.Error("Token still cached")
	}

	_, err = GetEmailByJWT(token)
	if err == nil {
		t.Error("Token is still cached")
	}

	Equal(t, true, ok)
	Equal(t, "jwt token not cached", err.Error())
}

func TestDeleteTokenWithEmailThatDoesntExist(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	ok, err := DeleteTokenByEmail(email, token)
	if err == nil || ok {
		t.Errorf("No Error thrown 4")
	}

	_, err = GetJWTByEmail(email)
	if err == nil {
		t.Errorf("No Error thrown 3")
	}

	_, err = GetEmailByJWT(token)
	if err == nil {
		t.Errorf("No Error thrown 2 ")
	}

	Equal(t, false, ok)
	Equal(t, "jwt token not cached", err.Error())
}

func TestIsTokenValidWithValidToken(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT Token %s", err)
	}

	valid, err := IsTokenValid(token)
	if err != nil {
		t.Errorf("Error while checking if token is valid %s", err)
	}

	Equal(t, true, valid)
	Equal(t, nil, err)
}

func TestIsTokenValidWithInvalidToken(t *testing.T) {
	Setup()

	token := StringWithCharset(245)

	valid, _ := IsTokenValid(token)

	Equal(t, false, valid)
}

func TestCacheUser(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	password := StringWithCharset(30)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()
	ip := RandomIPAddress()

	pwdHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while creating password hash: %s", err)
	}

	u := User{
		EMail:                   email,
		Username:                username,
		Password:                pwdHash,
		EmailVerified:           emailVerified,
		Rank:                    rank,
		TwoFactorAuthentication: twoFactorAuthentication,
		IPAddress:               ip,
	}

	err = u.CacheUser()
	if err != nil {
		t.Errorf("Error while caching user %s", err)
	}

	Equal(t, nil, err)
}

func TestGetUser(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	password := StringWithCharset(30)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()
	ip := RandomIPAddress()

	pwdHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while creating password hash: %s", err)
	}

	u := User{
		EMail:                   email,
		Username:                username,
		Password:                pwdHash,
		EmailVerified:           emailVerified,
		Rank:                    rank,
		TwoFactorAuthentication: twoFactorAuthentication,
		IPAddress:               ip,
	}

	err = u.CacheUser()
	if err != nil {
		t.Errorf("Error while caching user %s", err)
	}
	t.Log(email)
	t.Log(u.EMail)

	user, err := GetUser(email)
	if err != nil {
		t.Errorf("Error while getting user: %s", err)
	}

	Equal(t, nil, err)
	Equal(t, email, user.EMail)
	Equal(t, pwdHash, user.Password)
	Equal(t, emailVerified, user.EmailVerified)
	Equal(t, username, user.Username)
	Equal(t, rank, user.Rank)
	Equal(t, twoFactorAuthentication, user.TwoFactorAuthentication)
}

func TestGetUserThatDoesntExist(t *testing.T) {
	Setup()

	_, err := GetUser("dkjfgbksdjhfgbkjdhfsgb@gmail.com")

	NotEqual(t, nil, err)
}

func TestUpdateUser(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	password := StringWithCharset(30)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()
	ip := RandomIPAddress()

	pwdHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while creating password hash: %s", err)
	}

	u := User{
		EMail:                   email,
		Username:                username,
		Password:                pwdHash,
		EmailVerified:           emailVerified,
		Rank:                    rank,
		TwoFactorAuthentication: twoFactorAuthentication,
		IPAddress:               ip,
	}

	err = u.CacheUser()
	if err != nil {
		t.Errorf("Error while caching user %s", err)
	}

	user, err := GetUser(email)
	if err != nil {
		t.Errorf("Error while getting user: %s", err)
	}

	newUsername := StringWithCharset(10)
	newEmailVerified := RandomBoolean()
	newRank := RandomRank()
	newTwoFactorAuthentication := RandomBoolean()
	newIp := RandomIPAddress()
	newUser := User{
		EMail:                   email,
		Username:                newUsername,
		Password:                pwdHash,
		EmailVerified:           newEmailVerified,
		Rank:                    newRank,
		TwoFactorAuthentication: newTwoFactorAuthentication,
		IPAddress:               newIp,
	}

	err = UpdateUser(newUser)
	if err != nil {
		t.Errorf("Error while updating user: %s", err)
	}

	updatedUser, err := GetUser(email)
	if err != nil {
		t.Errorf("Error while getting updated user: %s", err)
	}

	Equal(t, nil, err)
	Equal(t, email, user.EMail)
	Equal(t, pwdHash, user.Password)
	Equal(t, emailVerified, user.EmailVerified)
	Equal(t, username, user.Username)
	Equal(t, rank, user.Rank)
	Equal(t, twoFactorAuthentication, user.TwoFactorAuthentication)
	Equal(t, email, updatedUser.EMail)
	Equal(t, newUsername, updatedUser.Username)
	Equal(t, pwdHash, updatedUser.Password)
	Equal(t, newEmailVerified, updatedUser.EmailVerified)
	Equal(t, newRank, updatedUser.Rank)
	Equal(t, newTwoFactorAuthentication, updatedUser.TwoFactorAuthentication)
}

func TestDeleteUser(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	password := StringWithCharset(30)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()
	ip := RandomIPAddress()

	pwdHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while creating password hash: %s", err)
	}

	u := User{
		EMail:                   email,
		Username:                username,
		Password:                pwdHash,
		EmailVerified:           emailVerified,
		Rank:                    rank,
		TwoFactorAuthentication: twoFactorAuthentication,
		IPAddress:               ip,
	}

	err = u.CacheUser()
	if err != nil {
		t.Errorf("Error while caching user: %s", err)
	}

	err = DeleteUser(email)
	if err != nil {
		t.Errorf("Error while deleting user: %s", err)
	}

	_, shouldErr := GetUser(email)
	if shouldErr == nil {
		t.Errorf("GetUser didn't throw an error after deleting")
	}

	Equal(t, nil, err)
	NotEqual(t, nil, shouldErr)
}

func TestGenerateSecretIdAndVerify(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	secretId, err := GenerateSecretId(email)
	if err != nil {
		t.Errorf("Error while generating secret id: %s", err)
	}

	ok, err := VerifySecretId(email, secretId)
	if err != nil {
		t.Errorf("Error while verifying secert id: %s", err)
	}

	Equal(t, true, ok)
	Equal(t, nil, err)
}

func TestVerifySecretIdWithWrongIdWithExistingAccount(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	_, err := GenerateSecretId(email)
	if err != nil {
		t.Errorf("Error while generating secret id: %s", err)
	}

	ok, err := VerifySecretId(email, "secretId")
	if err != nil {
		t.Errorf("Error while verifying secert id: %s", err)
	}

	Equal(t, false, ok)
	Equal(t, nil, err)
}

func TestVerifySecretIdWithWrongIdWithoutAccount(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	ok, err := VerifySecretId(email, "secretId")
	if err == nil {
		t.Errorf("No error thrown even though the secretId is wrong and doesn't exist")
	}

	Equal(t, false, ok)
	Equal(t, "secretid is not valid", err.Error())
}

func TestGetTwoFactorAuthenticationStatus(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	username := StringWithCharset(10)
	password := StringWithCharset(30)
	emailVerified := RandomBoolean()
	rank := RandomRank()
	twoFactorAuthentication := RandomBoolean()
	ip := RandomIPAddress()

	pwdHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		t.Errorf("Error while creating password hash: %s", err)
	}

	u := User{
		EMail:                   email,
		Username:                username,
		Password:                pwdHash,
		EmailVerified:           emailVerified,
		Rank:                    rank,
		TwoFactorAuthentication: twoFactorAuthentication,
		IPAddress:               ip,
	}

	err = u.CacheUser()
	if err != nil {
		t.Errorf("Error while caching user: %s", err)
	}

	status, err := GetTwoFactorAuthenticationStatus(email)
	if err != nil {
		t.Errorf("Error while getting two factor authentication status: %s", err)
	}

	Equal(t, twoFactorAuthentication, status)
	Equal(t, nil, err)
}

func TestHasSecretId(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	secretId, err := GenerateSecretId(email)
	if err != nil {
		t.Errorf("Error while generating secret id: %s", err)
	}

	ok, err := VerifySecretId(email, secretId)
	if err != nil {
		t.Errorf("Error while verifying secert id: %s", err)
	}

	key, has, err := HasSecretId(email)
	if err != nil {
		t.Errorf("Error while checking if the user still has a secretId: %s", err)
	}

	if !has {
		t.Errorf("User does not have a secretId even though he has one")
	}

	if key != secretId {
		t.Errorf("SecretId is not the same as the previously generated one")
	}

	Equal(t, true, ok)
	Equal(t, nil, err)
	Equal(t, secretId, key)
}

func TestHasSecretIdWhenTheUserDoesntHaveOne(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	_, has, _ := HasSecretId(email)

	if has {
		t.Errorf("User doesnt have a SecretId but it says it has")
	}

	Equal(t, false, has)
}

func TestInvalidateSecretId(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	secretId, err := GenerateSecretId(email)
	if err != nil {
		t.Errorf("Error while generating secret id: %s", err)
	}

	ok, err := VerifySecretId(email, secretId)
	if err != nil {
		t.Errorf("Error while verifying secert id: %s", err)
	}

	key, has, err := HasSecretId(email)
	if err != nil {
		t.Errorf("Error while checking if the user still has a secretId: %s", err)
	}

	if !has {
		t.Errorf("User does not have a secretId even though he has one")
	}

	if key != secretId {
		t.Errorf("SecretId is not the same as the previously generated one")
	}

	err = InvalidateSecretId(email)
	if err != nil {
		t.Errorf("Error while invalidating secretId: %s", err)
	}

	_, has2, err := HasSecretId(email)
	if err != nil {
		t.Errorf("Error while checking if the user still has a secretId 2: %s", err)
	}

	if has2 {
		t.Errorf("SecretId did not get invalidated!")
	}

	Equal(t, true, ok)
	Equal(t, nil, err)
	Equal(t, secretId, key)
}

func TestInvalidateAllJWTTokens(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	exists, err := EmailExists(email)
	if err != nil {
		t.Errorf("Error while checking if the email exists: %s", err)
	}

	secretId, err := GenerateSecretId(email)
	if err != nil {
		t.Errorf("Error while generating secret id: %s", err)
	}

	ok, err := VerifySecretId(email, secretId)
	if err != nil {
		t.Errorf("Error while verifying secert id: %s", err)
	}

	err = InvalidateAllJWTTokens(email)
	if err != nil {
		t.Errorf("Error while invalidating jwt tokens: %s", err)
	}

	ok2, _ := VerifySecretId(email, secretId)

	if ok2 {
		t.Errorf("SecretId did not get invalided!")
	}

	exists2, err2 := EmailExists(email)
	if err2 != nil {
		t.Errorf("Error while checking if the email exists: %s", err)
	}

	valid, err := IsTokenValid(token)

	if valid {
		t.Errorf("JWT Token did not get invalided!")
	}

	Equal(t, nil, err)
	Equal(t, true, exists)
	Equal(t, false, exists2)
	Equal(t, true, ok)
	Equal(t, false, ok2)
}

func TestDoesTokenBelongToEmail(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	exists, err := EmailExists(email)
	if err != nil {
		t.Errorf("Error while checking if the email exists: %s", err)
	}

	ok, err := DoesTokenBelongToEmail(email, token)
	if err != nil {
		t.Errorf("Error while checking if the token belongs to the email: %s", err)
	}

	Equal(t, true, exists)
	Equal(t, true, ok)
}

func TestInvalidateSpecificJWTToken(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	token := StringWithCharset(245)

	err := CacheJWT(email, token)
	if err != nil {
		t.Errorf("Error while caching JWT token %s", err)
	}

	exists, err := EmailExists(email)
	if err != nil {
		t.Errorf("Error while checking if the email exists: %s", err)
	}

	ok, err := DoesTokenBelongToEmail(email, token)
	if err != nil {
		t.Errorf("Error while checking if the token belongs to the email: %s", err)
	}

	err = InvalidateSpecificJWTToken(email, token)
	if err != nil {
		t.Errorf("Error while invalidating specific jwt token: %s", err)
	}

	ok2, _ := DoesTokenBelongToEmail(email, token)

	Equal(t, true, exists)
	Equal(t, true, ok)
	Equal(t, false, ok2)
}

func TestCacheAndGetTOTPSecret(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	secret := StringWithCharset(16)

	err := CacheTOTPSecret(email, secret)
	if err != nil {
		t.Errorf("Error while caching TOTP Secret: %s", err)
	}

	totpSecret, err := GetTOTPSecret(email)
	if err != nil {
		t.Errorf("Error while getting TOTP Secret: %s", err)
	}

	Equal(t, secret, totpSecret)
	Equal(t, nil, err)
}

func TestGetTOTPSecretThatDoesntExist(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	_, err := GetTOTPSecret(email)

	Equal(t, "totp secret is not cached", err.Error())
}

func TestDeleteTOTPSecret(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	secret := StringWithCharset(16)

	err := CacheTOTPSecret(email, secret)
	if err != nil {
		t.Errorf("Error while caching TOTP Secret: %s", err)
	}

	totpSecret, err := GetTOTPSecret(email)
	if err != nil {
		t.Errorf("Error while getting TOTP Secret: %s", err)
	}

	err = DeleteTOTPSecret(email)
	if err != nil {
		t.Errorf("Error while deleting TOTP Secret: %s", err)
	}

	_, delErr := GetTOTPSecret(email)

	Equal(t, secret, totpSecret)
	Equal(t, nil, err)
	Equal(t, "totp secret is not cached", delErr.Error())
}

func TestIsTOTPCached(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"
	secret := StringWithCharset(16)

	err := CacheTOTPSecret(email, secret)
	if err != nil {
		t.Errorf("Error while caching TOTP Secret: %s", err)
	}

	totpSecret, err := GetTOTPSecret(email)
	if err != nil {
		t.Errorf("Error while getting TOTP Secret: %s", err)
	}

	ok, err := IsTOTPSecretCached(email)
	if err != nil {
		t.Errorf("Error while checking if TOTP Secret is cached: %s", err)
	}

	Equal(t, secret, totpSecret)
	Equal(t, true, ok)
}

func TestIsTOTPCachedWhenItsNotCached(t *testing.T) {
	Setup()

	email := StringWithCharset(10) + "@gmail.com"

	_, err := GetTOTPSecret(email)

	ok, _ := IsTOTPSecretCached(email)

	Equal(t, false, ok)
	Equal(t, "totp secret is not cached", err.Error())
}

func StringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomBoolean() bool {
	nmb := rand.Intn(2)
	return nmb <= 1
}

func RandomRank() string {
	nmb := rand.Intn(2)
	if nmb <= 1 {
		return "admin"
	}
	return "default"
}

func RandomIPAddress() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
