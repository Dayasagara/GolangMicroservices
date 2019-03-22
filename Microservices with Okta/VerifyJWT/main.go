package verifyjwt

import(
	"net/http"
	"github.com/gorilla/sessions"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"fmt"
)

type Jwt struct {
	Claims map[string]interface{}
}

var sessionStore = sessions.NewCookieStore([]byte("okta-custom-login-session-store"))
var state = "ApplicationState"
var nonce = "NonceNotSetYet"

func VerifyHandler(r *http.Request) (*verifier.Jwt, error) {

	session, err := sessionStore.Get(r, "okta-custom-login-session-store")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(session.Values["access_token"])
	var tok string
	tok = (session.Values["id_token"]).(string)
	fmt.Println("tok",tok)
	result, verificationError := verifyToken(tok) //Token verification
 
	if verificationError != nil {
		fmt.Println("verf err",verificationError)
	}

	fmt.Println("result:jwt:",result)
	return result,verificationError
}

func verifyToken(t string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = 
	tv["aud"] = 
	fmt.Println("tv",tv)
	jv := verifier.JwtVerifier{
		Issuer:    ,
		ClaimsToValidate: tv,
	}
	fmt.Println("jv:",jv)
	result, err := jv.New().VerifyIdToken(t)

	if err != nil {
		return nil, fmt.Errorf("err:%s", err)
	}

	if result != nil {
		fmt.Println("res",result)
		fmt.Println(result.Claims)
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified: %s", "")
}