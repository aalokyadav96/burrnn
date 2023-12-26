package main

import (
    "fmt"
    "net/http"
    "log"
	"io"
	"os"
	"time"
	"strings"
	"unicode"
	"os/exec"
	"crypto/rand"
	"encoding/json"
	rndm "math/rand"
	"crypto/md5"
	"path/filepath"

    "github.com/julienschmidt/httprouter"
)

// Gif - We will be using this Gif type to perform crud operations
type GIF struct {
	Title  string
	Author string
	Hashags   []string
	Date   string
	URL    string
	Views  int
	Likes  int
}


type UsrSett struct {
	Name string
	Email string
	Password string
	SecretKey string
	Phone string
	PublicKey string
	PrivateKey string
}
	
func ViewUserSettings() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tmpl.ExecuteTemplate(w, "user.html", UsrSett{Name: ps.ByName("name"), Email : "a@a.com", Password : "pswd", SecretKey : "seecret", Phone : "5645435454", PublicKey : "fgr6wd8rt46378jtg674w89drg64d73gr76s34", PrivateKey : "fg4sd8th643s79tdr463ftrg64d3hnt7"})
       //~ fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
    }
}

func ManageContent() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tmpl.ExecuteTemplate(w, "manage.html", nil)
    }
}
	
type Usr struct {
	Name string
	Username string
	About string
	UserPic string
	Vids []string
}
	
func ViewUser() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		vids, _ := rdxGet("uz"+ps.ByName("name"))
		tg := strings.Split(vids,":::")
		tmpl.ExecuteTemplate(w, "user.html", Usr{Name: ps.ByName("name"), Username : ps.ByName("name"), About : "About to add", Vids : tg, UserPic : "userpic"})
       //~ fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
    }
}

func Me() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			//~ rdxHset("title",fileName,title)
		cookie, err := r.Cookie("exampleCookie")
		if err != nil {
			fmt.Fprintf(w,"MeMe")
		}
		//~ fmt.Fprintf(w,cookie.Value)
		vids, _ := rdxGet("uz"+cookie.Value)
		tg := strings.Split(vids,":::")
		tmpl.ExecuteTemplate(w, "meprofile.html", Usr{Name: cookie.Value, Username : cookie.Value, About : "About to add", Vids : tg, UserPic : "5645435454"})
	}
}

func MeProfile() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			//~ rdxHset("title",fileName,title)
		cookie, err := r.Cookie("exampleCookie")
		if err != nil {
			fmt.Fprintf(w,"MeMe")
		}
		//~ fmt.Fprintf(w,cookie.Value)
		vids, _ := rdxGet("uz"+cookie.Value)
		tg := strings.Split(vids,":::")
		tmpl.ExecuteTemplate(w, "meprofile.html", Usr{Name: cookie.Value, Username : cookie.Value, About : "About to add", Vids : tg, UserPic : "5645435454"})
	}
}



func UploadVideoFileHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "POST" {           
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "index.html", nil)
	} else {
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
		}

		var fileEndings string
		var fileName string

		files := r.MultipartForm.File["vidfile"]
		for _, fileHeader := range files {
			log.Println("hao")
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer file.Close()
			log.Println("file OK")
		    cookie, _ := r.Cookie("exampleCookie")
			cook := cookie.Value
			title := r.FormValue("title")
			tags := strings.ToLower(r.FormValue("tags"))
			fmt.Println(tags)
			// Get and print outfile size
			
			f := func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
			}
			titleArr := strings.FieldsFunc(strings.ToLower(title), f)
			fmt.Printf("Fields are: %q", titleArr)
			fileSize := fileHeader.Size
			fmt.Printf("File size (bytes): %v\n", fileSize)
			// validate file size
			if fileSize > maxUploadSize {
				renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			}
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				renderError(w, "INVALID_FILE"+http.DetectContentType(fileBytes), http.StatusBadRequest)
			}

			// check file type, detectcontenttype only needs the first 512 bytes
			detectedFileType := http.DetectContentType(fileBytes)
			switch detectedFileType {
			case "video/mp4":
				fileEndings = ".mp4"
				break
			case "video/webm":
				fileEndings = ".webm"
				break
			case "image/gif":
				fileEndings = ".gif"
				break
			default:
				renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			}
			fileName = GenerateName(rndmToken(12))
			// if fileName exists in Redis, again GenerateName(rndmToken(12))
			//		fileEndings, err := mime.ExtensionsByType(detectedFileType)

			if err != nil {
				renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			}
			newFileName := fileName + fileEndings

			newPath := filepath.Join(uploadPath, newFileName)
			fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

			// write file
			newFile, err := os.Create(newPath)
			if err != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			}
			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			}
			//			ffmpeg -i %1 -filter:v scale=-2:640:flags=lanczos -c:a copy -pix_fmt yuv420p %1_lcz.mp4
			FFConvert(fileName , fileEndings )
			FFPoster(fileName , fileEndings )
	//			var gif Gif{Title: , Author: , Tags: , Date: , URL: , Views: , Likes: }
			t := time.Now()
			tf := string(t.Format("2006-01-02"))
			//~ var newstr string = title+":::"+tags+":::"+cook+":::"+string(t.Format("2006-01-02"))
			rdxHset("title",fileName,title)
			rdxHset("tags",fileName,tags)
			rdxHset("date",fileName,tf)
			rdxHset("author",fileName,cook)
			rdxAppend("uz"+cook,":::"+fileName)
			fmt.Println(titleArr)/*
			for _,v := range titleArr {
				rdxAppend("tags"+v,fileName+":::")
			}*/
			NoSQLPostTagged(fileName, tags, tf)
			for _,v := range strings.Split(tags,",") {
				rdxAppend("tags"+v,fileName+":::")
			}
			NoSQLPostCreated(fileName,cook, title, tags, tf)
			http.Redirect(w, r, "/v/"+fileName, http.StatusSeeOther)
			}
			
//				tmpl.ExecuteTemplate(w, "show.html", fileName)
	}
}

func ViewPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == "GET" {
		fmt.Println(r.URL.Path)
		filenmae := ps.ByName("PostId")
		//~ res,_ := rdxHget("gif",ps.ByName("PostId"))
		a,_ := rdxHget("title",filenmae)
		b,_ := rdxHget("tags",filenmae)
		c,_ := rdxHget("date",filenmae)
		d,_ := rdxHget("author",filenmae)
		//~ fmt.Println(res)
		//~ arr := strings.Split(res,":::")
		//~ tg := strings.Split(arr[1],",")
		tg := strings.Split(b,",")
		//~ var gif = GIF{Title: arr[0], Tags: tg, Author: arr[2], Date: arr[3], URL: ps.ByName("PostId")}
		var gif = GIF{Title: a, Hashags: tg, Author: d, Date: c, URL: filenmae}
		fmt.Println(gif)
		tmpl.ExecuteTemplate(w, "show.html", gif)
	}
}

func ViewAllFiles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "viewall.html", SearchFiles(uploadPath))
	}
}

func MeCollectionsGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method=="GET" {
		var jsonresp string = `audio_gifs,reactions,likes`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w,jsonresp)
	}
}

func MeCollectionsPut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method=="PUT" {
		var jsonresp string = "200"
		fmt.Fprintf(w,jsonresp)
	}
}

func MeCollectionsPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method=="POST" {
		var jsonresp string = "200"
		fmt.Fprintf(w,jsonresp)
	}
}

func MeCollectionGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method=="GET" {
		var jsonresp string = "200"
		fmt.Fprintf(w,jsonresp)
	}
}

func MeCollectionPut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method=="PUT" {
		var jsonresp string = r.FormValue("collection")+","+r.FormValue("gifname")
		cookie, err := r.Cookie("exampleCookie")
		if err != nil {
			fmt.Fprintf(w,"MeMe")
		}
		NoSQLAddToCollection(r.FormValue("gifname"),cookie.Value,r.FormValue("collection")) 
		fmt.Fprintf(w,jsonresp)
	}
}

func MeCollectionPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method=="POST" {
		var jsonresp string = "200"
		fmt.Fprintf(w,jsonresp)
	}
}

func DeleteFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println(r.URL.Path)
	if r.Method == "POST" {
		fileName := ps.ByName("PostId")
		extn := fileName[:len(fileName)-len(filepath.Ext(fileName))]
		fmt.Println("./uploads/" + fileName + ".webm")
		os.Remove("./uploads/" + fileName + ".webm")
		os.Remove("./uploads/" + fileName + ".mp4")
		fmt.Println(streamPath + "/" + extn + ".mp4")
		os.Remove(streamPath + "/" + extn + ".mp4")
		fmt.Println("./posters/" + extn + ".png")
		os.Remove("./posters/" + extn + ".png")
		rdxHdel("title",ps.ByName("PostId"))
		rdxHdel("tags",ps.ByName("PostId"))
		rdxHdel("date",ps.ByName("PostId"))
		rdxHdel("author",ps.ByName("PostId"))
		XHRrespond(w, "Deleted")
	}
}

type Serp struct {
	Query string
	Poster []string
}

func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query().Get("q")
	var resArray []string
	//~ f := func(c rune) bool {
		//~ return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	//~ }
	//~ for _,v := range strings.FieldsFunc(strings.ToLower(query), f) {
	for _,v := range strings.Split(query, " ") {
		fmt.Println(v)
		redisRes,_ := rdxGet("tags"+v)
		fmt.Println(redisRes)
		res := strings.Split(redisRes,":::")
		fmt.Println(res)
		for _,k := range res {
			resArray = append(resArray,k)
		}
		fmt.Println(resArray)
	}
	log.Println(resArray)
	tmpl.ExecuteTemplate(w, "searchresults.html", Serp{Query:query,Poster:resArray[:len(resArray)-1]})
}

type Posters struct {
	Qtag string
	Poster []string
}

func SearchFiles(dir string) []string {
	var allFiles []string
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		switch file.Name() {
		case "$RECYCLE.BIN", "$Recycle.Bin":
			break
		case "System Volume Information":
			break
		default:
			allFiles = append(allFiles, file.Name())
		}
	}
	return allFiles
}

