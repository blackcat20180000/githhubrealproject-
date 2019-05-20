package realpro

import (
	"context"
	"realpro/api"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"realpro/errors"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.
func dbconn() (db *sql.DB) {
	dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
	dbName := "pmi"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
        panic(err.Error())
    }
    return db
}
type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) IPmiMans(ctx context.Context, limit *int, offset *int) ([]IPmiMan, error) {
	db :=dbconn()
	var ipmi IPmiMan
	var ipmis []IPmiMan
	rows,err :=db.Query("select * from i_pmi_man order by dat")
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ipmi.Dat, &ipmi.ManuIndex); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		ipmis = append(ipmis, ipmi)
	}
	return ipmis,nil
}
func (r *queryResolver) IPmiManIndustrys(ctx context.Context, limit *int, offset *int) ([]api.IPmiManIndustry, error) {
	db :=dbconn()
	var ipmivariable api.IPmiManIndustry
	var ipmivariables []api.IPmiManIndustry
	rows,err :=db.Query("select * from i_pmi_man_industries order by dat desc,section")
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ipmivariable.Dat,&ipmivariable.Industry,&ipmivariable.Rank,&ipmivariable.Comment,&ipmivariable.Section); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		ipmivariables = append(ipmivariables, ipmivariable)
	}
	return ipmivariables,nil
}
func (r *queryResolver) IPmiNomans(ctx context.Context, limit *int, offset *int) ([]IPmiNoman, error) {
	db :=dbconn()
	var ipminnoman IPmiNoman
	var ipminnomans []IPmiNoman
	rows,err :=db.Query("select * from i_pmi_nonman order by dat")
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ipminnoman.Dat,&ipminnoman.NonManu); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		ipminnomans = append(ipminnomans, ipminnoman)
	}
	return ipminnomans,nil
}
func (r *queryResolver) IPmiNonmanIndustries(ctx context.Context, limit *int, offset *int) ([]IPmiNonmanIndustry, error) {
	db :=dbconn()
	var ipminomanvariable IPmiNonmanIndustry
	var ipminomanvariables []IPmiNonmanIndustry
	rows,err :=db.Query("select * from i_pmi_nonman_industries order by dat desc,type")
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&ipminomanvariable.Dat,&ipminomanvariable.Industry,&ipminomanvariable.Rank,&ipminomanvariable.Comment,&ipminomanvariable.Type); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		ipminomanvariables = append(ipminomanvariables, ipminomanvariable)
	}
	return ipminomanvariables,nil
}
func (r *queryResolver) IRetailSales(ctx context.Context, limit *int, offset *int) ([]IRetailSale, error) {
	db :=dbconn()
	var iretailsale IRetailSale
	var iretailsales []IRetailSale
	rows,err :=db.Query("select * from i_retail_sales")
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&iretailsale.Period,&iretailsale.RetailSalesAndFoodServices,&iretailsale.FoodStores,&iretailsale.NonStoreRetail,&iretailsale.Health,&iretailsale.ClothingHobby,&iretailsale.GeneralMerch,&iretailsale.FoodServices,&iretailsale.GasStation,&iretailsale.Motors,&iretailsale.DomesticProducts); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		iretailsales = append(iretailsales, iretailsale)
	}
	return iretailsales,nil
}
func (r *queryResolver) IUmcsis(ctx context.Context, limit *int, offset *int) ([]IUmcsi, error) {
	db :=dbconn()
	var ucmis IUmcsi
	var ucmisdata []IUmcsi
	rows,err :=db.Query("select * from i_umcsi")
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {	
		if err := rows.Scan(&ucmis.Period,&ucmis.Csi); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		ucmisdata = append(ucmisdata, ucmis)
		
	}
	return ucmisdata,nil
}
func (r *queryResolver) IEsis(ctx context.Context, limit *int, offset *int) ([]IEsi, error) {
	db :=dbconn()
	var iesi IEsi
	var iesis []IEsi
	rows,err :=db.Query("select * from i_esi order by period asc")
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {	
		if err := rows.Scan(&iesi.Period,&iesi.CountryCode,&iesi.Industry,&iesi.Services,&iesi.Consumer,&iesi.Retail,&iesi.Building,&iesi.Esi); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		iesis = append(iesis,iesi)
		
	}
	return iesis,nil
}
func (r *queryResolver) IBuildingPermits(ctx context.Context, limit *int, offset *int) ([]IBuildingPermit, error) {
	db :=dbconn()
	var ibuildpermit IBuildingPermit
	var ibuildpermits []IBuildingPermit
	rows,err :=db.Query("select * from i_building_permits order by period")
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	defer rows.Close()
	for rows.Next() {	
		if err := rows.Scan(&ibuildpermit.Period,&ibuildpermit.Permits,&ibuildpermit.Starts,&ibuildpermit.Completed); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		ibuildpermits = append(ibuildpermits,ibuildpermit)
		
	}
	return ibuildpermits,nil
}
