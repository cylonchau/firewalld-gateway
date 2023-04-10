package auther

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/cylonchau/firewalld-gateway/config"
)

const UserIDKey = "userID"

var jsonSecret = []byte("jwt")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return jsonSecret, nil
}

// jwt包自带的jwt.StandardClaims只包含了官方字段，若需要额外记录其他字段，就可以自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中

type Token struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func GenToken(userID int64) (string, error) {
	// 创建一个我们自己的声明的数据
	c := Token{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(), // 过期时间
			Issuer:    config.CONFIG.AppName,                   // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(jsonSecret)
}

func GetInfo(token string) (int64, error) {
	var enconterError error
	var tokenOk *jwt.Token
	tokenOk, enconterError = jwt.ParseWithClaims(token, &Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jsonSecret), nil
	})
	if enconterError == nil {
		if claims, ok := tokenOk.Claims.(*Token); ok && tokenOk.Valid {
			return claims.UserID, nil
		}
	}
	return 0, enconterError
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*Token, error) {
	// 解析token
	var mc = new(Token)
	token, err := jwt.ParseWithClaims(tokenString, mc, keyFunc)
	if err != nil {
		return nil, err
	}
	// 校验token
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

// RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token无效直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	// 从旧access token中解析出claims数据
	var claims Token
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	// 当 access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		token, _ := GenToken(claims.UserID)
		return token, "", nil
	}
	return
}