func GenerateName(w int64) string {
	rndm.Seed(time.Now().Unix()) // initialize global pseudo random generator
	p1 := fmt.Sprintf(adjectives[rndm.Intn(len(adjectives))])
	p2 := fmt.Sprintf(adjectives[rndm.Intn(len(adjectives))])
	p3 := fmt.Sprintf(animals[rndm.Intn(len(animals))])
	return p1 + p2 + p3
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func rndmToken(len int) int64 {
	b := make([]byte, len)
	n, _ := rand.Read(b)
	return int64(n)
}

func XHRrespond(w http.ResponseWriter, message string) {
	fmt.Fprintf(w, message)
}

func EncrypIt(strToHash string) string {
	data := []byte(strToHash)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func SessionVerify(sessionKey string) string {
	return fmt.Sprintf(sessionKey)
}

func Ignore(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "favicon.png")
}


func FFConvert(fileName string, fileEndings string) {
	getFrom := uploadPath + "/" + fileName + fileEndings
	saveAs := streamPath + "/" + fileName + ".mp4"
	cmd := exec.Command("ffmpeg", "-i", getFrom, "-filter:v", "scale=-2:640:flags=lanczos", "-c:a", "copy", "-pix_fmt", "yuv420p", saveAs)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}


func FFPoster(fileName string, fileEndings string) {
	getFrom := streamPath + "/" + fileName + ".mp4"
	saveAs :=  "posters/" + fileName + ".png"
    out, err := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-count_frames", "-show_entries", "stream=nb_read_frames", "-print_format", "default=nokey=1:noprint_wrappers=1",getFrom).Output()
    if err != nil {
        log.Println(err)
    }
    fmt.Printf(string(out))
//	cmd := exec.Command("ffmpeg", "-i", getFrom, "-vf", "scale=-2:640:flags=lanczos", saveAs)
	cmd := exec.Command("ffmpeg", "-i", getFrom, "-vf", "scale=-2:480:flags=lanczos", saveAs)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}



func UploadImage(w http.ResponseWriter, r *http.Request, folder string, formName string) string {
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
		}
		var fileEndings string
		var fileName string
		var newFileName string
		
		files := r.MultipartForm.File[formName]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer file.Close()
			fileSize := fileHeader.Size
			fmt.Printf("File size (bytes): %v\n", fileSize)
			if fileSize > maxUploadSize {
				renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			}
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				renderError(w, "INVALID_FILE"+http.DetectContentType(fileBytes), http.StatusBadRequest)
			}
			detectedFileType := http.DetectContentType(fileBytes)
			switch detectedFileType {
			case "image/png":
				fileEndings = ".png"
				break
			case "image/webp":
				fileEndings = ".webp"
				break
			case "image/jpg":
				fileEndings = ".jpg"
				break
			case "image/jpeg":
				fileEndings = ".jpeg"
				break
			default:
				renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			}
			fileName = GenerateFileName(16)
			// if fileName exists in Redis, again GenerateName(rndmToken(12))
			//		fileEndings, err := mime.ExtensionsByType(detectedFileType)

			if err != nil {
				renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			}
			newFileName = fileName + fileEndings
			fmt.Println("nfn :::: ",newFileName)
			newPath := filepath.Join(folder, newFileName)
			fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

			// write file
			newFile, err := os.Create(newPath)
			if err != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			}
			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			}
		}
			fmt.Println("nfn :::: ",newFileName)
		return newFileName
}


func GeneratePostId(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZ")

    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rndm.Intn(len(letters))]
    }
    return string(b)
}


func GenerateFileName(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZ")

    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rndm.Intn(len(letters))]
    }
    return string(b)
}

func init() {
    rndm.Seed(time.Now().UnixNano())
}


func sendImageAsHTML(w http.ResponseWriter, r *http.Request, a string) {
	fmt.Fprintf(w,a)
}

func sendImageAsAttachment(w http.ResponseWriter, r *http.Request) {
    buf, err := os.ReadFile("F46ZPQ0bQAACFYs.jpg")
    if err != nil {
        log.Fatal(err)
    }
    w.Header().Set("Content-Type", "image/jpeg")
    w.Header().Set("Content-Disposition", `attachment;filename="F46ZPQ0bQAACFYs.jpg"`)
    w.Write(buf)
}

func sendImageAsBytes(w http.ResponseWriter, r *http.Request, a string) {
    buf, err := os.ReadFile("./uploads/"+a)
    if err != nil {
        log.Fatal(err)
    }
    w.Header().Set("Content-Type", "image/jpeg")
    w.Write(buf)
}

func DeleteEdit(filename string) {
	time.Sleep(10 * time.Second)  
	os.Remove(filename)
	fmt.Println(filename,"Deleted")
}



type Account struct {
	Username string
	Password string
	Email string
}

func CheckUname(username string) int {
	_,err := rdxHget("uzer", username)
	if err!=nil {
		return 404
	}
	return 200
}

func Register() {//POST
	var a Account
	//~ if 404 == CheckUname(r.FormValue("username")) {} else {}
	a.Username =  "newusername"
	a.Password = "passwordispassword"
	a.Email = "usera@user.com"
	rdxHset("uzer", a.Username, a.Password)
	rdxHset("uzeremail", a.Username, a.Email)
	CreateProfile(a.Username, a.Email)
	//~ fmt.Println(a)	
}

