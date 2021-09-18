package jwt

import (
	"net"
	"testing"
)

/*func TestJWT(t *testing.T) {
	cache.Setup()
	ip := GetOutboundIP(t)

	r := gin.Default()
	r.Use(JWT())
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"OPTIONS", "PUT", "GET", "POST", "DELETE", "PATCH"},
		AllowOrigins:     []string{"http://" + ip.String() + ":9998"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           9 * time.Hour,
		ExposeHeaders:    []string{"Content-Length"},
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	go r.Run(":9998")

	// create cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Errorf("Error while creating cookie jar: %v", err)
	}

	//prepare request
	c := &http.Client{Jar: jar}
	req, err := http.NewRequest("GET", "http://"+ip.String()+":9998", nil)
	if err != nil {
		t.Errorf("An error occurred while making the requests: %s", err)
	}

	// set cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "some_token",
		HttpOnly: true,
		MaxAge:   24 * 60 * 60,
		Domain:   "localhost",
		Path:     "/",
		Secure:   false,
	}
	req.AddCookie(cookie)

	// make request
	resp, err := c.Do(req)
	if err != nil {
		t.Errorf("Error while sending the request: %v", err)
	}
	defer resp.Body.Close()

	//check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status Code wasn't 200")
	}

	//check body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error while reading body: %v", err)
	}
	bodyString := string(bodyBytes)
	t.Log(bodyString)

	//TODO: NOT FINISHED: COOKIE NOT BEING SENT

	t.Errorf("dflkjgn")
}*/

func GetOutboundIP(t *testing.T) net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
