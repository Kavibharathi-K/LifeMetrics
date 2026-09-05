async function apiRequest(
    url,
    options = {}
){

    const token =
    localStorage.getItem(
        "token"
    )


    const headers = {
        ...(options.headers || {})
    }


    if(token){

        headers.Authorization =
        `Bearer ${token}`

    }


    const response =
    await fetch(
        url,
        {
            ...options,
            headers
        }
    )


    if(response.status === 204){
        return null
    }


    const contentType =
    response.headers.get(
        "content-type"
    )


    let data = null


    if(
        contentType &&
        contentType.includes(
            "application/json"
        )
    ){

        data =
        await response.json()

    }
    else{

        data =
        await response.text()

    }


    if(!response.ok){

        let message =
        "Request failed"


        if(
            data &&
            typeof data === "object" &&
            data.error
        ){

            message =
            data.error

        }
        else if(
            typeof data === "string" &&
            data.trim()
        ){

            message =
            data

        }


        throw new Error(
            message
        )

    }


    return data

}


async function apiGet(
    url
){

    return apiRequest(
        url,
        {
            method: "GET"
        }
    )

}


async function apiPost(
    url,
    data
){

    return apiRequest(
        url,
        {
            method: "POST",

            headers: {
                "Content-Type":
                    "application/json"
            },

            body:
                JSON.stringify(data)
        }
    )

}


async function apiPut(
    url,
    data
){

    return apiRequest(
        url,
        {
            method: "PUT",

            headers: {
                "Content-Type":
                    "application/json"
            },

            body:
                JSON.stringify(data)
        }
    )

}


async function apiDelete(
    url
){

    return apiRequest(
        url,
        {
            method: "DELETE"
        }
    )

}