func CreateProfile(username, email string) {
	var profilestring string = `{
  "username": {
    "email": "eamail@email.com",
    "name": "Display name",
    "about": "About needs to be edited",
    "followercount": "30",
    "followingcount": "30",
    "phone": "65478376483",
    "userid": "45rf654fe554ge74",
    "accountid": "4372184378146738",
    "publickey": "thfht7f7htf5h7t5fht5",
    "privatekey": "4gj78fvfv7f54fj",
    "password": "y4jyjyj7y7jyl",
    "secret": "y4jyj4y447jy78",
    "guardian": "8f584poyf8gro",
    "dob": "12-12-1212",
    "other_emails": "a@a.com,b@b.com",
    "address": "somewhere",
    "work": "knowhere llc",
    "devices": "phone,laptop",
    "fingerprint": "available",
    "iris": "NA"
  }
}`
	rdxHset("uprofile", ":username",profilestring)
}



func EditMyProfile(username string) {
	_,_ = rdxHget("uprofile", ":username")
	//~ fmt.Println(res)
}

func UserProfile() {//GET
	_,err := rdxHget("uprofile", ":username")
	if err != nil {
		return
	}
	//~ fmt.Println(uprof)
}


type Video struct {
	Author string
	Title string
	URL string
	Tags string
	Time string
}

func Upload() string {//POST
	var u Video
	u.Author =  "newusername"
	u.Title = "Title that can have spaces"
	u.URL = GenerateFileName(12)
	u.Tags = "atab,batang,tag with space"
	u.Time = (time.Now()).Format("202386")
	rdxHset("uploadauthor", u.URL, u.Author)
	rdxHset("uploadtitle", u.URL, u.Title)
	rdxHset("uploadtags", u.URL, u.Tags)
	rdxHset("uploadtime", u.URL, u.Time)
	NoSQLPostCreated(u.URL,u.Author, u.Title, u.Tags, u.Time)
	//~ fmt.Println(u)
	return u.URL
}
	

func Edit(filename string) {//PUT
	if !FileExists(filename) {
		//~ fmt.Println("Post no longer exist")
		return
	}
	var u Video
	u.Title = "Title that was modified"
	u.URL = filename
	u.Tags = "njf gbfd hg"
	u.Time = (time.Now()).Format("202386")
	rdxHset("uploadtitle", u.URL, u.Title)
	rdxHset("uploadtags", u.URL, u.Tags)
	rdxHset("uploadtime", u.URL, u.Time)
	//~ fmt.Println(u)
}

func View(filename string) {//GET
	var v Video
	v.Author , _ = rdxHget("uploadauthor", filename)
	v.Title , _ = rdxHget("uploadtitle", filename)
	v.Tags , _ = rdxHget("uploadtags", filename)
	v.Time , _ = rdxHget("uploadtime", filename)
	//~ fmt.Println(v)
}

func Delete(filename string) {//DELETE
	rdxHdel("uploadauthor", filename)
	rdxHdel("uploadtitle", filename)
	rdxHdel("uploadtags", filename)
	rdxHdel("uploadtime", filename)
	//~ fmt.Println("deleted")
	View(filename)
	Edit(filename)
}

func Save(filename string) {//PUT
	postid := filename
	userid := "newusername"
	collectionid := "gifts"
	NoSQLAddToCollection(postid,userid,collectionid)
}

func Like(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method=="PUT" {
		fmt.Println("Like : ", r.FormValue("gifname"))
		cookie, err := r.Cookie("exampleCookie")
		if err != nil {
			fmt.Fprintf(w,"404")
		}
		postid := r.FormValue("gifname")
		userid := cookie.Value
		collectionid := "likes"
		NoSQLAddToCollection(postid,userid,collectionid)
		fmt.Fprintf(w,"200")
	}
}

func Report(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method=="PUT" {
		fmt.Println("Like : ", r.FormValue("gifname"))
		cookie, err := r.Cookie("exampleCookie")
		if err != nil {
			fmt.Fprintf(w,"MeMe")
		}
		postid := r.FormValue("gifname")
		userid := cookie.Value
		reason := r.FormValue("reason")
		description := r.FormValue("description")
		email := r.FormValue("email")
		NoSQLPostReported(postid,userid,reason,email, description)
	}
}


