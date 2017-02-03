package main

import (
	"database/sql"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// We use DB as temporary storage
	// for testing only
	os.Remove("./main.db")

	// Init DB
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = InitTables(db)
	if err != nil {
		log.Fatal(err)
	}

	// Init Echo
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t
	e.Static("/assets", "assets")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "main", "World")
	})

	e.GET("/GetEmployeeList", func(c echo.Context) error {
		e, err := GetEmployees(db)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, e)
	})

	e.GET("/GetEmployeeRoles", func(c echo.Context) error {
		employeeId := c.QueryParam("employeeId")
		id, err := strconv.Atoi(employeeId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		e, err := GetEmployeeRoles(id, db)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, e)
	})

	e.POST("/ClockIn", func(c echo.Context) error {
		employeeId, err := strconv.Atoi(c.FormValue("employeeId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		roleId, err := strconv.Atoi(c.FormValue("roleId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		e := Attendance{
			EmployeeId: employeeId,
			RoleId:     roleId,
			ActionTime: int(time.Now().Unix()),
		}
		err = SaveAttendance(e, db)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, e)
	})

	e.Logger.Fatal(e.Start(":80"))
}
