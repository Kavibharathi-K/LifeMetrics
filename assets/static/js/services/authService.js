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

    return !!getToken()

}