async function loadFoods() {

    try {

        const response = await fetch("/foods/listfoods")

        const foods = await response.json()

        console.log("foods:", foods)

        const tableBody = document.querySelector("#foodTable tbody")

        tableBody.innerHTML = ""

        foods.forEach(food => {

            const row = `

                <tr>

                    <td>${food.name}</td>

                    <td>
                        ${food.base_quantity}
                        ${food.measurement_type}
                    </td>

                    <td>${food.calories || 0}</td>

                    <td>${food.protein || 0}</td>

                    <td>${food.carbs || 0}</td>

                    <td>${food.fat || 0}</td>

                </tr>

            `

            tableBody.innerHTML += row

        })

    }
    catch(error) {

        console.error("Error loading foods:", error)

    }

}


// show Available Foods table
function showFoodsView() {

    document.getElementById("pageTitle").innerText = "Available Foods"

    document.getElementById("foodsView").style.display = "block"

    document.getElementById("mealsView").style.display = "none"

}


// show Meals page
function showMealsView() {

    document.getElementById("pageTitle").innerText = "Meals"

    document.getElementById("foodsView").style.display = "none"

    document.getElementById("mealsView").style.display = "block"

}



let selectedMealType = ""


// open modal
function openMealModal(mealType) {

    selectedMealType = mealType

    document.getElementById("foodNameInput").value = ""
    document.getElementById("quantityInput").value = ""

    document.getElementById("mealModal").style.display = "block"

}


// close modal
function closeMealModal() {

    document.getElementById("mealModal").style.display = "none"

}


// submit data
function submitMealItem() {

    const foodName = document.getElementById("foodNameInput").value

    const quantity = document.getElementById("quantityInput").value

    const measurementType = document.getElementById("measurementTypeInput").value


    if (!foodName || !quantity) {

        alert("Please enter food name and quantity")
        return

    }


    const payload = {

        meal_date: new Date().toISOString().split('T')[0],

        meal_type: selectedMealType,

        food_name: foodName,

        quantity: parseFloat(quantity),

        measurement_type: measurementType

    }


    fetch("/mealitems/addmealitem", {

        method: "POST",

        headers: {
            "Content-Type": "application/json"
        },

        body: JSON.stringify(payload)

    })
    .then(res => res.json())
    .then(data => {

        console.log("Meal item added:", data)

        const msg = document.getElementById("mealSuccessMsg")

        msg.innerText = "Successfully logged meal item"

        msg.style.display = "block"

        // wait 2 seconds before closing modal
        setTimeout(() => {

            msg.innerText = ""

            closeMealModal()

        }, 2000)

    })
    .catch(err => console.error(err))

}


async function loadTodayNutrition(){

    try{

        const res = await fetch("/meals/gettodaynutrition");

        const data = await res.json();

        const nutrition = data[0];

        if(!nutrition) return;

        document.getElementById("calories-value").innerText =
            Math.round(nutrition.total_calories);

        document.getElementById("protein-value").innerText =
            nutrition.total_protein.toFixed(1) + " g";

        document.getElementById("carbs-value").innerText =
            nutrition.total_carbs.toFixed(1) + " g";

        document.getElementById("fat-value").innerText =
            nutrition.total_fat.toFixed(1) + " g";

        document.getElementById("fiber-value").innerText =
            nutrition.total_fiber.toFixed(1) + " g";

    }
    catch(err){

        console.error("Failed to load dashboard data", err);

    }

}


window.onload = function () {

    const subPage = document.body.dataset.subpage

    console.log("SubPage:", subPage)

    if (subPage === "today") {

        loadTodayNutrition()

    }
    else if (subPage === "meals") {

        showMealsView()

    }
    else {

        showFoodsView()
        loadFoods()

    }

}