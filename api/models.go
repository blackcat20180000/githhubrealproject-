package api
import (
	_"database/sql"
	_ "fmt"
	_"github.com/go-sql-driver/mysql"
	_ "github.com/graphql-go/graphql"
	_ "io"
	_ "realpro/api/errors"
	_ "strconv"
 	_ "time"
)
type IPmiManIndustry struct {
	Dat      string `json:"dat"`
	Industry string `json:"industry"`
	Rank     int    `json:"rank"`
	Comment  string  `json:"comment"`
	Section  int    `json:"section"`
}
