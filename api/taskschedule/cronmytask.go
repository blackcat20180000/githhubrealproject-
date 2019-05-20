package taskschedule

import(
	_ "errors"
	"fmt"
	"os"
	"encoding/csv"
	"github.com/extrame/xls"
	// "github.com/robfig/cron"
	"io/ioutil"
	"log"
	"archive/zip"
	"net/http"
	"path/filepath"
	"io"
	"strings"
	// "regexp"
	"github.com/grokify/html-strip-tags-go"
    // "github.com/extrame/xls"
	"crypto/tls"
	"github.com/tealeg/xlsx"
	"sort"
	"regexp"
	"time"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"strconv"
)
func GenerateError(err string) bool {
	fmt.Println("ok")
	return true
}
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string
	currentTime := time.Now()
	year := currentTime.Year()
	_ =year
    r, err := zip.OpenReader(src)
    if err != nil {
        return filenames, err
    }
    defer r.Close()

    for _, f := range r.File {

        // Store filename/path for returning and using later on
        fpath := filepath.Join(dest, f.Name)

        // Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return filenames, fmt.Errorf("%s: illegal file path", fpath)
        }

        filenames = append(filenames, fpath)

        if f.FileInfo().IsDir() {
            // Make Folder
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        // Make File
        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return filenames, err
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return filenames, err
        }

        rc, err := f.Open()
        if err != nil {
            return filenames, err
        }

        _, err = io.Copy(outFile, rc)

        // Close the file without defer to close before next iteration of loop
        outFile.Close()
        rc.Close()

        if err != nil {
            return filenames, err
        }
    }
    return filenames, nil
}
func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
	   err:= rows.Scan(&count)
	   checkErr(err)
   }   
   return count
}
func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}
var monthpattern =[12]string{"Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"}
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
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
func findmonth(searchmonstr string) int {
	var searchmon int;
	searchmon=-1;
	var curmonthpatt string;
	for i :=0;i<len(monthpattern);i++ {
		var pos int;
		pos=-1;
		curmonthpatt=monthpattern[i]
		pos=strings.Index(searchmonstr,curmonthpatt)
		if(pos>=0) {
			searchmon=i+1;
			break;
		}
	}
	return searchmon
}
func before(value string, a string) string {
    // Get substring before a string.
    pos := strings.Index(value, a)
    if pos == -1 {
        return ""
    }
    return value[0:pos]
}

