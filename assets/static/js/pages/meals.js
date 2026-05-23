let selectedMealType = ""


function showMealsView(){

    document.getElementById(
        "pageTitle"
    ).innerText =
    "Meals"

    document.getElementById(
        "foodsView"
    ).style.display =
    "none"

    document.getElementById(
        "mealsView"
    ).style.display =
    "block"

}


function openMealModal(mealType){

    selectedMealType = mealType

    document.querySelector(
        "#mealModal h3"
    ).innerText =
    "Add food to " + mealType

    document.getElementById(
        "mealSuccessMsg"
    ).innerText = ""

    document.getElementById(
        "foodNameInput"
    ).value = ""

    document.getElementById(
        "quantityInput"
    ).value = ""

    document.getElementById(
        "mealModal"
    ).style.display =
    "block"

}


function closeMealModal(){

    document.getElementById(
        "mealModal"
    ).style.display =
    "none"

}


async function submitMealItem(){

    const foodName =
    document.getElementById(
        "foodNameInput"
    ).value

    const quantity =
    document.getElementById(
        "quantityInput"
    ).value

    const measurementType =
    document.getElementById(
        "measurementTypeInput"
    ).value


    if(!foodName || !quantity){

        alert(
            "Please enter food name and quantity"
        )

        return

    }


    const payload = {

        meal_date:
        new Date()
        .toISOString()
        .split('T')[0],

        meal_type:
        selectedMealType,

        food_name:
        foodName,

        quantity:
        parseFloat(quantity),

        measurement_type:
        measurementType

    }


    try{

        const data =
        await addMealItem(payload)

        console.log(
            "Meal item added:",
            data
        )

        const msg =
        document.getElementById(
            "mealSuccessMsg"
        )

        msg.innerText =
        "Successfully logged meal item"

        msg.style.display =
        "block"


        setTimeout(() => {

            msg.innerText = ""

            closeMealModal()

            loadMeals()

        }, 1200)

    }
    catch(err){

        console.error(err)

    }

}


async function loadMeals(){

    try{

        const today =
        new Date()
        .toISOString()
        .split("T")[0]

        const data =
        await getMealsByDate(today)

        renderMealCard(
            "Breakfast",
            data.breakfast
        )

        renderMealCard(
            "Lunch",
            data.lunch
        )

        renderMealCard(
            "Dinner",
            data.dinner
        )

        renderMealCard(
            "Snacks",
            data.snack
        )

    }
    catch(err){

        console.error(
            "Failed loading meals",
            err
        )

    }

}