type Collectn struct {
	Operation string `json: "ops"`
	PostId string `json: "postid"`
	Collection string `json: "collection"`
	Username string `json: "username"`
}

var collector map[string]Collectn

func NoSQLAddToCollection(postid,userid,collectionid string) {
    collector := make(map[string]Collectn)
	var klc Collectn
	klc.Operation = "a2c"
	klc.PostId = postid 
	klc.Collection = collectionid
	klc.Username = userid
	collector["op"] = klc
	//~ jsonstr,_ := json.MarshalIndent(collector,"","    ")
	//~ fmt.Println(string(jsonstr))
	//~ SendJSONToPy(string(jsonstr))
}



type Reportd struct {
	Operation string `json: "ops"`
	PostId string `json: "postid"`
	Username string `json: "username"`
	Reason string `json: "reason"`
	Email string `json: "email"`
	Description string `json: "description"`
}

var reporter map[string]Reportd

func NoSQLPostReported(postid,userid,reason,email, description string) {
    reporter := make(map[string]Reportd)
	var klc Reportd
	klc.Operation = "rpt"
	klc.PostId = postid
	klc.Username = userid
	klc.Reason = reason
	klc.Email = email
	klc.Description = description
	reporter["op"] = klc
	//~ fmt.Println("repo", reporter[postid])
	//~ jsonstr,_ := json.MarshalIndent(reporter,"","    ")
	//~ fmt.Println(string(jsonstr))
	//~ SendJSONToPy(string(jsonstr))
	//~ fmt.Println("rted")
}


type Created struct {
	Operation string `json: "ops"`
	PostId string `json: "postid"`
	Username string `json: "username"`
	Title string `json: "title"`
	Tags string `json: "tags"`
	Date string `json: "date"`
}

var created map[string]Created

func NoSQLPostCreated(postid,username, title, tags, date string) { //~ {"user1":{"created":["gif7","gif8","gif9"]}}
    created := make(map[string]Created)
	var klc Created
	klc.Operation = "crt"
	klc.PostId = postid
	klc.Username = username
	klc.Title = title
	klc.Tags = tags
	klc.Date = date
	created["op"] = klc
	//~ fmt.Println("cree", created[postid])
	//~ jsonstr,_ := json.MarshalIndent(created,"","    ")
	//~ jsonstr,_ := json.Marshal(created)
	//~ fmt.Println(string(jsonstr))
	//~ SendJSONToPy(string(jsonstr))
}

func GetTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res,_ := rdxGet("tags"+ps.ByName("tag"))
	arr := strings.Split(res,":::")
	tmpl.ExecuteTemplate(w, "tags.html", Posters{Poster:arr[:len(arr)-1], Qtag:ps.ByName("tag")})
}


type Tags struct {
	Operation string `json: "ops"`
	TagId string `json: "tagid"`
	PostId string `json: "postid"`
	Date string `json: "date"`
}

var tagged map[string]Tags

func NoSQLPostTagged(postid, tag, date string) { 
    tagged := make(map[string]Tags)
	var klc Tags
	klc.Operation = "tag"
	klc.TagId = tag
	klc.PostId = postid
	klc.Date = date
	tagged["op"] = klc
	//~ jsonstr,_ := json.MarshalIndent(tagged,"","    ")
	//~ jsonstr,_ := json.Marshal(tagged)
	//~ fmt.Println(string(jsonstr))
	//~ SendJSONToPy(string(jsonstr))
}


func FileExists(filename string) bool {
	_ , err := rdxHget("uploadauthor", filename)
	if err != nil {
		return false
	}
	return true
}



//~ {"user1":{"collections":["coll1","coll2","coll3"]}}
//~ {"user1":{"coll1":["gif1","gif2","gif3"],"coll2":["gif4","gif5","gif6"]}}
//~ {"user1":{"likes":["gif4","gif5","gif6"]}}
//~ {"user1":{"created":["gif7","gif8","gif9"]}}

//~ func UserExists() bool {return false}