func after(value string, a string) string {
    // Get substring after a string.
    pos := strings.LastIndex(value, a)
    if pos == -1 {
        return ""
    }
    adjustedPos := pos + len(a)
    if adjustedPos >= len(value) {
        return ""
    }
    return value[adjustedPos:len(value)]
}
func finditem(a[]int,item int) int {
	for zz :=0;zz<len(a);zz++ {
		var positem int 
		positem=a[zz]
		if item==positem {
			return zz
		} 
	}
	return -1;
}
func Mainstr(link string,indus_index int){
	currentTime := time.Now()
	var year int
	var tablename string
	year= currentTime.Year()
	fmt.Println("year",year)
	// db :=dbconn()
	manu_industry :=[18]string{"Printing & Related Support Activities","Textile Mills","Computer & Electronic Products","Electrical Equipment, Appliances & Components","Fabricated Metal Products","Paper Products","Wood Products","Primary Metals","Chemical Products","Food, Beverage & Tobacco Products","Miscellaneous Manufacturing","Petroleum & Coal Products","Transportation Equipment","Machinery","Furniture & Related Products","Plastics & Rubber Products","Apparel, Leather & Allied Products","Nonmetallic Mineral Products"};
	nomanu_industry :=[18]string{"Real Estate, Rental & Leasing","Information","Transportation & Warehousing","Utilities","Arts, Entertainment & Recreation","Professional, Scientific & Technical Services","Construction","Health Care & Social Assistance","Management of Companies & Support Services","Wholesale Trade","Public Administration","Agriculture, Forestry, Fishing & Hunting","Accommodation & Food Services","Mining","Finance & Insurance","Retail Trade","Other Services","Educational Services"}
	strArray1 := [18]string{"Printing","Textile","Computer","Electrical","Fabricated","Paper","Wood","Primary","Chemical","Food","Miscellaneous","Petroleum","Transportation","Machinery","Furniture","Plastics","Apparel","Nonmetallic"}
	strArray2 := [18]string{"Real Estate","Information","Transportation","Utilities","Arts, Entertainment","Professional","Construction","Health Care","Management of Companies","Wholesale","Administration","Agriculture","Accommodation","Mining","Finance","Retail Trade","Other Services","Educational Services"}

	var manu_comment[] string
	var nonmanu_comment[] string
	var manu_index[] int
	var comment_manu_industid[] int
	var comment_nomanu_industid[] int
	var temp_menu_index[] int
	// var de_manu_index[] int 
	var zero_count int 
	zero_count=0
	var zero_cont int
	var nomanu_index[] int
	var temp_nomanu_index[] int 
	// var temp_nomenu_index[] int
	// var de_nomanu_index[] int
	var manu_mark[18] int 
	var nomanu_mark[18] int 
	var month int 
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	var date_str string
	var respstr string
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	respstr= string(robots);
	var matchstr,grostr string
	var pmistr string
	grostr=after(respstr,"<!-- Paragraph Three -->")
	grostr=before(grostr,"<!-- Respondent List Items -->")
	grostr=after(grostr,"order")
	prod_sta :=strings.Split(grostr,".")
	date_str=before(respstr,"Manufacturing ISM")
	date_str=strip.StripTags(date_str)
	manu_datstr :=strings.Split(date_str,"\n")
	fmt.Println(manu_datstr[len(manu_datstr)-1])
	date_str=manu_datstr[len(manu_datstr)-1]
	for yy :=0;yy<len(monthpattern);yy++ { 
	  var pos int
	  var mon string 
	  mon=monthpattern[yy] 
	  pos =strings.Index(date_str,mon)
	  if pos<0 {
		continue
	  } else {
		 month=yy+1
		 break
	  }

	}
	fmt.Println("month:",month)
	pmistr= after(respstr,"<!-- PMI Number -->")
	pmistr=before(pmistr,"<!-- Subtitles -->")
	pmistr=strip.StripTags(pmistr)
	pmistr_arr :=strings.Split(pmistr,";")
	var pmi_data string
	pmi_data=pmistr_arr[len(pmistr_arr)-1]
	re := regexp.MustCompile("[+-]?([0-9]*[.])?[0-9]+")
	ob :=re.FindAllString(pmi_data, -1)
	fmt.Println(ob)
	fmt.Println("pmi_data",pmi_data)
	matchstr=after(respstr,"<!-- Respondent List Items -->");
	matchstr=before(matchstr,"</ul>");
	stripes :=strip.StripTags(matchstr);
	comment :=strings.Split(stripes,"\n")

	for i :=1;i<len(comment)-1;i++ {
		comment[i]=strings.Replace(comment[i], "&#8221;", "", -1)
		comment[i]=strings.Replace(comment[i], "&#8220;", "", -1)
		comment[i]=strings.Replace(comment[i], "&#8212;", "", -1)
		comment[i]=strings.Replace(comment[i], "&amp; ", "", -1)
		comment[i]=strings.TrimSpace(comment[i])
		
	}
	if indus_index == 1 {
		var man_len int 
		man_len=len(comment)
		_=man_len
		for i :=0;i<len(comment);i++ {
	 		if comment[i] != "" {
				manu_comment=append(manu_comment,comment[i])
			}
			fmt.Println("manu_len",len(manu_comment))
		}
	} else {
			var noman_len int 
			noman_len=len(comment)
			for z:=0;z<noman_len;z++ {
				if comment[z] !="" {
					nonmanu_comment=append(nonmanu_comment,comment[z])
				}
			}
			 fmt.Println("no manu_comment length",len(nonmanu_comment))
		}
	if indus_index == 1{
		tablename="i_pmi_man"
		for mm :=0;mm<len(manu_comment);mm++ {
			var evcomment string
			evcomment=manu_comment[mm]
			for zz :=0;zz<len(strArray1);zz++ {
				var item string
				var pos int
				item=strArray1[zz] 
				pos =strings.Index(evcomment,item)
				if(pos<0) {
					continue
				} else {
					comment_manu_industid=append(comment_manu_industid,zz)
					break
				}
			}
		}
		fmt.Println("comment_manu_industid",comment_manu_industid)
	} else { 
		tablename="i_pmi_nonman"
		for yy :=0;yy<len(nonmanu_comment);yy++	{
			var evcomment string
			evcomment=nonmanu_comment[yy]
			for zz :=0;zz<len(strArray2);zz++ {
				
				var pos int 
				var item string 
				item=strArray2[zz]
				pos=strings.Index(evcomment,item)
				if pos<0 {
					continue
				} else {
					comment_nomanu_industid=append(comment_nomanu_industid,zz)
					break
				}

			}
		}
		fmt.Println("comment_nomanu_industid",comment_nomanu_industid)
	}
	if indus_index ==1{
		zero_count=0
		fmt.Println("prodsta0",prod_sta[0])
		for j :=0;j<len(strArray1); j++ {
			var temp string 
			var temp_index int
			temp=strArray1[j];
			temp_index=strings.Index(prod_sta[0],temp)
			fmt.Println("temp_index",temp_index)
			fmt.Println("\n")
			// fmt.Println("temp_index",temp_index)
			if temp_index >=0 {
				manu_index=append(manu_index,temp_index)
				temp_menu_index=append(temp_menu_index,temp_index)
			} else {

				temp_index=strings.Index(prod_sta[1],temp)
				// fmt.Println("----temp_index",temp_index)
				if temp_index < 0 {
					temp_index=0
				} else {
					zero_count=zero_count+1
				}
				temp_index=-1*temp_index
				manu_index=append(manu_index,temp_index)
				temp_menu_index=append(temp_menu_index,temp_index)
			}
		}

		sort.Ints(manu_index)
		fmt.Println("zero_count",zero_count)
		fmt.Println("temp_menu_index",temp_menu_index)
		fmt.Println("manu_index",manu_index)
		// sort.Ints(de_manu_index)
		for z :=0;z<len(manu_index);z++ {
			 var temp_val int
			 var pos int
			 temp_val=temp_menu_index[z]
			 pos=sort.IntSlice(manu_index).Search(temp_val)
			
			 if temp_val<0 {
				manu_mark[z]=-1*(zero_count-pos)
			 } else if temp_val == 0 {
				manu_mark[z]=0
			 } else {
				manu_mark[z]=len(manu_index)-pos
			}
		}
		fmt.Println("zero_count",zero_count)
			
		fmt.Println(manu_mark)
		fmt.Println("\n")	
		// fmt.Println(de_manu_index)
	} else {
		for j :=0;j<len(strArray2); j++ {
			var temp string 
			var temp_index int 
		    temp=strArray2[j];
			temp_index=strings.Index(prod_sta[0],temp)
			if temp_index >=0 {
				nomanu_index=append(nomanu_index,temp_index)
				temp_nomanu_index=append(temp_nomanu_index,temp_index)
				fmt.Println("temp_nomanu_index",temp_index)
			} else {
				temp_index=strings.Index(prod_sta[1],temp)
				if temp_index < 0 {
					temp_index=0
				} else {
					zero_cont=zero_cont+1;
				} 
				temp_index=temp_index*-1
				nomanu_index=append(nomanu_index,temp_index)
				temp_nomanu_index=append(temp_nomanu_index,temp_index)
				
			}
			
		}
		sort.Ints(nomanu_index)
	   fmt.Println("_________________________",prod_sta[0])
		fmt.Println("__________________________	")	
	   for m :=0;m<len(nomanu_index);m++{
		   var temp_val int 
		   var pos int   	
		   temp_val=temp_nomanu_index[m]
			pos=sort.IntSlice(nomanu_index).Search(temp_val)
			if temp_val<0 {
				nomanu_mark[m]=-1*(zero_cont-pos)
			} else if temp_val ==0 {
				nomanu_mark[m]=0
			} else {
				nomanu_mark[m]=len(nomanu_index)-pos
			}
		}
		fmt.Println(nomanu_index)
		fmt.Println("nomenupos:",nomanu_mark)
		//sql query 

	}
	var sqlquery string
	var dat string
	var temp string
	if month <10 {
		temp="0"+strconv.Itoa(month)
	} else {
		temp= strconv.Itoa(month)
	}
	fmt.Println(temp);
	dat=strconv.Itoa(year)+"-"+temp+"-01"
	db :=dbconn()
	var datestr string 
	datestr="select count(*) from "+tablename+" where dat =?"
	results,err:=db.Query(datestr,dat)
    var resultcount int
	if err != nil {
		panic(err.Error())
	}
	
	fmt.Println("results",results)
	resultcount=checkCount(results)
	 var insertquery string
	 if indus_index == 1{
		insertquery="insert into "+tablename+" (dat,manu_index) VALUES(?,?)"
	 } else {
		insertquery="insert into "+tablename+" (dat,non_manu) VALUES(?,?)"
	 }
	 fmt.Println("results",resultcount)
	if resultcount == 0 {
		stmtIns,err :=db.Prepare(insertquery)
		if err != nil {
			panic(err.Error())
		}
		defer stmtIns.Close()
		_, err = stmtIns.Exec(dat, ob[0])
		if err != nil {
			panic(err.Error())
		} 
	} else {
		var updatequery string
		if indus_index == 1{
			updatequery="update "+tablename+" set manu_index = ? where dat = ? and manu_index <> ?"
		} else {
			updatequery="update "+tablename+" set non_manu = ? where dat = ? and non_manu <> ?"
		}
	  stmt, err := db.Prepare(updatequery)
	  checkErr(err)
	  rest,err :=stmt.Exec(ob[0],dat,ob[0])
	  _ =rest
	  checkErr(err)
	  fmt.Println("rest",rest)
	  	
	}
	var selquery string
	// comment query
	var insertcomquery string
	var updatequery string   
	var currentindmarkrate int 
	var commentstring string 
	_ =insertcomquery
   if indus_index == 1 {
	var commentpos int
	   fmt.Println("commentids",comment_manu_industid)
		for zz :=0;zz<len(strArray1);zz++ {
			var ind string 
			commentpos=-1
			ind=strArray1[zz]
			selquery="select count(*) from i_pmi_man_industries where dat =? and industry like ?"
			_ =selquery
			ind="%"+ind+"%"
			results,err:=db.Query(selquery,dat,ind)
			checkErr(err)
			commentpos=finditem(comment_manu_industid,zz)
			fmt.Println("commentpos:",commentpos,"zz:",zz)
			var resultcount int 
			resultcount=checkCount(results)
			_=resultcount
			
			if commentpos >=0 {
				commentstring=manu_comment[commentpos]
			} else {
				commentstring=""
			}
			currentindmarkrate=manu_mark[zz]
			if resultcount == 0 {
				insertquery="insert into i_pmi_man_industries(dat,industry,rank,comment,section) values (?,?,?,?,?)"
				stmt,err :=db.Prepare(insertquery)
				checkErr(err)
				_,err=stmt.Exec(dat,manu_industry[zz],currentindmarkrate,commentstring,zz)
				checkErr(err)
			} else {
			    fmt.Println("updatecommentstring",commentstring)
				updatequery="update i_pmi_man_industries set rank=?,comment=? where section= ? and dat =?"
				stmtch,err :=db.Prepare(updatequery)
				checkErr(err)
				_,err=stmtch.Exec(currentindmarkrate,commentstring,zz,dat)
				checkErr(err)
			}
		}
		fmt.Println("resultcount",resultcount)
		fmt.Println("sqlquery",sqlquery)	
		defer results.Close()
		defer db.Close()
   } else {
		var indr string 
		var commentpos int
		var resultcount int
		var selquery string
		var  nomark_rate int
		var commnetstr string
		var updatequery string
		for yy :=0;yy<len(strArray2);yy++ {
			indr=strArray2[yy]
			indr="%"+indr+"%"
			selquery="select count(*) from i_pmi_nonman_industries where dat=? and industry like ?"
			results,err :=db.Query(selquery,dat,indr)
			nomark_rate=nomanu_mark[yy]
			checkErr(err)
			resultcount=checkCount(results)
			commentpos=finditem(comment_nomanu_industid,yy)
			
			if commentpos >=0 {
				commnetstr=nonmanu_comment[commentpos]

			} else {
				commnetstr=""
			}
			
			var insertquery string 
			if resultcount == 0 {
				insertquery="insert into i_pmi_nonman_industries(dat,industry,rank,comment,type) values(?,?,?,?,?)"
				 stmt,err :=db.Prepare(insertquery)
				 checkErr(err)
				 _,err=stmt.Exec(dat,nomanu_industry[yy],nomark_rate,commnetstr,yy)
				 checkErr(err)
			} else {
				updatequery="update i_pmi_nonman_industries set comment= ?,rank=? where type=? and dat=?"
				stmt,err :=db.Prepare(updatequery)
				checkErr(err)
				_,err=stmt.Exec(commnetstr,nomark_rate,yy,dat)
				checkErr(err)
			}
		}
		fmt.Println("nonmanu_comment",nonmanu_comment)
		fmt.Println("comment_nomanu_industid",comment_nomanu_industid)
   }
// fmt.Println(); result
}
func Insertumsci() {
	db :=dbconn()
	str,err :=readCSVFromUrl("http://www.sca.isr.umich.edu/files/tbmics.csv")
	if err != nil {
		checkErr(err)
	}
	rowcount :=len(str)
	lastrow :=str[rowcount-1]
	var monthstr string 
	var curmonth int
	monthstr=lastrow[0]
	fmt.Println("monthstr",monthstr)
	for m :=0;m<len(monthpattern);m++ {
		var pos int  
		var pomonth string 
		pomonth=monthpattern[m]
		
		pos=strings.Index(monthstr,pomonth)
	
		if pos>=0 {
			curmonth=m+1;
			break;
		} else {
			continue;
		}
	}
	var year int
	currentTime := time.Now()
    year=currentTime.Year()
	if strings.Index(monthstr,strconv.Itoa(year))<0 {
		year=year-1 
	}
	var tempmonth string;
	
	if curmonth<10 {
		tempmonth="0"+strconv.Itoa(curmonth)+"-01"
	} else {
		tempmonth=strconv.Itoa(curmonth)+"-01"
	}
	tempmonth=strconv.Itoa(year)+"-"+tempmonth
	fmt.Println("tempmonth",tempmonth)
	results,err :=db.Query("select count(*) from i_umcsi where period=?",tempmonth)
	checkErr(err)
	var chcount int 
	chcount=checkCount(results);
	fmt.Println("chcount",chcount)
	updatearr :=strings.Split(lastrow[0],",")
	fmt.Println("lastrow[2]",len(updatearr))
	if chcount == 1{
		stmtup,err :=db.Prepare("update i_umcsi set csi=? where period=?")
		checkErr(err)
		stmtup.Exec(updatearr[2],tempmonth)
	} else {
		stmtins,err :=db.Prepare("insert into i_umcsi(period,csi) values(?,?)")
		checkErr(err)
		
		_,err=stmtins.Exec(tempmonth,updatearr[2])
		checkErr(err)
	}
	fmt.Println(str)	
}

