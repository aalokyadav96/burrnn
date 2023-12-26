package main

import (
    "net/http"
    "log"
	"html/template"

    "github.com/julienschmidt/httprouter"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

const maxUploadSize = 10 * 1024 * 1024 // 10 mb
const uploadPath = "./uploads"
const streamPath = "./cmp"
var postersDir = "./posters"
var userpicPath = "./userpic"

func main() {
    router := httprouter.New()
    router.GET("/", HasAuthCookie(Index))

	router.GET("/search", Search)
	router.GET("/tag/:tag", GetTags)

	router.GET("/login", loginHandler)
    router.POST("/login", loginHandler)
    router.GET("/register", registerHandler)
    router.POST("/register", registerHandler)
    router.POST("/logout", HasAuthCookie(logoutHandler))

    router.GET("/@:name", HasAuthCookie(ViewUser()))
    router.GET("/me", HasAuthCookie(Me()))
    router.GET("/me/edit", HasAuthCookie(MeProfile()))
    router.GET("/manage", HasAuthCookie(ManageContent()))

	//~ GFYs
	router.GET("/new/video", HasAuthCookie(Index))
	router.POST("/new/video", HasAuthCookie(UploadVideoFileHandler))
	router.GET("/v/:PostId", ViewPost)
	router.POST("/del/:PostId", HasAuthCookie(DeleteFile))
	router.PUT("/like", HasAuthCookie(Like))
	router.PUT("/report", HasAuthCookie(Report))
	router.GET("/viewall", HasAuthCookie(ViewAllFiles))
	router.GET("/me/collections", HasAuthCookie(MeCollectionsGet))
	router.PUT("/me/collections", HasAuthCookie(MeCollectionsPut))
	router.POST("/me/collections", HasAuthCookie(MeCollectionsPost))
	router.GET("/me/collections/:collection", HasAuthCookie(MeCollectionGet))
	router.PUT("/me/collections/:collection", HasAuthCookie(MeCollectionPut))
	router.POST("/me/collections/:collection", HasAuthCookie(MeCollectionPost))

	//Misc
	
	router.GET("/fav/favicon.ico", Ignore)
	
	static := httprouter.New()
	static.ServeFiles("/files/*filepath", http.Dir(streamPath))
	static.ServeFiles("/giant/*filepath", http.Dir(uploadPath))
	static.ServeFiles("/assets/*filepath", http.Dir("./assets"))
	static.ServeFiles("/poster/*filepath", http.Dir(postersDir))
	static.ServeFiles("/userpic/*filepath", http.Dir(userpicPath))
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.ServeFiles("/usrimg/*filepath", http.Dir("usrimg"))
	router.ServeFiles("/img/*filepath", http.Dir("uploads"))
	router.ServeFiles("/public/*filepath", http.Dir("public"))
	router.NotFound = static
//~	router.NotFound = http.FileServer(http.Dir(""))


	log.Println("Starting Server")
    log.Fatal(http.ListenAndServe("localhost:4000", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tmpl.ExecuteTemplate(w,"head.html",nil)
	tmpl.ExecuteTemplate(w,"indexx.html",nil)
}
