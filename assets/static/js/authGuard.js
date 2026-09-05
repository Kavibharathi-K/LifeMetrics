function requireAuthentication(){

    if(!isAuthenticated()){
        window.location.href = "/login"
    }

}