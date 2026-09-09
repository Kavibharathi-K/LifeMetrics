async function register(
    data
){

    return apiPost(
        "/auth/register",
        data
    )

}


async function login(
    data
){

    const response =
    await apiPost(
        "/auth/login",
        data
    )

    localStorage.setItem(
        "token",
        response.token
    )

    return response

}


function logout(){

    localStorage.removeItem(
        "token"
    )

}


function getToken(){

    return localStorage.getItem(
        "token"
    )

}


function isAuthenticated(){

    const token = getToken()

    if(!token){
        return false
    }

    try{

        const payload =
        JSON.parse(
            atob(
                token.split(".")[1]
            )
        )

        if(
            payload.exp &&
            Date.now() >= payload.exp * 1000
        ){

            logout()

            return false

        }

        return true

    }
    catch(error){

        logout()

        return false

    }

}