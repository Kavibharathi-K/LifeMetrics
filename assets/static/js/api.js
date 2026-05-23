async function apiGet(url){

    const response = await fetch(url)

    if(!response.ok){
        throw new Error("Request failed")
    }

    return response.json()

}


async function apiPost(url, data){

    const response = await fetch(url, {

        method: "POST",

        headers: {
            "Content-Type": "application/json"
        },

        body: JSON.stringify(data)

    })

    if(!response.ok){
        throw new Error("Request failed")
    }

    return response.json()

}