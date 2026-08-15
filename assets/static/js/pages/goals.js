let goalsData = null


async function loadGoals(){

    try{

        const data =
        await getLatestGoals()

        if(!data) return

        document.querySelector(
            "[name=age]"
        ).value =
        data.age

        document.querySelector(
            "[name=gender]"
        ).value =
        data.gender

        document.querySelector(
            "[name=height_cm]"
        ).value =
        data.height_cm

        document.querySelector(
            "[name=weight_kg]"
        ).value =
        data.weight_kg

        document.querySelector(
            "[name=activity_level]"
        ).value =
        data.activity_level

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
    document.getElementById(
        "goals-form"
    )

    const data = {

        age:
        Number(form.age.value),

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

        const result =
        await saveGoals(data)

        document
        .getElementById(
            "goals-result"
        )
        .innerHTML =

        `
        <strong>Targets saved</strong>
        <br><br>
        Calories: ${result.maintenance_calories}
        <br>
        Protein: ${result.protein_goal} g
        <br>
        Carbs: ${result.carb_goal} g
        <br>
        Fat: ${result.fat_goal} g
        `

        setGoalsPageMessage(
            "Your goals were saved successfully.",
            "success"
        )

    }
    catch(err){

        console.error(
            "Failed to save goals",
            err
        )

    }

}


async function loadMacroGoals(){

    try{

        const data =
        await getLatestGoals()

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
            data.fiber_goal || 35

        }

    }
    catch(err){

        console.error(
            "Failed to load macro goals",
            err
        )

    }

}


function setGoalsPageMessage(
    message,
    type = ""
){

    const element =
    document.getElementById(
        "goalsPageMessage"
    )

    if(!element){
        return
    }

    element.textContent = message
    element.className = "page-message"

    if(type){
        element.classList.add(type)
    }

}