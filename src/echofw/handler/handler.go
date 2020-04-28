package handler

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/ajay340/SearchBreaches.me/database"
	"github.com/ajay340/SearchBreaches.me/pdf"
	"github.com/ajay340/SearchBreaches.me/recommender"

	"github.com/labstack/echo"
)

func noCacheContext(c echo.Context) {
	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")
}

func cookieExpire() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "auth-token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(0 * time.Hour)
	return cookie
}

func IndexHandler(c echo.Context) error {
	noCacheContext(c)
	cookie, err := c.Cookie("auth-token")
	if (err != nil) || !database.IsValidSession(cookie.Value) {
		cookie = cookieExpire()
		c.SetCookie(cookie)
		return c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"username": database.GetSessionUsername(cookie.Value),
		})
	}
}

func LoginGetHandler(c echo.Context) error {
	noCacheContext(c)
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"Error_MSG": "",
	})
}

func UserGetHandler(c echo.Context) error {
	noCacheContext(c)
	cookie, err := c.Cookie("auth-token")
	if (err != nil) || !database.IsValidSession(cookie.Value) {
		cookie = cookieExpire()
		c.SetCookie(cookie)
		return c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		username := database.GetSessionUsername(cookie.Value)
		return c.Render(http.StatusOK, "User.html", map[string]interface{}{
			"username": username,
		})
	}
}

func UserPostHandler(c echo.Context) error {
	noCacheContext(c)
	password := database.GenerateHash(c.FormValue("password"))
	newpassword1 := c.FormValue("newpassword1")
	newpassword2 := c.FormValue("newpassword2")
	cookie, err := c.Cookie("auth-token")
	if (err != nil) || !database.IsValidSession(cookie.Value) {
		cookie = cookieExpire()
		c.SetCookie(cookie)
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}
	username := database.GetSessionUsername(cookie.Value)
	user := database.ReadUserInfo(username)
	if newpassword1 != newpassword2 {
		return c.Render(http.StatusOK, "User.html", map[string]interface{}{
			"Error_MSG": "Passwords do not match",
		})
	} else if password != user.Password {
		return c.Render(http.StatusOK, "User.html", map[string]interface{}{
			"Error_MSG": "Error: Password is incorrect",
		})
	} else {
		database.UpdateUser(user.Username, database.GenerateHash(newpassword1))
		return c.Render(http.StatusOK, "User.html", map[string]interface{}{
			"Error_MSG": "Password has successfully changed",
		})
	}
}

func LoginPostHandler(c echo.Context) error {
	noCacheContext(c)
	username := c.FormValue("username")
	password := database.GenerateHash(c.FormValue("password"))
	user := database.ReadUserInfo(username)
	if user.Password == "" {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error_MSG": "Error: User \"" + username + "\" does not exist",
		})
	} else if password != user.Password {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error_MSG": "Error: Password is incorrect",
		})
	} else {
		cookie := new(http.Cookie)
		cookie.Name = "auth-token"
		cookie.Value = database.GenerateToken(username)
		cookie.Expires = time.Now().Add(72 * time.Hour)
		database.AddSessionKey(cookie.Value, username)
		c.SetCookie(cookie)
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
}

func LogoutGetHandler(c echo.Context) error {
	noCacheContext(c)
	cookie, _ := c.Cookie("auth-token")
	database.DeleteSessionKey(cookie.Value)
	cookie = cookieExpire()
	c.SetCookie(cookie)
	return c.Redirect(http.StatusMovedPermanently, "/login")
}

func RegisterGetHandler(c echo.Context) error {
	noCacheContext(c)
	return c.Render(http.StatusOK, "register.html", map[string]interface{}{
		"Error_MSG": "",
	})
}

func AboutHandler(c echo.Context) error {
	noCacheContext(c)
	return c.Render(http.StatusOK, "about.html", map[string]interface{}{})
}

func SearchQueryHandler(c echo.Context) error {
	noCacheContext(c)
	cookie, err := c.Cookie("auth-token")
	if (err != nil) || !database.IsValidSession(cookie.Value) {
		cookie = cookieExpire()
		c.SetCookie(cookie)
		type Error struct {
			ErrorMSG string `json:"error"`
		}
		return c.JSON(http.StatusUnauthorized, Error{ErrorMSG: "Unauthorized"})
	} else {
		search := c.Param("search")
		breaches := database.SearchDB(search)
		return c.JSON(http.StatusOK, breaches)
	}

}

func RegisterPostHandler(c echo.Context) error {
	noCacheContext(c)
	username := c.FormValue("username")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("password1")
	if password != confirmPassword {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"Error_MSG": "Error: Passwords do not match",
		})
	} else if database.UserExists(username) {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"Error_MSG": "Error: \"" + username + "\" already exists",
		})
	} else {
		password := database.GenerateHash(password)
		database.AddUser(username, password)
		return c.Render(http.StatusOK, "registersuccess.html", map[string]interface{}{})
	}
}

