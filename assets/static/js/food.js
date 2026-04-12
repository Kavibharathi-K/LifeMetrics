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

     loadFoods()

}


// show Meals page
function showMealsView() {

    document.getElementById("pageTitle").innerText = "Meals"

    document.getElementById("foodsView").style.display = "none"

    document.getElementById("mealsView").style.display = "block"

}



let selectedMealType = ""


function openMealModal(mealType){

    selectedMealType = mealType

    document.querySelector(
        "#mealModal h3"
    ).innerText =
    "Add food to " + mealType

    document.getElementById("mealSuccessMsg").innerText=""

    document.getElementById("foodNameInput").value=""
    document.getElementById("quantityInput").value=""

    document.getElementById("mealModal").style.display="block"
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

        await loadMacroGoals()

        const res =
        await fetch(
            "/meals/gettodaynutrition"
        )

        const data =
        await res.json()

        const nutrition =
        data[0]

        if(!nutrition) return


        const calories =
        Math.round(
            nutrition.total_calories
        )

        const protein =
        nutrition.total_protein

        const carbs =
        nutrition.total_carbs

        const fat =
        nutrition.total_fat

        const fiber =
        nutrition.total_fiber



        document
        .getElementById(
            "calories-value"
        )
        .innerText =
        calories


        document
        .getElementById(
            "protein-value"
        )
        .innerText =
        protein.toFixed(1) + " g"


        document
        .getElementById(
            "carbs-value"
        )
        .innerText =
        carbs.toFixed(1) + " g"


        document
        .getElementById(
            "fat-value"
        )
        .innerText =
        fat.toFixed(1) + " g"


        document
        .getElementById(
            "fiber-value"
        )
        .innerText =
        fiber.toFixed(1) + " g"



        updateProgress(
            "calories",
            calories,
            goalsData.calories,
            " kcal"
        )

        updateProgress(
            "protein",
            protein,
            goalsData.protein,
            " g"
        )

        updateProgress(
            "carbs",
            carbs,
            goalsData.carbs,
            " g"
        )

        updateProgress(
            "fat",
            fat,
            goalsData.fat,
            " g"
        )

        updateProgress(
            "fiber",
            fiber,
            goalsData.fiber,
            " g"
        )

    }
    catch(err){

        console.error(
            "Failed dashboard load",
            err
        )

    }

}


async function loadGoals(){

    try{

        const res =
        await fetch("/usermetrics/getlatestusermetrics")

        if(!res.ok) return

        const data =
        await res.json()

        if(!data) return

        document.querySelector(
            "[name=age]"
        ).value = data.age

        document.querySelector(
            "[name=gender]"
        ).value = data.gender

        document.querySelector(
            "[name=height_cm]"
        ).value = data.height_cm

        document.querySelector(
            "[name=weight_kg]"
        ).value = data.weight_kg

        document.querySelector(
            "[name=activity_level]"
        ).value = data.activity_level

    }
    catch(err){

        console.error(
            "Failed to load goals",
            err
        )

    }

}

async function submitGoals(event){

    event.preventDefault()

    const form =
    document.getElementById("goals-form")

    const data = {

        age:
        Number(
            form.age.value
        ),

        gender:
        form.gender.value,

        height_cm:
        Number(
            form.height_cm.value
        ),

        weight_kg:
        Number(
            form.weight_kg.value
        ),

        activity_level:
        form.activity_level.value
    }


    try{

        const res =
        await fetch(

            "/usermetrics/addusermetrics",

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


        const result =
        await res.json()


        document
        .getElementById("goals-result")

        .innerHTML =

        `
        Targets Saved

        <br><br>

        Calories:
        ${result.maintenance_calories}

        <br>

        Protein:
        ${result.protein_goal} g

        <br>

        Carbs:
        ${result.carb_goal} g

        <br>

        Fat:
        ${result.fat_goal} g
        `

    }
    catch(err){

        console.error(
            "Failed to save goals",
            err
        )

    }

}


window.onload = function () {

    const subPage =
    document.body.dataset.subpage

    console.log("SubPage:", subPage)


    if (subPage === "today") {

        loadTodayNutrition()
        return
    }


    if (subPage === "goals") {

        loadGoals()

        const form =
        document.getElementById("goals-form")

        if(form){

            form.addEventListener(
                "submit",
                submitGoals
            )
        }

        return
    }


    if(subPage === "foods"){

        loadFoods()

    }

}


let goalsData = null;



async function loadMacroGoals(){

    try{

        const res =
        await fetch("/usermetrics/getlatestusermetrics")

        const data =
        await res.json()

        goalsData = {

            calories:
            data.maintenance_calories,

            protein:
            data.protein_goal,

            carbs:
            data.carb_goal,

            fat:
            data.fat_goal,

            fiber:
            data.fiber_goal || 30

        }

    }
    catch(err){

        console.error(
            "Failed to load macro goals",
            err
        )

    }

}



function updateProgress(
    id,
    value,
    goal,
    unit
){

    if(!goal) return

    const percent =
    (value / goal) * 100


    const bar =
    document.getElementById(
        `${id}-progress`
    )

    const text =
    document.getElementById(
        `${id}-remaining`
    )


    /* width */

    bar.style.width =
    Math.min(percent, 100) + "%"



    /* status text */

    if(value < goal){

        const remaining =
        goal - value

        text.innerText =
        remaining.toFixed(0) +
        unit +
        " left"

        bar.style.background =
        "#2563eb"   // theme blue

    }

    else if(value === goal){

        text.innerText =
        "Goal reached"

        bar.style.background =
        "#16a34a"   // subtle green

    }

    else{

        const extra =
        value - goal

        text.innerText =
        extra.toFixed(0) +
        unit +
        " extra"

        bar.style.background =
        "#dc2626"   // red


        /* allow bar overflow visually */

        bar.style.width =
        Math.min(percent, 130) + "%"

    }

}