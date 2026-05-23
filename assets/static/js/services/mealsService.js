async function getMealsByDate(date){

    return apiGet(
        `/meals/getmealsbydate?date=${date}`
    )

}


async function addMealItem(payload){

    return apiPost(
        "/mealitems/addmealitem",
        payload
    )

}


async function getTodayNutrition(){

    return apiGet(
        "/meals/gettodaynutrition"
    )

}