func savexlsumcsidata(excelfilename string){
	db :=dbconn()
	currentTime := time.Now() 
	
	var target1 int 
	var target2 int 
	_ = target2 
	var target3 int
	_ = target3 
	var flag int
	var month int 
	flag=0; 
	if xlFile, err := xls.Open("data/BuildPsdata.xls", "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			year := currentTime.Year()
			var tempyear=year-1;
			for i :=1;i<=int(sheet1.MaxRow);i++ {
				Row := sheet1.Row(i)
				if Row ==nil {
					flag=0
					continue;
				} else {
					if flag == 0 {
						col0 :=Row.Col(0)
						s := strconv.Itoa(tempyear)
						m :=strings.Contains(col0,s);
						if m == false {
							continue;
						}
						if m == true  {
							flag=1
							continue;
						}
                             
					}	
					if flag == 1 {
						temptarget,err :=strconv.Atoi(Row.Col(1))
					
						if err != nil {
							flag=0
							tempyear=tempyear+1
							continue
							
						}
						
					if temptarget > 200 {
							target1=temptarget
							fmt.Println("month",Row.Col(0))
							month=findmonth(Row.Col(0))
							fmt.Println("target1",target1)
						} else {
							flag=0
							tempyear=tempyear+1
						}
									
					}
					fmt.Println("target1",target1)		
					
					
				}
				fmt.Println("month",month) 
				
			}

		}
	   if sheet3 :=xlFile.GetSheet(2); sheet3 != nil {
			year := currentTime.Year()
			var tempyear=year-1;
			for i :=1;i<=int(sheet3.MaxRow);i++ {
				Row := sheet3.Row(i)
				
				if Row ==nil {
					flag=0
					continue;
				} else {
					if flag == 0 {
						col0 :=Row.Col(0)
						s := strconv.Itoa(tempyear)
						m :=strings.Contains(col0,s);
						if m == false {
							continue;
						}
						if m == true  {
							flag=1
							continue;
						}
					}	
					if flag == 1 {
						temptarget,err :=strconv.Atoi(Row.Col(1))
					
						if err != nil {
							flag=0
							tempyear=tempyear+1
							continue
							
						}
						if temptarget > 200 {
							target2=temptarget
						} else {
							flag=0
							tempyear=tempyear+1
						}
					}
					fmt.Println("target2",target2)		
				} 
				
			}
		}
	    if sheet5 :=xlFile.GetSheet(4); sheet5 != nil {
			year := currentTime.Year()
			var tempyear=year-1;
			for i :=1;i<=int(sheet5.MaxRow);i++ {
				Row := sheet5.Row(i)
				if Row ==nil {
					
					flag=0
					continue;
				} else {
					if flag == 0 {
						col0 :=Row.Col(0)
						// fmt.Println(col0)
						s := strconv.Itoa(tempyear)
						m :=strings.Contains(col0,s);
						if m == false {
							continue;
						}
						if m == true  {
							flag=1
							continue;
						}

					}	
					if flag == 1 {
						temptarget,err :=strconv.Atoi(Row.Col(1))
					
						if err != nil {
							flag=0
							tempyear=tempyear+1
							continue
							
						}
						
				    	if temptarget > 200 {
							target3=temptarget
							
						} else {
							flag=0
							tempyear=tempyear+1
						}
									
					}
			
					
				} 
				
			}
		}
	  
	}
	fmt.Println("target3",target3)
	var temp string;
	var count int;
	 if month<10 {
		 temp="0"+strconv.Itoa(month)
	 } else {
		temp=strconv.Itoa(month)
	 }
	temp="2019-"+temp+"-01"	 
	results,err :=db.Query("select count(*) from i_building_permits where period=?",temp)
	checkErr(err)
	count=checkCount(results)
	if count == 0 {
		query :="insert into i_building_permits values(?,?,?,?)"
		stmp,err :=db.Prepare(query)
		checkErr(err)
		_,err=stmp.Exec(temp,target1,target2,target3)
		checkErr(err)	
	} else {
		query :="update i_building_permits set permits=?,starts=?,completed=? where period=?"
		stmup,err:=db.Prepare(query)
		_,err=stmup.Exec(target1,target2,target3,temp)
		checkErr(err)
		fmt.Println(count)
		
	}
}
func DownloadFile(filepath string,  url string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	 }
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)
	resp, err := client.Do(req)
	checkErr(err)
	if err != nil {
		// handle err
	}
	out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}
