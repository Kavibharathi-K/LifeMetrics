async function getFoods(){

    return apiGet("/foods/listfoods")

}

async function searchUSDAFoods(query) {

    const response =
    await fetch(
        `/usda/searchfoods?query=${encodeURIComponent(query)}`
    )

    if(!response.ok){

        throw new Error(
            "Failed to search foods"
        )

    }

    return await response.json()

}