func ListingHandler(c echo.Context) error {
	noCacheContext(c)
	cookie, err := c.Cookie("auth-token")
	if (err != nil) || !database.IsValidSession(cookie.Value) {
		cookie = cookieExpire()
		c.SetCookie(cookie)
		return c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		searchTerm := c.Param("search")
		isSafe := regexp.MustCompile(`^[A-Za-z ]+$`).MatchString
		if !isSafe(searchTerm) {
			return c.File("ErrorPages/400.html")
		} else {
			var searchResults []database.Breach = database.SearchDB(searchTerm)
			if len(searchResults) == 0 {
				return c.Render(http.StatusOK, "SearchNotFound.html", searchResults)
			} else {
				return c.Render(http.StatusOK, "listings.html", map[string]interface{}{})
			}
		}
	}
}

func PDFDownloaderHandler(c echo.Context) error {
	noCacheContext(c)
	cookie, err := c.Cookie("auth-token")
	if (err != nil) || !database.IsValidSession(cookie.Value) {
		cookie = cookieExpire()
		c.SetCookie(cookie)
		return c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		id := c.Param("id")
		Breach := database.FindRowBreach(id, "ID")
		if len(Breach.Name_of_Covered_Entity) <= 0 {
			return c.File("ErrorPages/404.html")
		} else {
			if len(Breach.Summary) <= 0 {
				Breach.Summary = "Breach of " + Breach.Name_of_Covered_Entity + " by " + Breach.Type_of_Breach
			}

			pdf, _ := pdf.PDFgen(Breach)
			return c.Blob(http.StatusOK, "application/pdf", pdf.Bytes())
		}
	}
}

func ProductHandler(c echo.Context) error {
	noCacheContext(c)
	cookie, err := c.Cookie("auth-token")
	if (err != nil) || !database.IsValidSession(cookie.Value) {
		cookie = cookieExpire()
		c.SetCookie(cookie)
		return c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		id := c.Param("id")
		Breach := database.FindRowBreach(id, "ID")
		if len(Breach.Name_of_Covered_Entity) <= 0 {
			return c.File("ErrorPages/400.html")
		} else {
			if len(Breach.Summary) <= 0 {
				Breach.Summary = "Breach of " + Breach.Name_of_Covered_Entity + " by " + Breach.Type_of_Breach
			}
			recommendations := recommender.RecommendTopBreach(id)
			recommendedBreach1 := database.FindRowBreach(strconv.Itoa(recommendations[0]), "ID")
			recommendedBreach2 := database.FindRowBreach(strconv.Itoa(recommendations[1]), "ID")
			recommendedBreach3 := database.FindRowBreach(strconv.Itoa(recommendations[2]), "ID")
			recommendedBreach4 := database.FindRowBreach(strconv.Itoa(recommendations[3]), "ID")
			return c.Render(http.StatusOK, "product.html", map[string]interface{}{
				"Breach_ID":                               Breach.ID,
				"Breach_Name_of_Covered_Entity":           Breach.Name_of_Covered_Entity,
				"Breach_State":                            Breach.State,
				"Breach_Business_Associate_Involved":      Breach.Business_Associate_Involved,
				"Breach_Individuals_Affected":             Breach.Individuals_Affected,
				"Breach_Date_of_Breach":                   Breach.Date_of_Breach,
				"Breach_Type_of_Breach":                   Breach.Type_of_Breach,
				"Breach_Location_of_Breached_Information": Breach.Location_of_Breached_Information,
				"Breach_Date_Posted_or_Updated":           Breach.Date_Posted_or_Updated,
				"Breach_Summary":                          Breach.Summary,
				"Breach_breach_start":                     Breach.Breach_start,
				"Breach_breach_end":                       Breach.Breach_end,
				"Breach_year":                             Breach.Year,
				"Breach_Industry":                         Breach.Industry,
				"recommend_1_id":                          recommendedBreach1.ID,
				"recommend_1_Name_of_Covered_Entity":      recommendedBreach1.Name_of_Covered_Entity,
				"recommend_2_id":                          recommendedBreach2.ID,
				"recommend_2_Name_of_Covered_Entity":      recommendedBreach2.Name_of_Covered_Entity,
				"recommend_3_id":                          recommendedBreach3.ID,
				"recommend_3_Name_of_Covered_Entity":      recommendedBreach3.Name_of_Covered_Entity,
				"recommend_4_id":                          recommendedBreach4.ID,
				"recommend_4_Name_of_Covered_Entity":      recommendedBreach4.Name_of_Covered_Entity,
			})
		}
	}
}