func Insertbuildingshit(){
	
	fileUrl := "https://www.census.gov/construction/nrc/xls/newresconst.xls"
	excelFileName := "data/BuildPsdata.xls"
	if err := DownloadFile(excelFileName, fileUrl); err != nil {
        panic(err)
	}
	savexlsumcsidata(excelFileName);
}	
func Eu_data() {
	url :="https://www.census.gov/retail/marts/www/marts_current.xls"
	filepath  :="data/marts_current.xls"
   DownloadFile(filepath,url)
	db :=dbconn()
	_ =db 
	currentTime := time.Now() 
	year := currentTime.Year()
	curmonth :=currentTime.Month();
   _ = year
   var pos int;
   var month int ;
   var flag int ;
   var variable_arr[12] int;
	_ =variable_arr
   month = -1
	if xlFile, err := xls.Open(filepath, "utf-8"); err == nil {
	   if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
		   for i :=1;i<=int(sheet1.MaxRow);i++ {
			   Row := sheet1.Row(i)
			   if Row == nil {
				   continue
			   }
			   col1 :=Row.Col(1)
			   if flag == 1 {
				   col4 :=Row.Col(4)
					   item,err:=strconv.Atoi(col4)
					   if err != nil {
						   flag=1;
						   continue;
					   }
					   if item <1 {
						   flag=1;
						   continue;
					   }
					   flag=0
					   variable_arr[pos]=item
				   
			   }
			   if findtargetstr(col1) >=0 {
				   pos=findtargetstr(col1)
				   if variable_arr[pos] >0 {
					   continue;
				   } else {
					   col4 :=Row.Col(4)
					   item,err:=strconv.Atoi(col4)
					   if err != nil {
						   flag=1;
						   continue;
					   }
					   if item <1 {
						   flag=1;
						   continue;
					   }
					   flag=0
					   variable_arr[pos]=item
					   
				   }
				   
			   }
			   col2 :=Row.Col(4)
			   if month < 0 {
				   temp :=findmonth(col2)
				   if temp >0 {
					   month=temp	
					   fmt.Println("month",month)
				   } 
			   }
			   
		   }
	   }
   }
   if int(curmonth)+1 < month {
	   year=year-1
   }
   var temp string 
   if month<10 {
	   temp="0"+strconv.Itoa(month)
   } else {
	   temp=strconv.Itoa(month)
   }
   temp=strconv.Itoa(year)+"-"+temp+"-01"
   results,err :=db.Query("select count(*) from i_retail_sales where period=?",temp)
   checkErr(err)
   var count int;
   count=checkCount(results)
	   item1 :=variable_arr[0]
	   item2 :=variable_arr[1]
	   item3 :=variable_arr[2]+variable_arr[3]
	   item4 :=variable_arr[4]
	   item5 :=variable_arr[5]
	   item6 :=variable_arr[6]
	   item7 :=variable_arr[7]+variable_arr[8]
	   item8 :=variable_arr[9]
	   item9 :=variable_arr[10]
	   item10 :=variable_arr[11]
   if count == 0 {
	   stmt,err :=db.Prepare("insert into i_retail_sales values(?,?,?,?,?,?,?,?,?,?,?)")
	   checkErr(err)
	   _,err=stmt.Exec(temp,item1,item4,item9,item5,item7,item8,item10,item6,item2,item3)	
	   checkErr(err)
   } else {
	   stmt,err :=db.Prepare("update i_retail_sales set retail_sales_and_food_services=?,food_stores=?,non_store_retail=?,health=?,clothing_hobby=?,general_merch=?,food_services=?,gas_station=?,motors=?,domestic_products=? where period=?")
	   checkErr(err)
	   _,err=stmt.Exec(item1,item4,item9,item5,item7,item8,item10,item6,item2,item3,temp)
	   checkErr(err)
   }
   fmt.Println(variable_arr)
}
func Esi_data() { 
	currentTime := time.Now()
	db :=dbconn()
	filepath :="data/euesi.zip"
	url :="http://ec.europa.eu/economy_finance/db_indicators/surveys/documents/series/nace2_ecfin_1904/esi_nace2.zip"
	DownloadFile(filepath,url)
	files,err :=Unzip(filepath,"data")
	fmt.Println(files);
	checkErr(err)
	xlFile, err := xlsx.OpenFile(files[0])
	checkErr(err)
	sheet :=xlFile.Sheets
	fmt.Println("sheet:",sheet)
	sheet3 :=sheet[3]
	lens :=len(sheet3.Rows)
	fmt.Println("lastrow",lens);
	lastrow :=sheet3.Rows[lens-1]
	var month int;
	var cuountrypos int;
	var itempos int;
	var item1,item2,item3,item4,item5,item6 float64;
	var countryname string 
	cuountrypos=0;
	monthstr :=lastrow.Cells[0].String()
	curmonth :=findmonth(monthstr)
	country_arr :=[31]string{"EU","","","BE","BG","CZ","DK","DE","EE","IE","EL","ES","FR","HR","IT","CY","LV","LT","LU","HU","MT","NL","AT","PL","PT","RO","SI","SK","FI","SE","UK"}
	year := currentTime.Year()
	month=int(currentTime.Month())
	if curmonth > month {
		year=year-1
	}
	var temp string
	if curmonth < 10 {
		temp="0"+strconv.Itoa(curmonth)
	} else {
		temp=strconv.Itoa(curmonth)
	}
	temp=strconv.Itoa(year)+"-"+temp+"-01"
	for i :=2;i<len(lastrow.Cells);i++ {
		cell :=lastrow.Cells[i]
		text := cell.String()
		countryname=country_arr[cuountrypos]
		fmt.Println("")
		var item float64;
		itemtemp,err:=strconv.ParseFloat(text, 64);
		if err != nil {
				
		} else {
			item=itemtemp
		}
		
		if itempos ==6 {
			
			cuountrypos=cuountrypos+1
			itempos=0 
			continue;
		}
		itempos=itempos+1;
		if cuountrypos == 1 || cuountrypos == 2 {
			continue
		}
		if itempos==1{
			item1=item
		} else if itempos==2 {
			item2=item
		} else if itempos ==3 {
			item3=item
		} else if itempos== 4{
			item4=item
		} else if itempos ==5 {
			item5=item
			
		} else {
			item6=item
			results,err :=db.Query("select count(*) from i_esi where period=? and country_code = ?",temp,countryname)
			checkErr(err)
			var count int 
			count=checkCount(results)
			if count ==0 {
				stmt,err :=db.Prepare("insert into i_esi values(?,?,?,?,?,?,?,?)")
				checkErr(err)
				_,err=stmt.Exec(temp,countryname,item1,item2,item3,item4,item5,item6)
				checkErr(err)
			} else {
				stmup,err :=db.Prepare("update i_esi set industry=?,services=?,consumer=?,retail=?,building=?,esi=? where period=? and country_code=?")
				checkErr(err)
				_,err=stmup.Exec(item1,item2,item3,item4,item5,item6,temp,countryname)
				checkErr(err)
			}
		}
		fmt.Printf("%s\n", text)
	}
}
func findtargetstr(str string) int {
	var targetpattern=[12]string{"Retail","Motor","Furniture","Electronics","Food &","Health","Gasoline","Clothing ","Sporting","merchandise","Nonstore","Food services"}
	var pos int;
	fmt.Println(str)
	for z :=0;z<len(targetpattern);z++ {
		var targetstr string;
		targetstr=targetpattern[z]
		pos=strings.Index(str,targetstr)
		if pos < 0{
			continue;
		} else {
			fmt.Println("pos",z)
			return z;
		}
	}
	return -1;
}